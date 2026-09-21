package deploy

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/compose-spec/compose-go/v2/types"
)

type recordRunner struct {
	cmds [][]string
	err  error
}

func (r *recordRunner) Run(cmd []string, verbose bool) error {
	r.cmds = append(r.cmds, cmd)
	return r.err
}

func (r *recordRunner) RunWithInput(cmd []string, input []byte, verbose bool) error {
	r.cmds = append(r.cmds, cmd)
	return r.err
}

func (r *recordRunner) RunOutput(cmd []string, verbose bool) (string, error) {
	r.cmds = append(r.cmds, cmd)
	return "", r.err
}

func writeCompose(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	fn := filepath.Join(dir, "compose.apps.yaml")
	if err := os.WriteFile(fn, []byte(content), 0644); err != nil {
		t.Fatalf("write compose: %v", err)
	}
	return fn
}

func strptr(s string) *string { return &s }

func TestBuildCommandFromService(t *testing.T) {
	svc := types.ServiceConfig{Image: "repo/app:old"}
	dockerfile := "Dockerfile.prod"
	svc.Build = &types.BuildConfig{
		Context:    "./app",
		Dockerfile: dockerfile,
		Args: map[string]*string{
			"FOO":   strptr("bar"),
			"EMPTY": nil,
		},
	}

	cmd := buildCommandFromService(svc)
	wantDf := filepath.Join(svc.Build.Context, svc.Build.Dockerfile)

	has := func(x string) bool {
		for _, a := range cmd {
			if a == x {
				return true
			}
		}
		return false
	}
	if !has("docker") || !has("build") {
		t.Fatalf("missing docker build: %#v", cmd)
	}
	if !has("-t") || !has("repo/app:old") {
		t.Fatalf("missing tag: %#v", cmd)
	}
	if !has("-f") || !has(wantDf) {
		t.Fatalf("missing dockerfile: %#v", cmd)
	}
	if !has(svc.Build.Context) {
		t.Fatalf("missing context: %#v", cmd)
	}
	if !has("--build-arg") || !has("FOO=bar") || !has("EMPTY=") {
		t.Fatalf("missing build args: %#v", cmd)
	}
}

func TestBuildCommandFromService_NoBuild(t *testing.T) {
	cmd := buildCommandFromService(types.ServiceConfig{Image: "repo/app:1"})
	want := []string{"docker", "build", "--network", "host", "-t", "repo/app:1", "."}
	if !reflect.DeepEqual(cmd, want) {
		t.Fatalf("cmd = %#v, want %#v", cmd, want)
	}
}

func TestRollout_BuildsStackDeployArgs(t *testing.T) {
	files := []string{"a.yml", "b.yml"}
	runner := &recordRunner{}
	if err := NewService(runner).Rollout(files, "minipaas", true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"docker", "stack", "deploy", "--with-registry-auth", "-c", "a.yml", "-c", "b.yml", "minipaas"}
	if len(runner.cmds) != 1 || !reflect.DeepEqual(runner.cmds[0], want) {
		t.Fatalf("cmds = %#v, want %#v", runner.cmds, want)
	}
}

func TestBuild_RunsDockerForServicesWithBuild(t *testing.T) {
	fn := writeCompose(t, "services:\n  api:\n    image: ghcr.io/example/api:1\n    build:\n      context: ../api\n  web:\n    image: ghcr.io/example/web:1\n")
	runner := &recordRunner{}

	if err := NewService(runner).Build([]string{fn}, "ghcr.io/example", "minipaas", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(runner.cmds) != 2 {
		t.Fatalf("expected build and push commands, got %#v", runner.cmds)
	}
	wantBuild := []string{"docker", "build", "--network", "host", "-t", "ghcr.io/example/api:1", "-f", "../api/Dockerfile", "../api"}
	if !reflect.DeepEqual(runner.cmds[0], wantBuild) {
		t.Fatalf("build cmd = %#v, want %#v", runner.cmds[0], wantBuild)
	}
	wantPush := []string{"docker", "push", "ghcr.io/example/api:1"}
	if !reflect.DeepEqual(runner.cmds[1], wantPush) {
		t.Fatalf("push cmd = %#v, want %#v", runner.cmds[1], wantPush)
	}
}

func TestBuild_WithoutRegistrySkipsPush(t *testing.T) {
	fn := writeCompose(t, "services:\n  api:\n    image: ghcr.io/example/api:1\n    build:\n      context: ../api\n")
	runner := &recordRunner{}

	if err := NewService(runner).Build([]string{fn}, "", "minipaas", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(runner.cmds) != 1 {
		t.Fatalf("expected only a build command, got %#v", runner.cmds)
	}
	if runner.cmds[0][0] != "docker" || runner.cmds[0][1] != "build" {
		t.Fatalf("unexpected command: %#v", runner.cmds[0])
	}
}

func TestCanary_UpdatesService(t *testing.T) {
	fn := writeCompose(t, "services:\n  api:\n    image: ghcr.io/example/api:1\n")
	runner := &recordRunner{}

	if err := NewService(runner).Canary([]string{fn}, []string{"api"}, 2, "minipaas", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"docker", "service", "update", "--with-registry-auth", "--replicas", "2", "--image", "ghcr.io/example/api:1", "minipaas_api"}
	if len(runner.cmds) != 1 || !reflect.DeepEqual(runner.cmds[0], want) {
		t.Fatalf("cmds = %#v, want %#v", runner.cmds, want)
	}
}

func TestRollout_CustomNamespace(t *testing.T) {
	runner := &recordRunner{}
	if err := NewService(runner).Rollout([]string{"a.yml"}, "acme", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"docker", "stack", "deploy", "--with-registry-auth", "-c", "a.yml", "acme"}
	if len(runner.cmds) != 1 || !reflect.DeepEqual(runner.cmds[0], want) {
		t.Fatalf("cmds = %#v, want %#v", runner.cmds, want)
	}
}

func TestCanary_CustomNamespace(t *testing.T) {
	fn := writeCompose(t, "services:\n  api:\n    image: ghcr.io/example/api:1\n")
	runner := &recordRunner{}

	if err := NewService(runner).Canary([]string{fn}, []string{"api"}, 1, "acme", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"docker", "service", "update", "--with-registry-auth", "--replicas", "1", "--image", "ghcr.io/example/api:1", "acme_api"}
	if len(runner.cmds) != 1 || !reflect.DeepEqual(runner.cmds[0], want) {
		t.Fatalf("cmds = %#v, want %#v", runner.cmds, want)
	}
}

func TestCanary_MissingService(t *testing.T) {
	fn := writeCompose(t, "services:\n  api:\n    image: ghcr.io/example/api:1\n")
	runner := &recordRunner{}

	if err := NewService(runner).Canary([]string{fn}, []string{"missing"}, 1, "minipaas", false); err == nil {
		t.Fatalf("expected error for missing service")
	}
	if len(runner.cmds) != 0 {
		t.Fatalf("no command should run: %#v", runner.cmds)
	}
}
