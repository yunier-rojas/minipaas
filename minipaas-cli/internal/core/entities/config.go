package entities

const (
	ConfigFile       = "minipaas.yaml"
	AppsFile         = "compose.apps.yaml"
	DefaultNamespace = "minipaas"
)

type ComposeConfig struct {
	Files []string `yaml:"files"`
}

type DockerConfig struct {
	Host string `yaml:"host,omitempty"`
}

type MinipaasConfig struct {
	Compose   ComposeConfig     `yaml:"compose"`
	Docker    DockerConfig      `yaml:"docker,omitempty"`
	Namespace string            `yaml:"namespace,omitempty"`
	Vars      map[string]string `yaml:"vars,omitempty"`
}

func (c MinipaasConfig) DeployNamespace() string {
	if c.Namespace == "" {
		return DefaultNamespace
	}
	return c.Namespace
}

type BaseFiles struct {
	Common    string
	Postgres  string
	Caddy     string
	CaddyJSON string
}
