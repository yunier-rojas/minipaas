package caddy

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/entities"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/usecases"
)

type ConfigService struct{}

func NewConfigService() *ConfigService {
	return &ConfigService{}
}

func (s *ConfigService) AddRoute(env, publicURL, target, namespace string) (string, error) {
	fn := filepath.Join(env, configFile)

	parsed, err := parsePublicURL(publicURL)
	if err != nil {
		return fn, err
	}

	domain := parsed.Hostname()
	scheme := parsed.Scheme
	publicPort := parsed.Port()
	if domain == "" {
		scheme = "http"
		if publicPort == "" {
			publicPort = "80"
		}
	} else if scheme == "" {
		scheme = "https"
	}
	normPath := normalizeCaddyPath(parsed.Path)

	service, svcPort := splitTarget(target)
	upstreamDial := fmt.Sprintf("%s_%s:%s", namespace, service, svcPort)

	data, err := os.ReadFile(fn)
	if err != nil {
		return fn, err
	}

	var root map[string]interface{}
	if err = json.Unmarshal(data, &root); err != nil {
		return fn, err
	}

	server := ensureServerAndListen(root, scheme, publicPort)

	hostOnly := domain
	if h, _, err := net.SplitHostPort(domain); err == nil {
		hostOnly = h
	}
	newRoute := buildRoute(hostOnly, normPath, upstreamDial)

	replaceOrAppendRoute(server, newRoute)

	out, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return fn, err
	}
	if err = os.WriteFile(fn, out, 0644); err != nil {
		return fn, err
	}
	return fn, nil
}

func parsePublicURL(input string) (*url.URL, error) {
	if strings.HasPrefix(input, "/") {
		return url.Parse(input)
	}
	if !strings.Contains(input, "://") {
		input = "https://" + input
	}
	parsed, err := url.Parse(input)
	if err != nil {
		return nil, err
	}
	if parsed.Port() == "" {
		if parsed.Scheme == "http" {
			parsed.Host = parsed.Host + ":80"
		} else {
			parsed.Host = parsed.Host + ":443"
		}
	}
	return parsed, nil
}

func normalizeCaddyPath(p string) string {
	if p == "" || p == "/" {
		return "/*"
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	if !strings.HasSuffix(p, "*") {
		if strings.HasSuffix(p, "/") {
			p = p + "*"
		} else {
			p = p + "/*"
		}
	}
	return p
}

func splitTarget(target string) (service, port string) {
	return entities.ParseServiceTarget(target)
}

func ensureServerAndListen(root map[string]interface{}, scheme, publicPort string) map[string]interface{} {
	ensureMap := func(parent map[string]interface{}, key string) map[string]interface{} {
		if parent == nil {
			return nil
		}
		if v, ok := parent[key]; ok {
			if m, ok := v.(map[string]interface{}); ok {
				return m
			}
		}
		m := map[string]interface{}{}
		parent[key] = m
		return m
	}

	apps := ensureMap(root, "apps")
	httpApp := ensureMap(apps, "http")
	servers := ensureMap(httpApp, "servers")

	server, _ := servers["minipaas"].(map[string]interface{})
	if server == nil {
		server = map[string]interface{}{}
		servers["minipaas"] = server
	}

	if publicPort != "" {
		listenAddr := ":" + publicPort
		var listenSlice []interface{}
		if v, ok := server["listen"].([]interface{}); ok {
			listenSlice = v
		}
		hasListen := false
		for _, l := range listenSlice {
			if s, ok := l.(string); ok && s == listenAddr {
				hasListen = true
				break
			}
		}
		if !hasListen {
			listenSlice = append(listenSlice, listenAddr)
			server["listen"] = listenSlice
		}
	}

	if scheme == "http" {
		aut, _ := server["automatic_https"].(map[string]interface{})
		if aut == nil {
			aut = map[string]interface{}{}
		}
		aut["disable"] = true
		aut["disable_redirects"] = true
		server["automatic_https"] = aut
	}

	return server
}

func buildRoute(host, path, upstreamDial string) map[string]interface{} {
	match := map[string]interface{}{
		"path": []interface{}{path},
	}
	if host != "" {
		match["host"] = []interface{}{host}
	}
	return map[string]interface{}{
		"match": []interface{}{match},
		"handle": []interface{}{
			map[string]interface{}{
				"handler": "reverse_proxy",
				"upstreams": []interface{}{
					map[string]interface{}{
						"dial": upstreamDial,
					},
				},
			},
		},
		"terminal": true,
	}
}

func routeMatch(route map[string]interface{}) (host, path string) {
	v, ok := route["match"].([]interface{})
	if !ok || len(v) == 0 {
		return "", ""
	}
	m, ok := v[0].(map[string]interface{})
	if !ok {
		return "", ""
	}
	if hs, ok := m["host"].([]interface{}); ok && len(hs) > 0 {
		if s, ok := hs[0].(string); ok {
			host = s
		}
	}
	if ps, ok := m["path"].([]interface{}); ok && len(ps) > 0 {
		if s, ok := ps[0].(string); ok {
			path = s
		}
	}
	return
}

func replaceOrAppendRoute(server map[string]interface{}, newRoute map[string]interface{}) {
	var routes []interface{}
	if v, ok := server["routes"].([]interface{}); ok {
		routes = v
	}

	replaced := false
	nh, np := routeMatch(newRoute)
	for i := range routes {
		if r, ok := routes[i].(map[string]interface{}); ok {
			rh, rp := routeMatch(r)
			if rh == nh && rp == np {
				routes[i] = newRoute
				replaced = true
				break
			}
		}
	}
	if !replaced {
		routes = append(routes, newRoute)
	}
	server["routes"] = routes
}

var _ usecases.RoutePort = (*ConfigService)(nil)
