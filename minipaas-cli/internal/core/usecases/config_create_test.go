package usecases

import (
	"testing"

	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/entities"
)

func TestConfigCreateInteractor_CreateAndPatch(t *testing.T) {
	content := []byte("hello")
	input := &fakeInput{baseName: "app.conf", content: content}
	env := &fakeEnvironment{}
	configs := &fakeConfigStoreDocker{}
	compose := &fakeCompose{files: []string{"a.yml"}, grouped: map[string][]string{"a.yml": {"api"}}}
	store := &fakeConfigStore{cfg: entities.MinipaasConfig{Compose: entities.ComposeConfig{Files: []string{"a.yml"}}}}
	interactor := NewConfigCreateInteractor(store, env, input, configs, compose)

	result, err := interactor.CreateConfig("envdir", "", "", []string{"api"}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := configHashName("app.conf", content)
	if result.Name != want {
		t.Fatalf("name = %q, want %q", result.Name, want)
	}
	if !result.Compose || !result.Attached {
		t.Fatalf("unexpected result: %#v", result)
	}
	if len(configs.created) != 1 || configs.created[0] != want {
		t.Fatalf("config not created: %#v", configs.created)
	}
	if len(compose.addedConfigs) != 1 || compose.addedConfigs[0] != "a.yml" {
		t.Fatalf("compose not patched: %#v", compose.addedConfigs)
	}
	if !env.applied {
		t.Fatalf("environment not applied")
	}
}

func TestConfigCreateInteractor_ExistingConfigIsNotRecreated(t *testing.T) {
	input := &fakeInput{baseName: "app.conf", content: []byte("x")}
	configs := &fakeConfigStoreDocker{exists: true}
	store := &fakeConfigStore{cfg: entities.MinipaasConfig{Compose: entities.ComposeConfig{Files: []string{"a.yml"}}}}
	interactor := NewConfigCreateInteractor(store, &fakeEnvironment{}, input, configs, &fakeCompose{})

	if _, err := interactor.CreateConfig("envdir", "", "", nil, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(configs.created) != 0 {
		t.Fatalf("config should not be created: %#v", configs.created)
	}
}

func TestConfigCreateInteractor_EmptyInputSkips(t *testing.T) {
	input := &fakeInput{baseName: "app.conf", content: nil}
	configs := &fakeConfigStoreDocker{}
	compose := &fakeCompose{}
	interactor := NewConfigCreateInteractor(&fakeConfigStore{}, &fakeEnvironment{}, input, configs, compose)

	result, err := interactor.CreateConfig("envdir", "", "", nil, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "" {
		t.Fatalf("name should be empty, got %q", result.Name)
	}
	if len(configs.created) != 0 || len(compose.addedConfigs) != 0 {
		t.Fatalf("nothing should happen on empty input")
	}
}

func TestConfigCreateInteractor_MissingServiceDoesNotCreate(t *testing.T) {
	input := &fakeInput{baseName: "app.conf", content: []byte("x")}
	compose := &fakeCompose{missing: []string{"worker"}}
	configs := &fakeConfigStoreDocker{}
	store := &fakeConfigStore{cfg: entities.MinipaasConfig{Compose: entities.ComposeConfig{Files: []string{"a.yml"}}}}
	interactor := NewConfigCreateInteractor(store, &fakeEnvironment{}, input, configs, compose)

	if _, err := interactor.CreateConfig("envdir", "", "", []string{"worker"}, false); err == nil {
		t.Fatalf("expected error for missing service")
	}
	if len(configs.created) != 0 {
		t.Fatalf("config must not be created before validation: %#v", configs.created)
	}
}
