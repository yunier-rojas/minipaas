package usecases

import (
	"errors"
	"testing"
)

func TestCodeRouteInteractor_AddRoute(t *testing.T) {
	env := &fakeEnvironment{}
	routes := &fakeRoute{file: "dev/caddy.json"}
	compose := &fakeCompose{appsFile: "dev/compose.apps.yaml"}
	interactor := NewCodeRouteInteractor(&fakeConfigStore{}, env, routes, compose)

	routeFile, appsFile, err := interactor.AddRoute("dev", "https://example.com/app", "api:8080", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if routeFile != "dev/caddy.json" {
		t.Fatalf("routeFile = %q, want dev/caddy.json", routeFile)
	}
	if appsFile != "dev/compose.apps.yaml" {
		t.Fatalf("appsFile = %q, want dev/compose.apps.yaml", appsFile)
	}
	if !env.applied {
		t.Fatalf("environment not applied")
	}
	if !routes.called {
		t.Fatalf("route port not called")
	}
	if routes.env != "dev" || routes.url != "https://example.com/app" || routes.target != "api:8080" {
		t.Fatalf("route args mismatch: %#v", routes)
	}
	if compose.routeEnv != "dev" || compose.routeService != "api" || compose.routePort != "8080" {
		t.Fatalf("route deploy args mismatch: %#v", compose)
	}
}

func TestCodeRouteInteractor_ConfigError(t *testing.T) {
	loadErr := errors.New("load boom")
	routes := &fakeRoute{}
	compose := &fakeCompose{}
	interactor := NewCodeRouteInteractor(&fakeConfigStore{err: loadErr}, &fakeEnvironment{}, routes, compose)

	if _, _, err := interactor.AddRoute("dev", "example.com", "api", false); !errors.Is(err, loadErr) {
		t.Fatalf("err = %v, want %v", err, loadErr)
	}
	if routes.called {
		t.Fatalf("route port must not be called on config error")
	}
}
