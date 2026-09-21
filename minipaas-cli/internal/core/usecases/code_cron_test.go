package usecases

import (
	"testing"

	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/entities"
)

func TestCodeCronInteractor_Configure(t *testing.T) {
	compose := &fakeCompose{appsFile: entities.AppsFile}
	env := &fakeEnvironment{}
	interactor := NewCodeCronInteractor(&fakeConfigStore{}, env, compose)

	file, err := interactor.Configure("dev", []string{"cleanup"}, "*/5 * * * *", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if file != entities.AppsFile {
		t.Fatalf("file = %q", file)
	}
	if compose.cronEnv != "dev" || compose.cronSchedule != "*/5 * * * *" {
		t.Fatalf("cron args mismatch: %#v", compose)
	}
	if len(compose.cronServices) != 1 || compose.cronServices[0] != "cleanup" {
		t.Fatalf("cron services mismatch: %#v", compose.cronServices)
	}
	if !env.applied {
		t.Fatalf("environment not applied")
	}
}
