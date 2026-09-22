package usecases

import "github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/entities"

type fakeConfigStore struct {
	cfg    entities.MinipaasConfig
	err    error
	saved  entities.MinipaasConfig
	saveAt string
}

func (f *fakeConfigStore) Load(env string) (entities.MinipaasConfig, error) {
	return f.cfg, f.err
}

func (f *fakeConfigStore) Save(env string, cfg entities.MinipaasConfig) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	f.saved = cfg
	f.saveAt = env
	return env + "/" + entities.ConfigFile, nil
}

type fakeScaffold struct {
	base      entities.BaseFiles
	namespace string
	err       error
}

func (f *fakeScaffold) WriteBase(env string, namespace string) (entities.BaseFiles, error) {
	f.namespace = namespace
	return f.base, f.err
}

type fakeEnvironment struct {
	applied bool
	env     string
	verbose bool
}

func (f *fakeEnvironment) Apply(env string, cfg entities.MinipaasConfig, verbose bool) {
	f.applied = true
	f.env = env
	f.verbose = verbose
}

type fakeInput struct {
	baseName string
	content  []byte
	err      error
}

func (f *fakeInput) ReadContent(file, name string) (string, []byte, error) {
	return f.baseName, f.content, f.err
}

type fakeSecretStore struct {
	exists    bool
	created   []string
	createErr error
}

func (f *fakeSecretStore) Exists(name string) bool {
	return f.exists
}

func (f *fakeSecretStore) Create(name string, content []byte, verbose bool) error {
	if f.createErr != nil {
		return f.createErr
	}
	f.created = append(f.created, name)
	return nil
}

type fakeConfigStoreDocker struct {
	exists    bool
	created   []string
	createErr error
}

func (f *fakeConfigStoreDocker) Exists(name string) bool {
	return f.exists
}

func (f *fakeConfigStoreDocker) Create(name string, content []byte, verbose bool) error {
	if f.createErr != nil {
		return f.createErr
	}
	f.created = append(f.created, name)
	return nil
}

type fakeCompose struct {
	files          []string
	discovered     []string
	grouped        map[string][]string
	missing        []string
	addedSecrets   []string
	addedConfigs   []string
	appsFile       string
	buildEnv       string
	buildFiles     []string
	buildRegistry  string
	buildNamespace string
	workerEnv      string
	workerServices []string
	jobEnv         string
	jobServices    []string
	cronEnv        string
	cronServices   []string
	cronSchedule   string
	routeEnv       string
	routeService   string
	routePort      string
	addErr         error
}

func (f *fakeCompose) BuildApps(env string, files []string, registry string, namespace string) (string, error) {
	if f.addErr != nil {
		return "", f.addErr
	}
	f.buildEnv = env
	f.buildFiles = files
	f.buildRegistry = registry
	f.buildNamespace = namespace
	return f.appsFile, nil
}

func (f *fakeCompose) DiscoverProjectFiles(env string) []string {
	return f.discovered
}

func (f *fakeCompose) FilesForEnv(env string, cfg entities.MinipaasConfig) []string {
	return f.files
}

func (f *fakeCompose) AddWorkerDeploy(env string, services []string, namespace string) (string, error) {
	if f.addErr != nil {
		return "", f.addErr
	}
	f.workerEnv = env
	f.workerServices = services
	return f.appsFile, nil
}

func (f *fakeCompose) AddJobDeploy(env string, services []string, namespace string) (string, error) {
	if f.addErr != nil {
		return "", f.addErr
	}
	f.jobEnv = env
	f.jobServices = services
	return f.appsFile, nil
}

func (f *fakeCompose) AddCronDeploy(env string, services []string, cron string, namespace string) (string, error) {
	if f.addErr != nil {
		return "", f.addErr
	}
	f.cronEnv = env
	f.cronServices = services
	f.cronSchedule = cron
	return f.appsFile, nil
}

func (f *fakeCompose) AddRouteDeploy(env, service, port string, namespace string) (string, error) {
	if f.addErr != nil {
		return "", f.addErr
	}
	f.routeEnv = env
	f.routeService = service
	f.routePort = port
	return f.appsFile, nil
}

func (f *fakeCompose) GroupServicesByFile(files, services []string, namespace string) (map[string][]string, []string) {
	return f.grouped, f.missing
}

func (f *fakeCompose) AddConfig(file, configName, baseName string, services []string, namespace string) error {
	if f.addErr != nil {
		return f.addErr
	}
	f.addedConfigs = append(f.addedConfigs, file)
	return nil
}

func (f *fakeCompose) AddSecret(file, secretName, baseName string, services []string, namespace string) error {
	if f.addErr != nil {
		return f.addErr
	}
	f.addedSecrets = append(f.addedSecrets, file)
	return nil
}

type fakeShell struct {
	launched bool
	err      error
}

func (f *fakeShell) Launch() error {
	f.launched = true
	return f.err
}

type fakeDeploy struct {
	buildFiles       []string
	buildRegistry    string
	buildNamespace   string
	rolloutFiles     []string
	rolloutNamespace string
	canaryFiles      []string
	canaryServices   []string
	canaryReplicas   int
	canaryNamespace  string
	verbose          bool
	err              error
}

func (f *fakeDeploy) Build(files []string, registry string, namespace string, verbose bool) error {
	f.buildFiles = files
	f.buildRegistry = registry
	f.buildNamespace = namespace
	f.verbose = verbose
	return f.err
}

func (f *fakeDeploy) Rollout(files []string, namespace string, verbose bool) error {
	f.rolloutFiles = files
	f.rolloutNamespace = namespace
	f.verbose = verbose
	return f.err
}

func (f *fakeDeploy) Canary(files []string, services []string, replicas int, namespace string, verbose bool) error {
	f.canaryFiles = files
	f.canaryServices = services
	f.canaryReplicas = replicas
	f.canaryNamespace = namespace
	f.verbose = verbose
	return f.err
}

type fakeRouting struct {
	env       string
	namespace string
	verbose   bool
	err       error
}

func (f *fakeRouting) Apply(env string, namespace string, verbose bool) error {
	f.env = env
	f.namespace = namespace
	f.verbose = verbose
	return f.err
}

type fakeRoute struct {
	env       string
	url       string
	target    string
	namespace string
	file      string
	called    bool
	err       error
}

func (f *fakeRoute) AddRoute(env, url, target string, namespace string) (string, error) {
	f.called = true
	f.env = env
	f.url = url
	f.target = target
	f.namespace = namespace
	return f.file, f.err
}
