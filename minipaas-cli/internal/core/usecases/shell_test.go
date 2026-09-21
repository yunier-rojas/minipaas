package usecases

import "testing"

func TestShellInteractor_AppliesEnvironmentAndLaunches(t *testing.T) {
	env := &fakeEnvironment{}
	shell := &fakeShell{}
	interactor := NewShellInteractor(&fakeConfigStore{}, env, shell)

	if err := interactor.Launch("/env", true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !env.applied {
		t.Fatalf("environment not applied")
	}
	if !shell.launched {
		t.Fatalf("shell not launched")
	}
}
