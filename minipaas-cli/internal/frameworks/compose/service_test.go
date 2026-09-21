package compose

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/compose-spec/compose-go/v2/types"

	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/entities"
)

func TestFilesForEnv_ResolvesRelativeToEnv(t *testing.T) {
	dir := t.TempDir()
	cfg := entities.MinipaasConfig{Compose: entities.ComposeConfig{Files: []string{
		"compose.common.yml",
		"sub/compose.apps.yml",
		"/abs/compose.other.yml",
	}}}

	got := NewManagerService().FilesForEnv(dir, cfg)
	want := []string{
		filepath.Join(dir, "compose.common.yml"),
		filepath.Join(dir, "sub", "compose.apps.yml"),
		"/abs/compose.other.yml",
	}
	if len(got) != len(want) {
		t.Fatalf("unexpected compose files: %#v", got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("file %d = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestFilesForEnv_PreservesOrderAndEmpty(t *testing.T) {
	dir := t.TempDir()
	cfg := entities.MinipaasConfig{Compose: entities.ComposeConfig{Files: []string{"b.yml", "a.yml"}}}
	got := NewManagerService().FilesForEnv(dir, cfg)
	if len(got) != 2 || got[0] != filepath.Join(dir, "b.yml") || got[1] != filepath.Join(dir, "a.yml") {
		t.Fatalf("order not preserved: %#v", got)
	}

	if empty := NewManagerService().FilesForEnv(dir, entities.MinipaasConfig{}); len(empty) != 0 {
		t.Fatalf("expected no files, got %#v", empty)
	}
}

func TestDiscoverProjectFiles(t *testing.T) {
	root := t.TempDir()
	env := filepath.Join(root, "minipaas")
	if err := os.MkdirAll(env, 0755); err != nil {
		t.Fatal(err)
	}

	base := filepath.Join(root, "compose.yaml")
	build := filepath.Join(root, "compose.build.yaml")
	for _, f := range []string{base, build} {
		if err := os.WriteFile(f, []byte("services: {}\n"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	got := NewManagerService().DiscoverProjectFiles(env)
	want := []string{base, build}
	if len(got) != len(want) {
		t.Fatalf("discovered = %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("file %d = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestDiscoverProjectFiles_None(t *testing.T) {
	root := t.TempDir()
	if got := NewManagerService().DiscoverProjectFiles(filepath.Join(root, "dev")); len(got) != 0 {
		t.Fatalf("expected no files, got %#v", got)
	}
}

func TestGroupServicesByFile_LastDefinitionWins(t *testing.T) {
	dir := t.TempDir()
	base := filepath.Join(dir, "compose.yaml")
	apps := filepath.Join(dir, "compose.apps.yaml")

	if err := os.WriteFile(base, []byte("services:\n  api:\n    image: x\n  web:\n    image: y\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(apps, []byte("services:\n  api:\n    image: z\n"), 0644); err != nil {
		t.Fatal(err)
	}

	svcMap, missing := NewManagerService().GroupServicesByFile([]string{base, apps}, []string{"api", "web"}, "minipaas")
	if len(missing) != 0 {
		t.Fatalf("unexpected missing: %#v", missing)
	}
	if len(svcMap[apps]) != 1 || svcMap[apps][0] != "api" {
		t.Fatalf("api should map to the last defining file: %#v", svcMap)
	}
	if len(svcMap[base]) != 1 || svcMap[base][0] != "web" {
		t.Fatalf("web should map to base: %#v", svcMap)
	}
}

func TestGroupServicesByFile(t *testing.T) {
	dir := t.TempDir()
	f1 := filepath.Join(dir, "a.yml")
	f2 := filepath.Join(dir, "b.yml")

	if err := os.WriteFile(f1, []byte("version: '3.9'\nservices:\n  api:\n    image: x\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(f2, []byte("version: '3.9'\nservices:\n  web:\n    image: y\n"), 0644); err != nil {
		t.Fatal(err)
	}

	svcMap, missing := NewManagerService().GroupServicesByFile([]string{f1, f2}, []string{"api", "web", "worker"}, "minipaas")
	if len(missing) != 1 || missing[0] != "worker" {
		t.Fatalf("missing mismatch: %#v", missing)
	}
	if len(svcMap[f1]) != 1 || svcMap[f1][0] != "api" {
		t.Fatalf("f1 map mismatch: %#v", svcMap[f1])
	}
	if len(svcMap[f2]) != 1 || svcMap[f2][0] != "web" {
		t.Fatalf("f2 map mismatch: %#v", svcMap[f2])
	}
}

func TestLoadComposeFile_StripsDefaultNetwork(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "one.yml")
	data := "version: '3.9'\nnetworks:\n  default: {}\nservices:\n  web:\n    image: x\n    networks:\n      - default\n"
	if err := os.WriteFile(f, []byte(data), 0644); err != nil {
		t.Fatal(err)
	}
	p, _, err := loadComposeFile(f, "minipaas")
	if err != nil {
		t.Fatalf("loadComposeFile: %v", err)
	}
	if _, ok := p.Networks["default"]; ok {
		t.Fatalf("default network should be stripped")
	}
	if svc := p.Services["web"]; svc.Networks != nil {
		if _, ok := svc.Networks["default"]; ok {
			t.Fatalf("service default network should be stripped")
		}
	}
}

func TestSaveComposeFile_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "round.yml")
	p := &types.Project{Services: make(types.Services)}
	p.Services["a"] = types.ServiceConfig{Image: "busybox"}
	if _, err := saveComposeFile(f, p); err != nil {
		t.Fatalf("saveComposeFile: %v", err)
	}
	if _, err := os.Stat(f); err != nil {
		t.Fatalf("file not written: %v", err)
	}
	p2, _, err := loadComposeFile(f, "minipaas")
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	if _, ok := p2.Services["a"]; !ok {
		t.Fatalf("service missing after roundtrip")
	}
}

func TestAddComposeSecretAndConfig_NoDupes(t *testing.T) {
	p := &types.Project{Services: make(types.Services)}
	p.Services["api"] = types.ServiceConfig{Image: "x"}

	if err := addComposeSecret(p, "secret-hash", "env", []string{"api"}); err != nil {
		t.Fatalf("add secret: %v", err)
	}
	if err := addComposeSecret(p, "secret-hash", "env", []string{"api"}); err != nil {
		t.Fatalf("add secret twice: %v", err)
	}
	if len(p.Secrets) != 1 || len(p.Services["api"].Secrets) != 1 {
		t.Fatalf("secret dupes detected: secrets=%d svcRefs=%d", len(p.Secrets), len(p.Services["api"].Secrets))
	}

	if err := addComposeConfig(p, "cfg-hash", "app.conf", []string{"api"}); err != nil {
		t.Fatalf("add config: %v", err)
	}
	if err := addComposeConfig(p, "cfg-hash", "app.conf", []string{"api"}); err != nil {
		t.Fatalf("add config twice: %v", err)
	}
	if len(p.Configs) != 1 || len(p.Services["api"].Configs) != 1 {
		t.Fatalf("config dupes detected: configs=%d svcRefs=%d", len(p.Configs), len(p.Services["api"].Configs))
	}
}

func TestAddComposeSecret_ReplacesByTarget(t *testing.T) {
	p := &types.Project{Services: make(types.Services)}
	p.Services["api"] = types.ServiceConfig{Image: "x"}

	if err := addComposeSecret(p, "password.11111111", "password", []string{"api"}); err != nil {
		t.Fatalf("first add: %v", err)
	}
	if err := addComposeSecret(p, "password.22222222", "password", []string{"api"}); err != nil {
		t.Fatalf("second add: %v", err)
	}

	if len(p.Secrets) != 1 {
		t.Fatalf("old secret not pruned: %#v", p.Secrets)
	}
	if _, ok := p.Secrets["password.22222222"]; !ok {
		t.Fatalf("latest secret missing: %#v", p.Secrets)
	}
	if len(p.Services["api"].Secrets) != 1 {
		t.Fatalf("service has stale secret refs: %#v", p.Services["api"].Secrets)
	}
	ref := p.Services["api"].Secrets[0]
	if ref.Source != "password.22222222" || ref.Target != "password" {
		t.Fatalf("service ref = %#v, want source password.22222222 target password", ref)
	}
}

func TestAddComposeConfig_ReplacesByTarget(t *testing.T) {
	p := &types.Project{Services: make(types.Services)}
	p.Services["api"] = types.ServiceConfig{Image: "x"}

	if err := addComposeConfig(p, "app.11111111", "app.conf", []string{"api"}); err != nil {
		t.Fatalf("first add: %v", err)
	}
	if err := addComposeConfig(p, "app.22222222", "app.conf", []string{"api"}); err != nil {
		t.Fatalf("second add: %v", err)
	}

	if len(p.Configs) != 1 {
		t.Fatalf("old config not pruned: %#v", p.Configs)
	}
	if _, ok := p.Configs["app.22222222"]; !ok {
		t.Fatalf("latest config missing: %#v", p.Configs)
	}
	if len(p.Services["api"].Configs) != 1 {
		t.Fatalf("service has stale config refs: %#v", p.Services["api"].Configs)
	}
	ref := p.Services["api"].Configs[0]
	if ref.Source != "app.22222222" || ref.Target != "app.conf" {
		t.Fatalf("service ref = %#v, want source app.22222222 target app.conf", ref)
	}
}

func TestBuildCommonImage(t *testing.T) {
	tests := []struct {
		name     string
		image    string
		registry string
		want     string
	}{
		{"plain with registry", "example:latest", "ghcr.io/org", "${MINIPAAS_IMAGE_PREFIX}/example:${MINIPAAS_DEPLOY_VERSION}"},
		{"plain without registry", "example:latest", "", "example:${MINIPAAS_DEPLOY_VERSION}"},
		{"no tag with registry", "example", "ghcr.io/org", "${MINIPAAS_IMAGE_PREFIX}/example:${MINIPAAS_DEPLOY_VERSION}"},
		{"qualified image", "ghcr.io/acme/example:1.2.3", "ghcr.io/org", "ghcr.io/acme/example:${MINIPAAS_DEPLOY_VERSION}"},
	}
	for _, tt := range tests {
		if got := buildCommonImage(tt.image, tt.registry); got != tt.want {
			t.Fatalf("buildCommonImage(%q, %q) = %q, want %q", tt.image, tt.registry, got, tt.want)
		}
	}
}

func TestCreateDeployService(t *testing.T) {
	svc := types.ServiceConfig{Image: "example:latest"}

	built := createDeployService(svc, true, "ghcr.io/org", "minipaas")
	if built.Image != "${MINIPAAS_IMAGE_PREFIX}/example:${MINIPAAS_DEPLOY_VERSION}" {
		t.Fatalf("built image = %q", built.Image)
	}
	if _, ok := built.Networks["minipaas_network"]; !ok {
		t.Fatalf("network not set: %#v", built.Networks)
	}
	if built.Labels[entities.AlertCPUThresholdLabel] != entities.DefaultAlertCPUThreshold {
		t.Fatalf("cpu label = %q, want %q", built.Labels[entities.AlertCPUThresholdLabel], entities.DefaultAlertCPUThreshold)
	}
	if built.Labels[entities.AlertMemoryThresholdLabel] != entities.DefaultAlertMemoryThreshold {
		t.Fatalf("memory label = %q, want %q", built.Labels[entities.AlertMemoryThresholdLabel], entities.DefaultAlertMemoryThreshold)
	}

	notBuilt := createDeployService(svc, false, "ghcr.io/org", "minipaas")
	if notBuilt.Image != "" {
		t.Fatalf("non-build image = %q, want empty", notBuilt.Image)
	}
}

func TestCreateDeployService_SkipsExistingAlertLabels(t *testing.T) {
	svc := types.ServiceConfig{
		Image: "example:latest",
		Labels: map[string]string{
			entities.AlertCPUThresholdLabel: "75",
		},
	}

	got := createDeployService(svc, false, "", "minipaas")
	if _, ok := got.Labels[entities.AlertCPUThresholdLabel]; ok {
		t.Fatalf("cpu label should keep its source value: %#v", got.Labels)
	}
	if got.Labels[entities.AlertMemoryThresholdLabel] != entities.DefaultAlertMemoryThreshold {
		t.Fatalf("memory label = %q, want %q", got.Labels[entities.AlertMemoryThresholdLabel], entities.DefaultAlertMemoryThreshold)
	}
}

func TestCreateDeployService_CustomNamespace(t *testing.T) {
	built := createDeployService(types.ServiceConfig{Image: "example:latest"}, true, "", "acme")
	if _, ok := built.Networks["acme_network"]; !ok {
		t.Fatalf("network not namespaced: %#v", built.Networks)
	}
	if _, ok := built.Networks["minipaas_network"]; ok {
		t.Fatalf("default network should be absent: %#v", built.Networks)
	}
}

func TestBuildApps(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "compose.yaml")
	data := "services:\n  example:\n    image: example:latest\n    build:\n      context: .\n  db:\n    image: postgres:17\n"
	if err := os.WriteFile(src, []byte(data), 0644); err != nil {
		t.Fatal(err)
	}

	env := filepath.Join(dir, "dev")
	name, err := NewManagerService().BuildApps(env, []string{src}, "ghcr.io/org", "minipaas")
	if err != nil {
		t.Fatalf("BuildApps: %v", err)
	}
	if name != entities.AppsFile {
		t.Fatalf("name = %q, want %q", name, entities.AppsFile)
	}

	content, err := os.ReadFile(filepath.Join(env, entities.AppsFile))
	if err != nil {
		t.Fatalf("read apps: %v", err)
	}
	out := string(content)
	if !strings.Contains(out, "${MINIPAAS_IMAGE_PREFIX}/example:${MINIPAAS_DEPLOY_VERSION}") {
		t.Fatalf("built image missing:\n%s", out)
	}
	if !strings.Contains(out, "minipaas_network") {
		t.Fatalf("network missing:\n%s", out)
	}
	if !strings.Contains(out, entities.AlertCPUThresholdLabel) || !strings.Contains(out, entities.AlertMemoryThresholdLabel) {
		t.Fatalf("monitoring labels missing:\n%s", out)
	}
	if strings.Contains(out, "postgres:17") {
		t.Fatalf("non-build image should be dropped:\n%s", out)
	}
}

func TestAddComposeWorkerDeploy(t *testing.T) {
	p := &types.Project{Services: make(types.Services)}
	p.Services["api"] = types.ServiceConfig{Image: "x"}

	if err := addComposeWorkerDeploy(p, "api"); err != nil {
		t.Fatalf("addComposeWorkerDeploy: %v", err)
	}
	d := p.Services["api"].Deploy
	if d == nil || d.Mode != "replicated" || d.Replicas == nil || *d.Replicas != 1 {
		t.Fatalf("unexpected deploy: %#v", d)
	}
	if d.RestartPolicy == nil || d.RestartPolicy.Condition != "any" {
		t.Fatalf("unexpected restart policy: %#v", d.RestartPolicy)
	}
	if err := addComposeWorkerDeploy(p, "missing"); err == nil {
		t.Fatalf("expected error for missing service")
	}
}

func TestAddComposeJobDeploy(t *testing.T) {
	p := &types.Project{Services: make(types.Services)}
	p.Services["migrate"] = types.ServiceConfig{Image: "x"}

	if err := addComposeJobDeploy(p, "migrate"); err != nil {
		t.Fatalf("addComposeJobDeploy: %v", err)
	}
	d := p.Services["migrate"].Deploy
	if d == nil || d.RestartPolicy == nil || d.RestartPolicy.Condition != "on-failure" {
		t.Fatalf("unexpected deploy: %#v", d)
	}
}

func TestAddComposeCronDeploy(t *testing.T) {
	p := &types.Project{Services: make(types.Services)}
	p.Services["cleanup"] = types.ServiceConfig{Image: "x"}

	if err := addComposeCronDeploy(p, "cleanup", "*/5 * * * *"); err != nil {
		t.Fatalf("addComposeCronDeploy: %v", err)
	}
	d := p.Services["cleanup"].Deploy
	if d == nil || d.Replicas == nil || *d.Replicas != 0 {
		t.Fatalf("unexpected deploy: %#v", d)
	}
	if d.Labels["swarm.cronjob.enable"] != "true" || d.Labels["swarm.cronjob.schedule"] != "*/5 * * * *" {
		t.Fatalf("unexpected labels: %#v", d.Labels)
	}
}

func TestAddComposeRouteDeploy(t *testing.T) {
	p := &types.Project{Services: make(types.Services)}
	p.Services["api"] = types.ServiceConfig{Image: "x"}

	if err := addComposeRouteDeploy(p, "api", "8080"); err != nil {
		t.Fatalf("addComposeRouteDeploy: %v", err)
	}
	svc := p.Services["api"]
	if svc.Deploy == nil || svc.Deploy.Replicas == nil || *svc.Deploy.Replicas != 2 {
		t.Fatalf("unexpected deploy: %#v", svc.Deploy)
	}
	if svc.HealthCheck == nil || len(svc.HealthCheck.Test) != 2 {
		t.Fatalf("unexpected healthcheck: %#v", svc.HealthCheck)
	}
	if svc.HealthCheck.Test[1] != "wget -qO- --spider http://127.0.0.1:8080" {
		t.Fatalf("unexpected healthcheck test: %#v", svc.HealthCheck.Test)
	}
	if err := addComposeRouteDeploy(p, "missing", "80"); err == nil {
		t.Fatalf("expected error for missing service")
	}
}

func TestAddRouteDeploy_WritesAppsFile(t *testing.T) {
	env := t.TempDir()
	apps := filepath.Join(env, entities.AppsFile)
	if err := os.WriteFile(apps, []byte("services:\n  api:\n    image: x\n"), 0644); err != nil {
		t.Fatal(err)
	}

	fn, err := NewManagerService().AddRouteDeploy(env, "api", "8080", "minipaas")
	if err != nil {
		t.Fatalf("AddRouteDeploy: %v", err)
	}
	if fn != apps {
		t.Fatalf("file = %q, want %q", fn, apps)
	}

	p, _, err := loadComposeFile(apps, "minipaas")
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if p.Services["api"].HealthCheck == nil || p.Services["api"].Deploy == nil {
		t.Fatalf("route deploy config not written")
	}
}

func TestAddWorkerDeploy_WritesAppsFile(t *testing.T) {
	env := t.TempDir()
	apps := filepath.Join(env, entities.AppsFile)
	if err := os.WriteFile(apps, []byte("services:\n  api:\n    image: x\n"), 0644); err != nil {
		t.Fatal(err)
	}

	fn, err := NewManagerService().AddWorkerDeploy(env, []string{"api"}, "minipaas")
	if err != nil {
		t.Fatalf("AddWorkerDeploy: %v", err)
	}
	if fn != apps {
		t.Fatalf("file = %q, want %q", fn, apps)
	}

	p, _, err := loadComposeFile(apps, "minipaas")
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if p.Services["api"].Deploy == nil {
		t.Fatalf("deploy config not written")
	}
}

func TestAddComposeSecretMissingService(t *testing.T) {
	p := &types.Project{Services: make(types.Services)}
	if err := addComposeSecret(p, "sec", "file", []string{"missing"}); err == nil {
		t.Fatalf("expected error for missing service")
	}
}

func TestAddComposeConfigMissingService(t *testing.T) {
	p := &types.Project{Services: make(types.Services)}
	if err := addComposeConfig(p, "cfg", "file", []string{"missing"}); err == nil {
		t.Fatalf("expected error for missing service")
	}
}
