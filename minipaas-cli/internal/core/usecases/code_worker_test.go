package usecases

import (
	"testing"

	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/entities"
)

func TestCodeWorkerInteractor_Configure(t *testing.T) {
	compose := &fakeCompose{appsFile: entities.AppsFile}
	env := &fakeEnvironment{}
	interactor := NewCodeWorkerInteractor(&fakeConfigStore{}, env, compose)

	file, err := interactor.Configure("dev", []string{"api", "worker"}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if file != entities.AppsFile {
		t.Fatalf("file = %q", file)
	}
	if compose.workerEnv != "dev" || len(compose.workerServices) != 2 {
		t.Fatalf("worker args mismatch: %#v", compose)
	}
	if !env.applied {
		t.Fatalf("environment not applied")
	}
}
