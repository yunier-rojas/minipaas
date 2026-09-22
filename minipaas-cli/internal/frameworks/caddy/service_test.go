package caddy

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

type recordRunner struct {
	cmds   [][]string
	output string
	runErr error
	outErr error
}

func (r *recordRunner) Run(cmd []string, verbose bool) error {
	r.cmds = append(r.cmds, cmd)
	return r.runErr
}

func (r *recordRunner) RunWithInput(cmd []string, input []byte, verbose bool) error {
	r.cmds = append(r.cmds, cmd)
	return r.runErr
}

func (r *recordRunner) RunOutput(cmd []string, verbose bool) (string, error) {
	r.cmds = append(r.cmds, cmd)
	return r.output, r.outErr
}

func TestApply_LoadsConfigIntoContainer(t *testing.T) {
	env := t.TempDir()
	payload := `{"apps":{}}`
	if err := os.WriteFile(filepath.Join(env, configFile), []byte(payload), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	runner := &recordRunner{output: "cid1\ncid2\n"}
	if err := NewService(runner).Apply(env, "minipaas", true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(runner.cmds) != 2 {
		t.Fatalf("expected 2 commands, got %#v", runner.cmds)
	}
	wantPS := []string{"docker", "ps", "--filter", "label=com.docker.swarm.service.name=minipaas_caddy", "--format", "{{.ID}}"}
	if !reflect.DeepEqual(runner.cmds[0], wantPS) {
		t.Fatalf("ps cmd = %#v, want %#v", runner.cmds[0], wantPS)
	}
	exec := runner.cmds[1]
	if len(exec) != 11 || exec[0] != "docker" || exec[1] != "exec" || exec[2] != "-i" || exec[3] != "cid1" {
		t.Fatalf("exec cmd = %#v", exec)
	}
	if exec[4] != "/usr/bin/wget" || exec[10] != "http://127.0.0.1:2019/load" {
		t.Fatalf("exec cmd = %#v", exec)
	}
	if exec[9] != "--post-data="+payload {
		t.Fatalf("payload arg = %q", exec[9])
	}
}

func TestApply_MissingConfig(t *testing.T) {
	runner := &recordRunner{}
	if err := NewService(runner).Apply(t.TempDir(), "minipaas", false); err == nil {
		t.Fatalf("expected error for missing config")
	}
	if len(runner.cmds) != 0 {
		t.Fatalf("no command should run: %#v", runner.cmds)
	}
}

func TestApply_NoContainer(t *testing.T) {
	env := t.TempDir()
	if err := os.WriteFile(filepath.Join(env, configFile), []byte("{}"), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	runner := &recordRunner{output: "\n"}
	if err := NewService(runner).Apply(env, "minipaas", false); err == nil {
		t.Fatalf("expected error when no container is running")
	}
}
