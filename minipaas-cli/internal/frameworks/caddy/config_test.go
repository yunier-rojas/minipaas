package caddy

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

type testCaddyConfig struct {
	Apps struct {
		HTTP struct {
			Servers struct {
				Minipaas struct {
					Listen         []string `json:"listen"`
					AutomaticHTTPS struct {
						Disable          bool `json:"disable"`
						DisableRedirects bool `json:"disable_redirects"`
					} `json:"automatic_https"`
					Routes []struct {
						Match []struct {
							Host []string `json:"host"`
							Path []string `json:"path"`
						} `json:"match"`
						Handle []struct {
							Upstreams []struct {
								Dial string `json:"dial"`
							} `json:"upstreams"`
						} `json:"handle"`
					} `json:"routes"`
				} `json:"minipaas"`
			} `json:"servers"`
		} `json:"http"`
	} `json:"apps"`
}

func writeEmptyCaddy(t *testing.T, dir string) string {
	t.Helper()
	fn := filepath.Join(dir, configFile)
	if err := os.WriteFile(fn, []byte("{}"), 0644); err != nil {
		t.Fatalf("write caddy.json: %v", err)
	}
	return fn
}

func readTestCaddy(t *testing.T, fn string) testCaddyConfig {
	t.Helper()
	data, err := os.ReadFile(fn)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	var cfg testCaddyConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return cfg
}

func TestAddRoute_HTTPS(t *testing.T) {
	dir := t.TempDir()
	fn := writeEmptyCaddy(t, dir)

	out, err := NewConfigService().AddRoute(dir, "example.com/app", "api:8080", "minipaas")
	if err != nil {
		t.Fatalf("AddRoute: %v", err)
	}
	if out != fn {
		t.Fatalf("file = %q, want %q", out, fn)
	}

	srv := readTestCaddy(t, fn).Apps.HTTP.Servers.Minipaas
	if len(srv.Listen) != 1 || srv.Listen[0] != ":443" {
		t.Fatalf("listen = %#v, want [:443]", srv.Listen)
	}
	if srv.AutomaticHTTPS.Disable || srv.AutomaticHTTPS.DisableRedirects {
		t.Fatalf("automatic https should stay enabled for https")
	}
	if len(srv.Routes) != 1 {
		t.Fatalf("routes = %d, want 1", len(srv.Routes))
	}
	r := srv.Routes[0]
	if r.Match[0].Host[0] != "example.com" || r.Match[0].Path[0] != "/app/*" {
		t.Fatalf("match mismatch: %#v", r.Match)
	}
	if r.Handle[0].Upstreams[0].Dial != "minipaas_api:8080" {
		t.Fatalf("upstream mismatch: %#v", r.Handle)
	}
}

func TestAddRoute_HTTP_CustomPort(t *testing.T) {
	dir := t.TempDir()
	fn := writeEmptyCaddy(t, dir)

	if _, err := NewConfigService().AddRoute(dir, "http://example.com:8081/root", "web:80", "minipaas"); err != nil {
		t.Fatalf("AddRoute: %v", err)
	}
	srv := readTestCaddy(t, fn).Apps.HTTP.Servers.Minipaas
	if len(srv.Listen) != 1 || srv.Listen[0] != ":8081" {
		t.Fatalf("listen = %#v, want [:8081]", srv.Listen)
	}
	if !srv.AutomaticHTTPS.Disable || !srv.AutomaticHTTPS.DisableRedirects {
		t.Fatalf("automatic https should be disabled for http")
	}
	if srv.Routes[0].Match[0].Host[0] != "example.com" || srv.Routes[0].Match[0].Path[0] != "/root/*" {
		t.Fatalf("match mismatch: %#v", srv.Routes[0].Match)
	}
}

func TestAddRoute_PathOnly_DefaultsToHTTP80(t *testing.T) {
	dir := t.TempDir()
	fn := writeEmptyCaddy(t, dir)

	if _, err := NewConfigService().AddRoute(dir, "/web", "web:80", "minipaas"); err != nil {
		t.Fatalf("AddRoute: %v", err)
	}

	srv := readTestCaddy(t, fn).Apps.HTTP.Servers.Minipaas
	if len(srv.Listen) != 1 || srv.Listen[0] != ":80" {
		t.Fatalf("listen = %#v, want [:80]", srv.Listen)
	}
	if !srv.AutomaticHTTPS.Disable || !srv.AutomaticHTTPS.DisableRedirects {
		t.Fatalf("automatic https should be disabled for hostless url")
	}
	if len(srv.Routes) != 1 {
		t.Fatalf("routes = %d, want 1", len(srv.Routes))
	}
	r := srv.Routes[0]
	if len(r.Match[0].Host) != 0 {
		t.Fatalf("host constraint should be omitted: %#v", r.Match)
	}
	if r.Match[0].Path[0] != "/web/*" {
		t.Fatalf("path mismatch: %#v", r.Match)
	}
	if r.Handle[0].Upstreams[0].Dial != "minipaas_web:80" {
		t.Fatalf("upstream mismatch: %#v", r.Handle)
	}
}

func TestAddRoute_HostlessWithPort(t *testing.T) {
	dir := t.TempDir()
	fn := writeEmptyCaddy(t, dir)

	if _, err := NewConfigService().AddRoute(dir, "http://:8000/web", "web:80", "minipaas"); err != nil {
		t.Fatalf("AddRoute: %v", err)
	}

	srv := readTestCaddy(t, fn).Apps.HTTP.Servers.Minipaas
	if len(srv.Listen) != 1 || srv.Listen[0] != ":8000" {
		t.Fatalf("listen = %#v, want [:8000]", srv.Listen)
	}
	r := srv.Routes[0]
	if len(r.Match[0].Host) != 0 || r.Match[0].Path[0] != "/web/*" {
		t.Fatalf("match mismatch: %#v", r.Match)
	}
}

func TestAddRoute_ReplaceAndRootPath(t *testing.T) {
	dir := t.TempDir()
	fn := writeEmptyCaddy(t, dir)
	svc := NewConfigService()

	if _, err := svc.AddRoute(dir, "https://example.org/", "web:80", "minipaas"); err != nil {
		t.Fatalf("first add: %v", err)
	}
	srv := readTestCaddy(t, fn).Apps.HTTP.Servers.Minipaas
	if len(srv.Routes) != 1 || srv.Routes[0].Match[0].Path[0] != "/*" {
		t.Fatalf("root path normalization failed: %#v", srv.Routes)
	}

	if _, err := svc.AddRoute(dir, "https://example.org/", "web:81", "minipaas"); err != nil {
		t.Fatalf("replace: %v", err)
	}
	srv = readTestCaddy(t, fn).Apps.HTTP.Servers.Minipaas
	if len(srv.Routes) != 1 {
		t.Fatalf("route should be replaced, not duplicated: %d", len(srv.Routes))
	}
	if srv.Routes[0].Handle[0].Upstreams[0].Dial != "minipaas_web:81" {
		t.Fatalf("upstream not updated: %#v", srv.Routes[0].Handle)
	}
}

func TestAddRoute_SecondRouteAndListenDedup(t *testing.T) {
	dir := t.TempDir()
	fn := writeEmptyCaddy(t, dir)
	svc := NewConfigService()

	if _, err := svc.AddRoute(dir, "https://example.org/a", "web:80", "minipaas"); err != nil {
		t.Fatalf("first: %v", err)
	}
	if _, err := svc.AddRoute(dir, "https://example.org/b", "api:80", "minipaas"); err != nil {
		t.Fatalf("second: %v", err)
	}

	srv := readTestCaddy(t, fn).Apps.HTTP.Servers.Minipaas
	if len(srv.Routes) != 2 {
		t.Fatalf("routes = %d, want 2", len(srv.Routes))
	}
	count := 0
	for _, l := range srv.Listen {
		if l == ":443" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf(":443 should appear once, got %d", count)
	}
}

func TestAddRoute_MissingConfig(t *testing.T) {
	if _, err := NewConfigService().AddRoute(t.TempDir(), "example.com", "api", "minipaas"); err == nil {
		t.Fatalf("expected error for missing caddy.json")
	}
}

func TestAddRoute_CustomNamespace(t *testing.T) {
	dir := t.TempDir()
	fn := writeEmptyCaddy(t, dir)

	if _, err := NewConfigService().AddRoute(dir, "https://example.org/app", "web:80", "acme"); err != nil {
		t.Fatalf("AddRoute: %v", err)
	}

	srv := readTestCaddy(t, fn).Apps.HTTP.Servers.Minipaas
	if srv.Routes[0].Handle[0].Upstreams[0].Dial != "acme_web:80" {
		t.Fatalf("upstream mismatch: %#v", srv.Routes[0].Handle)
	}
}

func TestParsePublicURL(t *testing.T) {
	tests := []struct {
		in   string
		host string
		port string
		path string
	}{
		{"example.com", "example.com", "443", ""},
		{"http://example.com", "example.com", "80", ""},
		{"https://example.com", "example.com", "443", ""},
		{"https://example.com/app", "example.com", "443", "/app"},
		{"http://example.com:8080/app", "example.com", "8080", "/app"},
		{"/web", "", "", "/web"},
	}
	for _, tt := range tests {
		u, err := parsePublicURL(tt.in)
		if err != nil {
			t.Fatalf("parsePublicURL(%q): %v", tt.in, err)
		}
		if u.Hostname() != tt.host || u.Port() != tt.port || u.Path != tt.path {
			t.Fatalf("parsePublicURL(%q) = %q,%q,%q", tt.in, u.Hostname(), u.Port(), u.Path)
		}
	}
}

func TestNormalizeCaddyPath(t *testing.T) {
	cases := map[string]string{
		"":       "/*",
		"/":      "/*",
		"api":    "/api/*",
		"/api":   "/api/*",
		"/api/":  "/api/*",
		"/api/*": "/api/*",
	}
	for in, want := range cases {
		if got := normalizeCaddyPath(in); got != want {
			t.Fatalf("normalizeCaddyPath(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestSplitTarget(t *testing.T) {
	if s, p := splitTarget("web"); s != "web" || p != "80" {
		t.Fatalf("default port failed: %s %s", s, p)
	}
	if s, p := splitTarget("web:8080"); s != "web" || p != "8080" {
		t.Fatalf("explicit port failed: %s %s", s, p)
	}
}
