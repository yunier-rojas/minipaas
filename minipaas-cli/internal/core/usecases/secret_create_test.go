package usecases

import (
	"testing"

	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/entities"
)

func TestSecretCreateInteractor_CreateAndPatch(t *testing.T) {
	content := []byte("supersecret")
	input := &fakeInput{baseName: "env", content: content}
	env := &fakeEnvironment{}
	secrets := &fakeSecretStore{}
	compose := &fakeCompose{files: []string{"a.yml"}, grouped: map[string][]string{"a.yml": {"api"}}}
	store := &fakeConfigStore{cfg: entities.MinipaasConfig{Compose: entities.ComposeConfig{Files: []string{"a.yml"}}}}
	interactor := NewSecretCreateInteractor(store, env, input, secrets, compose)

	result, err := interactor.CreateSecret("envdir", "", "", []string{"api"}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := secretHashName("env", content)
	if result.Name != want {
		t.Fatalf("name = %q, want %q", result.Name, want)
	}
	if !result.Compose || !result.Attached {
		t.Fatalf("unexpected result: %#v", result)
	}
	if len(secrets.created) != 1 || secrets.created[0] != want {
		t.Fatalf("secret not created: %#v", secrets.created)
	}
	if len(compose.addedSecrets) != 1 || compose.addedSecrets[0] != "a.yml" {
		t.Fatalf("compose not patched: %#v", compose.addedSecrets)
	}
	if !env.applied {
		t.Fatalf("environment not applied")
	}
}

func TestSecretCreateInteractor_ExistingSecretIsNotRecreated(t *testing.T) {
	input := &fakeInput{baseName: "env", content: []byte("x")}
	secrets := &fakeSecretStore{exists: true}
	store := &fakeConfigStore{cfg: entities.MinipaasConfig{Compose: entities.ComposeConfig{Files: []string{"a.yml"}}}}
	interactor := NewSecretCreateInteractor(store, &fakeEnvironment{}, input, secrets, &fakeCompose{})

	if _, err := interactor.CreateSecret("envdir", "", "", nil, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(secrets.created) != 0 {
		t.Fatalf("secret should not be created: %#v", secrets.created)
	}
}

func TestSecretCreateInteractor_EmptyInputSkips(t *testing.T) {
	input := &fakeInput{baseName: "env", content: nil}
	secrets := &fakeSecretStore{}
	compose := &fakeCompose{}
	interactor := NewSecretCreateInteractor(&fakeConfigStore{}, &fakeEnvironment{}, input, secrets, compose)

	result, err := interactor.CreateSecret("envdir", "", "", nil, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "" {
		t.Fatalf("name should be empty, got %q", result.Name)
	}
	if len(secrets.created) != 0 || len(compose.addedSecrets) != 0 {
		t.Fatalf("nothing should happen on empty input")
	}
}

func TestSecretCreateInteractor_MissingServiceDoesNotCreate(t *testing.T) {
	input := &fakeInput{baseName: "env", content: []byte("x")}
	compose := &fakeCompose{missing: []string{"worker"}}
	secrets := &fakeSecretStore{}
	store := &fakeConfigStore{cfg: entities.MinipaasConfig{Compose: entities.ComposeConfig{Files: []string{"a.yml"}}}}
	interactor := NewSecretCreateInteractor(store, &fakeEnvironment{}, input, secrets, compose)

	if _, err := interactor.CreateSecret("envdir", "", "", []string{"worker"}, false); err == nil {
		t.Fatalf("expected error for missing service")
	}
	if len(secrets.created) != 0 {
		t.Fatalf("secret must not be created before validation: %#v", secrets.created)
	}
}

func TestSecretCreateInteractor_NoComposeFilesSkipsAttach(t *testing.T) {
	input := &fakeInput{baseName: "postgres_password", content: []byte("x")}
	secrets := &fakeSecretStore{}
	compose := &fakeCompose{}
	interactor := NewSecretCreateInteractor(&fakeConfigStore{}, &fakeEnvironment{}, input, secrets, compose)

	result, err := interactor.CreateSecret("envdir", "", "", []string{"postgres"}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Compose || result.Attached {
		t.Fatalf("result = %#v, want unattached", result)
	}
	if len(secrets.created) != 1 {
		t.Fatalf("secret not created: %#v", secrets.created)
	}
	if len(compose.addedSecrets) != 0 {
		t.Fatalf("compose must not be patched: %#v", compose.addedSecrets)
	}
}
