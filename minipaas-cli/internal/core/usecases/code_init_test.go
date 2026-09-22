package usecases

import (
	"errors"
	"reflect"
	"testing"

	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/entities"
)

const (
	testRegistry  = "ghcr.io/example"
	testNamespace = "minipaas"
)

func newFakeScaffold() *fakeScaffold {
	return &fakeScaffold{base: entities.BaseFiles{
		Common:    "compose.common.yml",
		Postgres:  "compose.postgres.yml",
		Caddy:     "compose.caddy.yml",
		CaddyJSON: "caddy.json",
	}}
}

func TestCodeInitInteractor_InitLocal(t *testing.T) {
	store := &fakeConfigStore{}
	compose := &fakeCompose{appsFile: entities.AppsFile}
	scaffold := newFakeScaffold()
	interactor := NewCodeInitInteractor(store, compose, scaffold)

	if err := interactor.Init("dev", []string{"compose.yaml"}, "", true, testRegistry, testNamespace); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if compose.buildEnv != "dev" || !reflect.DeepEqual(compose.buildFiles, []string{"compose.yaml"}) {
		t.Fatalf("BuildApps args mismatch: %#v", compose)
	}
	if compose.buildRegistry != testRegistry {
		t.Fatalf("BuildApps registry = %q, want %q", compose.buildRegistry, testRegistry)
	}
	if compose.buildNamespace != testNamespace {
		t.Fatalf("BuildApps namespace = %q, want %q", compose.buildNamespace, testNamespace)
	}
	if scaffold.namespace != testNamespace {
		t.Fatalf("WriteBase namespace = %q, want %q", scaffold.namespace, testNamespace)
	}

	want := []string{
		"../compose.yaml",
		"compose.common.yml",
		"compose.apps.yaml",
		"compose.postgres.yml",
		"compose.caddy.yml",
	}
	if !reflect.DeepEqual(store.saved.Compose.Files, want) {
		t.Fatalf("files = %#v, want %#v", store.saved.Compose.Files, want)
	}
	if store.saved.Vars["MINIPAAS_DEPLOY_VERSION"] != defaultDeployVersion {
		t.Fatalf("version = %q", store.saved.Vars["MINIPAAS_DEPLOY_VERSION"])
	}
	if store.saved.Vars["MINIPAAS_IMAGE_PREFIX"] != testRegistry {
		t.Fatalf("image prefix = %q", store.saved.Vars["MINIPAAS_IMAGE_PREFIX"])
	}
	if store.saved.DeployNamespace() != testNamespace {
		t.Fatalf("namespace = %q, want %q", store.saved.Namespace, testNamespace)
	}
	if store.saved.Docker.Host != "" {
		t.Fatalf("expected local docker config to be empty: %#v", store.saved.Docker)
	}
	if store.saveAt != "dev" {
		t.Fatalf("saved at %q, want dev", store.saveAt)
	}
}

func TestCodeInitInteractor_InitRemote(t *testing.T) {
	store := &fakeConfigStore{}
	interactor := NewCodeInitInteractor(store, &fakeCompose{appsFile: entities.AppsFile}, newFakeScaffold())

	if err := interactor.Init("dev", []string{"/abs/compose.yaml"}, "deploy@example.com", false, testRegistry, testNamespace); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if store.saved.Docker.Host != "ssh://deploy@example.com" {
		t.Fatalf("docker mismatch: %#v", store.saved.Docker)
	}
	if store.saved.Compose.Files[0] != "/abs/compose.yaml" {
		t.Fatalf("absolute file not preserved: %#v", store.saved.Compose.Files)
	}
}

func TestCodeInitInteractor_InitRemoteKeepsScheme(t *testing.T) {
	store := &fakeConfigStore{}
	interactor := NewCodeInitInteractor(store, &fakeCompose{appsFile: entities.AppsFile}, newFakeScaffold())

	if err := interactor.Init("dev", []string{"compose.yaml"}, "ssh://deploy@example.com:2222", false, testRegistry, testNamespace); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if store.saved.Docker.Host != "ssh://deploy@example.com:2222" {
		t.Fatalf("docker mismatch: %#v", store.saved.Docker)
	}
}

func TestCodeInitInteractor_RemoteRequiresHost(t *testing.T) {
	store := &fakeConfigStore{}
	compose := &fakeCompose{appsFile: entities.AppsFile}
	interactor := NewCodeInitInteractor(store, compose, newFakeScaffold())

	if err := interactor.Init("dev", []string{"compose.yaml"}, "", false, testRegistry, testNamespace); err == nil {
		t.Fatalf("expected error when host is empty")
	}
	if compose.buildEnv != "" {
		t.Fatalf("compose must not run when host is missing")
	}
	if store.saved.Compose.Files != nil {
		t.Fatalf("config must not be saved when host is missing")
	}
}

func TestCodeInitInteractor_RegistersDiscoveredProjectFiles(t *testing.T) {
	store := &fakeConfigStore{}
	compose := &fakeCompose{
		appsFile:   entities.AppsFile,
		discovered: []string{"compose.yaml", "compose.build.yaml"},
	}
	interactor := NewCodeInitInteractor(store, compose, newFakeScaffold())

	if err := interactor.Init("dev", []string{"compose.build.yaml", "extra.yaml"}, "", true, testRegistry, testNamespace); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantBuild := []string{"compose.yaml", "compose.build.yaml", "extra.yaml"}
	if !reflect.DeepEqual(compose.buildFiles, wantBuild) {
		t.Fatalf("build files = %#v, want %#v", compose.buildFiles, wantBuild)
	}

	want := []string{
		"../compose.yaml",
		"../compose.build.yaml",
		"../extra.yaml",
		"compose.common.yml",
		"compose.apps.yaml",
		"compose.postgres.yml",
		"compose.caddy.yml",
	}
	if !reflect.DeepEqual(store.saved.Compose.Files, want) {
		t.Fatalf("files = %#v, want %#v", store.saved.Compose.Files, want)
	}
}

func TestCodeInitInteractor_RequiresComposeFiles(t *testing.T) {
	store := &fakeConfigStore{}
	compose := &fakeCompose{appsFile: entities.AppsFile}
	interactor := NewCodeInitInteractor(store, compose, newFakeScaffold())

	if err := interactor.Init("dev", nil, "", true, testRegistry, testNamespace); err == nil {
		t.Fatalf("expected error when no compose files are found")
	}
}

func TestCodeInitInteractor_CustomNamespace(t *testing.T) {
	store := &fakeConfigStore{}
	compose := &fakeCompose{appsFile: entities.AppsFile}
	scaffold := newFakeScaffold()
	interactor := NewCodeInitInteractor(store, compose, scaffold)

	if err := interactor.Init("dev", []string{"compose.yaml"}, "", true, "", "acme"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if store.saved.Namespace != "acme" {
		t.Fatalf("namespace = %q, want acme", store.saved.Namespace)
	}
	if compose.buildNamespace != "acme" {
		t.Fatalf("BuildApps namespace = %q, want acme", compose.buildNamespace)
	}
	if scaffold.namespace != "acme" {
		t.Fatalf("WriteBase namespace = %q, want acme", scaffold.namespace)
	}
}

func TestCodeInitInteractor_RejectsInvalidNamespace(t *testing.T) {
	store := &fakeConfigStore{}
	compose := &fakeCompose{appsFile: entities.AppsFile}
	interactor := NewCodeInitInteractor(store, compose, newFakeScaffold())

	if err := interactor.Init("dev", []string{"compose.yaml"}, "", true, "", "bad namespace"); err == nil {
		t.Fatalf("expected error for invalid namespace")
	}
	if store.saved.Compose.Files != nil {
		t.Fatalf("config must not be saved for invalid namespace")
	}
}

func TestCodeInitInteractor_RegistryIsOptional(t *testing.T) {
	store := &fakeConfigStore{}
	compose := &fakeCompose{appsFile: entities.AppsFile}
	interactor := NewCodeInitInteractor(store, compose, newFakeScaffold())

	if err := interactor.Init("dev", []string{"compose.yaml"}, "", true, "", testNamespace); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if compose.buildRegistry != "" {
		t.Fatalf("BuildApps registry = %q, want empty", compose.buildRegistry)
	}
	if _, ok := store.saved.Vars["MINIPAAS_IMAGE_PREFIX"]; ok {
		t.Fatalf("image prefix must not be saved when no registry is given: %#v", store.saved.Vars)
	}
}

func TestCodeInitInteractor_Errors(t *testing.T) {
	buildErr := errors.New("build boom")
	store := &fakeConfigStore{}
	compose := &fakeCompose{addErr: buildErr}
	interactor := NewCodeInitInteractor(store, compose, newFakeScaffold())

	if err := interactor.Init("dev", []string{"compose.yaml"}, "", true, testRegistry, testNamespace); !errors.Is(err, buildErr) {
		t.Fatalf("err = %v, want %v", err, buildErr)
	}
	if store.saved.Compose.Files != nil {
		t.Fatalf("config must not be saved on compose error")
	}

	scaffoldErr := errors.New("scaffold boom")
	store = &fakeConfigStore{}
	interactor = NewCodeInitInteractor(store, &fakeCompose{appsFile: entities.AppsFile}, &fakeScaffold{err: scaffoldErr})

	if err := interactor.Init("dev", []string{"compose.yaml"}, "", true, testRegistry, testNamespace); !errors.Is(err, scaffoldErr) {
		t.Fatalf("err = %v, want %v", err, scaffoldErr)
	}
	if store.saved.Compose.Files != nil {
		t.Fatalf("config must not be saved on scaffold error")
	}
}
