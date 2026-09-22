package usecases

import (
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"

	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/entities"
)

const defaultDeployVersion = "0.1.0"

type CodeInitCase interface {
	Init(env string, files []string, host string, local bool, registry string, namespace string) error
}

type CodeInitInteractor struct {
	config   ConfigStorePort
	compose  ComposePort
	scaffold ScaffoldPort
}

func (i *CodeInitInteractor) Init(env string, files []string, host string, local bool, registry string, namespace string) error {
	if !local && host == "" {
		return errors.New("host is required when local is false")
	}

	registry = strings.TrimSuffix(strings.TrimSpace(registry), "/")
	namespace = strings.TrimSpace(namespace)
	if namespace == "" {
		namespace = entities.DefaultNamespace
	}
	if strings.IndexFunc(namespace, invalidNamespaceRune) != -1 {
		return errors.New("namespace may only contain letters, digits, '-', '_' and '.', got " + namespace)
	}

	sourceFiles := mergeSourceFiles(i.compose.DiscoverProjectFiles(env), files)
	if len(sourceFiles) == 0 {
		return errors.New("no compose files found; add compose.yaml to the project root or pass -c")
	}

	appsFile, err := i.compose.BuildApps(env, sourceFiles, registry, namespace)
	if err != nil {
		return err
	}

	base, err := i.scaffold.WriteBase(env, namespace)
	if err != nil {
		return err
	}

	relFiles, err := filesRelativeToEnv(env, sourceFiles)
	if err != nil {
		return err
	}

	composeFiles := make([]string, 0, len(relFiles)+4)
	composeFiles = append(composeFiles, relFiles...)
	composeFiles = append(composeFiles, base.Common, appsFile, base.Postgres, base.Caddy)

	vars := map[string]string{"MINIPAAS_DEPLOY_VERSION": defaultDeployVersion}
	if registry != "" {
		vars["MINIPAAS_IMAGE_PREFIX"] = registry
	}

	cfg := entities.MinipaasConfig{
		Compose:   entities.ComposeConfig{Files: composeFiles},
		Namespace: namespace,
		Vars:      vars,
	}

	if !local {
		cfg.Docker = entities.DockerConfig{Host: sshHost(host)}
	}

	_, err = i.config.Save(env, cfg)
	return err
}

func mergeSourceFiles(discovered, extra []string) []string {
	seen := make(map[string]struct{}, len(discovered)+len(extra))
	out := make([]string, 0, len(discovered)+len(extra))

	for _, f := range append(append([]string{}, discovered...), extra...) {
		key, err := filepath.Abs(f)
		if err != nil {
			key = f
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, f)
	}

	return out
}

func invalidNamespaceRune(r rune) bool {
	switch {
	case r >= 'a' && r <= 'z',
		r >= 'A' && r <= 'Z',
		r >= '0' && r <= '9',
		r == '-', r == '_', r == '.':
		return false
	}
	return true
}

func sshHost(host string) string {
	if host == "" || strings.Contains(host, "://") {
		return host
	}
	return "ssh://" + host
}

func filesRelativeToEnv(env string, files []string) ([]string, error) {
	envAbs, err := filepath.Abs(env)
	if err != nil {
		return nil, err
	}

	rel := make([]string, 0, len(files))
	for _, f := range files {
		if filepath.IsAbs(f) {
			rel = append(rel, f)
			continue
		}

		fileAbs, err := filepath.Abs(f)
		if err != nil {
			return nil, err
		}

		path, err := filepath.Rel(envAbs, fileAbs)
		if err != nil {
			return nil, err
		}
		rel = append(rel, filepath.ToSlash(path))
	}

	return rel, nil
}

func NewCodeInitInteractor(config ConfigStorePort, compose ComposePort, scaffold ScaffoldPort) *CodeInitInteractor {
	return &CodeInitInteractor{config: config, compose: compose, scaffold: scaffold}
}

var _ CodeInitCase = (*CodeInitInteractor)(nil)
