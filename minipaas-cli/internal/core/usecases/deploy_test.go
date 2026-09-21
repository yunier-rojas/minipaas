package usecases

import (
	"errors"
	"reflect"
	"testing"

	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/entities"
)

func TestDeployBuildInteractor(t *testing.T) {
	env := &fakeEnvironment{}
	compose := &fakeCompose{files: []string{"dev/compose.apps.yaml"}}
	deploy := &fakeDeploy{}
	store := &fakeConfigStore{cfg: entities.MinipaasConfig{Vars: map[string]string{"MINIPAAS_IMAGE_PREFIX": "ghcr.io/org"}}}
	interactor := NewDeployBuildInteractor(store, env, compose, deploy)

	if err := interactor.Build("dev", true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !env.applied {
		t.Fatalf("environment not applied")
	}
	if !reflect.DeepEqual(deploy.buildFiles, compose.files) {
		t.Fatalf("files = %#v, want %#v", deploy.buildFiles, compose.files)
	}
	if deploy.buildRegistry != "ghcr.io/org" {
		t.Fatalf("registry = %q, want ghcr.io/org", deploy.buildRegistry)
	}
	if !deploy.verbose {
		t.Fatalf("verbose not forwarded")
	}
}

func TestDeployRolloutInteractor(t *testing.T) {
	env := &fakeEnvironment{}
	compose := &fakeCompose{files: []string{"dev/compose.apps.yaml"}}
	deploy := &fakeDeploy{}
	interactor := NewDeployRolloutInteractor(&fakeConfigStore{}, env, compose, deploy)

	if err := interactor.Rollout("dev", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !env.applied {
		t.Fatalf("environment not applied")
	}
	if !reflect.DeepEqual(deploy.rolloutFiles, compose.files) {
		t.Fatalf("files = %#v, want %#v", deploy.rolloutFiles, compose.files)
	}
}

func TestDeployCanaryInteractor(t *testing.T) {
	env := &fakeEnvironment{}
	compose := &fakeCompose{files: []string{"dev/compose.apps.yaml"}}
	deploy := &fakeDeploy{}
	interactor := NewDeployCanaryInteractor(&fakeConfigStore{}, env, compose, deploy)

	services := []string{"api", "worker"}
	if err := interactor.Canary("dev", services, 3, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(deploy.canaryServices, services) {
		t.Fatalf("services = %#v, want %#v", deploy.canaryServices, services)
	}
	if deploy.canaryReplicas != 3 {
		t.Fatalf("replicas = %d, want 3", deploy.canaryReplicas)
	}
	if !reflect.DeepEqual(deploy.canaryFiles, compose.files) {
		t.Fatalf("files = %#v, want %#v", deploy.canaryFiles, compose.files)
	}
}

func TestDeployRoutingInteractor(t *testing.T) {
	env := &fakeEnvironment{}
	routing := &fakeRouting{}
	interactor := NewDeployRoutingInteractor(&fakeConfigStore{}, env, routing)

	if err := interactor.Apply("dev", true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !env.applied {
		t.Fatalf("environment not applied")
	}
	if routing.env != "dev" || !routing.verbose {
		t.Fatalf("routing args mismatch: %#v", routing)
	}
}

func TestDeployInteractors_ConfigError(t *testing.T) {
	loadErr := errors.New("load boom")
	store := &fakeConfigStore{err: loadErr}
	deploy := &fakeDeploy{}
	routing := &fakeRouting{}

	if err := NewDeployBuildInteractor(store, &fakeEnvironment{}, &fakeCompose{}, deploy).Build("dev", false); !errors.Is(err, loadErr) {
		t.Fatalf("build err = %v, want %v", err, loadErr)
	}
	if err := NewDeployRolloutInteractor(store, &fakeEnvironment{}, &fakeCompose{}, deploy).Rollout("dev", false); !errors.Is(err, loadErr) {
		t.Fatalf("rollout err = %v, want %v", err, loadErr)
	}
	if err := NewDeployCanaryInteractor(store, &fakeEnvironment{}, &fakeCompose{}, deploy).Canary("dev", []string{"api"}, 1, false); !errors.Is(err, loadErr) {
		t.Fatalf("canary err = %v, want %v", err, loadErr)
	}
	if err := NewDeployRoutingInteractor(store, &fakeEnvironment{}, routing).Apply("dev", false); !errors.Is(err, loadErr) {
		t.Fatalf("routing err = %v, want %v", err, loadErr)
	}

	if deploy.buildFiles != nil || deploy.rolloutFiles != nil || deploy.canaryFiles != nil || routing.env != "" {
		t.Fatalf("ports must not be called on config error")
	}
}
