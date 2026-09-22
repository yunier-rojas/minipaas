package usecases

import (
	"testing"

	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/entities"
)

func TestCodeJobInteractor_Configure(t *testing.T) {
	compose := &fakeCompose{appsFile: entities.AppsFile}
	env := &fakeEnvironment{}
	interactor := NewCodeJobInteractor(&fakeConfigStore{}, env, compose)

	file, err := interactor.Configure("dev", []string{"migrate"}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if file != entities.AppsFile {
		t.Fatalf("file = %q", file)
	}
	if compose.jobEnv != "dev" || len(compose.jobServices) != 1 || compose.jobServices[0] != "migrate" {
		t.Fatalf("job args mismatch: %#v", compose)
	}
	if !env.applied {
		t.Fatalf("environment not applied")
	}
}
