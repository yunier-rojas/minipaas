package compose

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/compose-spec/compose-go/v2/cli"
	"github.com/compose-spec/compose-go/v2/types"

	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/entities"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/usecases"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/logger"
)

type ManagerService struct{}

func NewManagerService() *ManagerService {
	return &ManagerService{}
}

func (m *ManagerService) FilesForEnv(env string, cfg entities.MinipaasConfig) []string {
	files := make([]string, 0, len(cfg.Compose.Files))

	for _, f := range cfg.Compose.Files {
		if filepath.IsAbs(f) {
			files = append(files, f)
			continue
		}
		files = append(files, filepath.Join(env, f))
	}

	return files
}

func (m *ManagerService) DiscoverProjectFiles(env string) []string {
	root := filepath.Dir(filepath.Clean(env))

	files := make([]string, 0, 2)
	for _, names := range [][]string{
		{"compose.yaml", "compose.yml"},
		{"compose.build.yaml", "compose.build.yml"},
	} {
		for _, name := range names {
			path := filepath.Join(root, name)
			if info, err := os.Stat(path); err == nil && !info.IsDir() {
				files = append(files, path)
				break
			}
		}
	}

	return files
}

func (m *ManagerService) GroupServicesByFile(files, services []string, namespace string) (map[string][]string, []string) {
	svcPerFile := map[string][]string{}
	var missing []string

	for _, s := range services {
		target := ""
		for _, file := range files {
			project, _, err := loadComposeFile(file, namespace)
			if err != nil {
				continue
			}
			if _, ok := project.Services[s]; ok {
				target = file
			}
		}
		if target == "" {
			missing = append(missing, s)
			continue
		}
		svcPerFile[target] = append(svcPerFile[target], s)
	}

	return svcPerFile, missing
}

func (m *ManagerService) AddConfig(file, configName, baseName string, services []string, namespace string) error {
	project, _, err := loadComposeFile(file, namespace)
	if err != nil {
		return err
	}

	if err = addComposeConfig(project, configName, baseName, services); err != nil {
		return err
	}

	if _, err = saveComposeFile(file, project); err != nil {
		return err
	}

	logger.Info("Updated compose file with config: " + file)
	return nil
}

func (m *ManagerService) AddSecret(file, secretName, baseName string, services []string, namespace string) error {
	project, _, err := loadComposeFile(file, namespace)
	if err != nil {
		return err
	}

	if err = addComposeSecret(project, secretName, baseName, services); err != nil {
		return err
	}

	if _, err = saveComposeFile(file, project); err != nil {
		return err
	}

	logger.Info("Updated compose file with secret: " + file)
	return nil
}

func (m *ManagerService) AddWorkerDeploy(env string, services []string, namespace string) (string, error) {
	return m.updateApps(env, services, namespace, addComposeWorkerDeploy)
}

func (m *ManagerService) AddJobDeploy(env string, services []string, namespace string) (string, error) {
	return m.updateApps(env, services, namespace, addComposeJobDeploy)
}

func (m *ManagerService) AddCronDeploy(env string, services []string, cron string, namespace string) (string, error) {
	return m.updateApps(env, services, namespace, func(project *types.Project, name string) error {
		return addComposeCronDeploy(project, name, cron)
	})
}

func (m *ManagerService) AddRouteDeploy(env, service, port string, namespace string) (string, error) {
	return m.updateApps(env, []string{service}, namespace, func(project *types.Project, name string) error {
		return addComposeRouteDeploy(project, name, port)
	})
}

func (m *ManagerService) updateApps(env string, services []string, namespace string, apply func(*types.Project, string) error) (string, error) {
	fn := filepath.Join(env, entities.AppsFile)

	project, _, err := loadComposeFile(fn, namespace)
	if err != nil {
		return fn, err
	}

	for _, name := range services {
		if err = apply(project, name); err != nil {
			return fn, err
		}
	}

	if _, err = saveComposeFile(fn, project); err != nil {
		return fn, err
	}

	logger.Info("Updated deploy file: " + fn)
	return fn, nil
}

func (m *ManagerService) BuildApps(env string, files []string, registry string, namespace string) (string, error) {
	project, err := loadComposeProject(files, namespace)
	if err != nil {
		return "", err
	}

	stripDefaultNetworks(project)

	buildImages := map[string]struct{}{}
	for _, svc := range project.Services {
		if svc.Build != nil && svc.Image != "" {
			buildImages[svc.Image] = struct{}{}
		}
	}

	apps := &types.Project{Services: make(types.Services)}
	for name, svc := range project.Services {
		_, requiresBuild := buildImages[svc.Image]
		apps.Services[name] = createDeployService(svc, requiresBuild, registry, namespace)
	}

	apps.Name = ""
	data, err := apps.MarshalYAML()
	if err != nil {
		return "", err
	}

	if err = os.MkdirAll(env, 0755); err != nil {
		return "", err
	}

	target := filepath.Join(env, entities.AppsFile)
	if err = os.WriteFile(target, data, 0644); err != nil {
		return "", err
	}

	logger.Info("Wrote compose file: " + target)
	return entities.AppsFile, nil
}

func loadComposeProject(files []string, namespace string) (*types.Project, error) {
	options, err := cli.NewProjectOptions(
		files,
		cli.WithName(namespace),
		cli.WithResolvedPaths(false),
		cli.WithConsistency(false),
		cli.WithoutEnvironmentResolution,
	)
	if err != nil {
		return nil, err
	}

	return options.LoadProject(context.Background())
}

func createDeployService(svc types.ServiceConfig, requiresBuild bool, registry string, namespace string) types.ServiceConfig {
	srv := types.ServiceConfig{
		Networks: map[string]*types.ServiceNetworkConfig{
			namespace + "_network": {},
		},
		Labels: monitoringLabels(svc.Labels),
	}
	if requiresBuild {
		srv.Image = buildCommonImage(svc.Image, registry)
	}
	return srv
}

func monitoringLabels(existing map[string]string) map[string]string {
	labels := map[string]string{}
	if _, ok := existing[entities.AlertCPUThresholdLabel]; !ok {
		labels[entities.AlertCPUThresholdLabel] = entities.DefaultAlertCPUThreshold
	}
	if _, ok := existing[entities.AlertMemoryThresholdLabel]; !ok {
		labels[entities.AlertMemoryThresholdLabel] = entities.DefaultAlertMemoryThreshold
	}
	if len(labels) == 0 {
		return nil
	}
	return labels
}

func buildCommonImage(image, registry string) string {
	version := "${MINIPAAS_DEPLOY_VERSION}"

	if idx := strings.LastIndex(image, ":"); idx != -1 {
		image = image[:idx]
	}

	if registry != "" && !strings.Contains(image, "/") {
		image = "${MINIPAAS_IMAGE_PREFIX}/" + image
	}

	return image + ":" + version
}

func loadComposeFile(file string, namespace string) (*types.Project, string, error) {
	options, err := cli.NewProjectOptions(
		[]string{file},
		cli.WithEnv([]string{"MINIPAAS_DEPLOY_VERSION=${MINIPAAS_DEPLOY_VERSION}"}),
		cli.WithName(namespace),
		cli.WithResolvedPaths(false),
		cli.WithConsistency(false),
		cli.WithoutEnvironmentResolution,
	)
	if err != nil {
		return nil, file, err
	}

	project, err := options.LoadProject(context.Background())
	if err != nil {
		return nil, file, err
	}

	stripDefaultNetworks(project)

	return project, file, nil
}

func saveComposeFile(file string, project *types.Project) (string, error) {
	project.Name = ""
	data, err := project.MarshalYAML()
	if err != nil {
		return file, err
	}
	return file, os.WriteFile(file, data, 0644)
}

func stripDefaultNetworks(project *types.Project) {
	if project.Networks != nil {
		delete(project.Networks, "default")
	}

	for name, svc := range project.Services {
		if svc.Networks != nil {
			if _, exists := svc.Networks["default"]; exists {
				delete(svc.Networks, "default")
				project.Services[name] = svc
			}
		}
	}
}

func addComposeWorkerDeploy(project *types.Project, service string) error {
	svc, ok := project.Services[service]
	if !ok {
		return fmt.Errorf("service %s not found in project", service)
	}

	replicas := 1
	updateParallelism := uint64(1)
	updateDelay := types.Duration(10 * time.Second)
	rollbackParallelism := uint64(0)
	restartDelay := types.Duration(10 * time.Second)
	rollbackDelay := types.Duration(10 * time.Second)

	svc.Deploy = &types.DeployConfig{
		Mode:     "replicated",
		Replicas: &replicas,
		UpdateConfig: &types.UpdateConfig{
			Parallelism:   &updateParallelism,
			Order:         "start-first",
			FailureAction: "rollback",
			Delay:         updateDelay,
		},
		RestartPolicy: &types.RestartPolicy{
			Condition: "any",
			Delay:     &restartDelay,
		},
		RollbackConfig: &types.UpdateConfig{
			Parallelism: &rollbackParallelism,
			Order:       "stop-first",
			Delay:       rollbackDelay,
		},
	}

	project.Services[service] = svc
	return nil
}

func addComposeJobDeploy(project *types.Project, serviceName string) error {
	svc, ok := project.Services[serviceName]
	if !ok {
		return fmt.Errorf("service %s not found in project", serviceName)
	}

	replicas := 1
	updateParallelism := uint64(0)
	updateDelay := types.Duration(10 * time.Second)
	restartDelay := types.Duration(10 * time.Second)
	restartAttempts := uint64(10)

	svc.Deploy = &types.DeployConfig{
		Mode:     "replicated",
		Replicas: &replicas,
		UpdateConfig: &types.UpdateConfig{
			Parallelism:   &updateParallelism,
			Order:         "stop-first",
			FailureAction: "pause",
			Delay:         updateDelay,
		},
		RestartPolicy: &types.RestartPolicy{
			Condition:   "on-failure",
			Delay:       &restartDelay,
			MaxAttempts: &restartAttempts,
		},
	}

	project.Services[serviceName] = svc
	return nil
}

func addComposeCronDeploy(project *types.Project, serviceName string, cron string) error {
	svc, ok := project.Services[serviceName]
	if !ok {
		return fmt.Errorf("service %s not found in project", serviceName)
	}

	replicas := 0

	svc.Deploy = &types.DeployConfig{
		Mode:     "replicated",
		Replicas: &replicas,
		RestartPolicy: &types.RestartPolicy{
			Condition: "none",
		},
		Labels: map[string]string{
			"swarm.cronjob.enable":       "true",
			"swarm.cronjob.schedule":     cron,
			"swarm.cronjob.skip-running": "true",
		},
	}

	project.Services[serviceName] = svc
	return nil
}

func addComposeRouteDeploy(project *types.Project, service, port string) error {
	svc, ok := project.Services[service]
	if !ok {
		return fmt.Errorf("service %s not found in project", service)
	}

	replicas := 2

	healthInterval := types.Duration(10 * time.Second)
	healthTimeout := types.Duration(10 * time.Second)
	healthRetries := uint64(5)
	healthStart := types.Duration(10 * time.Second)

	updateParallelism := uint64(1)
	updateDelay := types.Duration(10 * time.Second)
	rollbackParallelism := uint64(0)
	restartDelay := types.Duration(10 * time.Second)
	rollbackDelay := types.Duration(10 * time.Second)

	svc.HealthCheck = &types.HealthCheckConfig{
		Test:        []string{"CMD-SHELL", fmt.Sprintf("wget -qO- --spider http://127.0.0.1:%s", port)},
		Interval:    &healthInterval,
		Timeout:     &healthTimeout,
		Retries:     &healthRetries,
		StartPeriod: &healthStart,
	}

	svc.Deploy = &types.DeployConfig{
		Mode:     "replicated",
		Replicas: &replicas,
		UpdateConfig: &types.UpdateConfig{
			Parallelism:   &updateParallelism,
			Order:         "start-first",
			FailureAction: "rollback",
			Delay:         updateDelay,
		},
		RestartPolicy: &types.RestartPolicy{
			Condition: "any",
			Delay:     &restartDelay,
		},
		RollbackConfig: &types.UpdateConfig{
			Parallelism: &rollbackParallelism,
			Order:       "stop-first",
			Delay:       rollbackDelay,
		},
	}

	project.Services[service] = svc
	return nil
}

func addComposeConfig(project *types.Project, config, name string, services []string) error {
	if project.Configs == nil {
		project.Configs = make(map[string]types.ConfigObjConfig)
	}

	project.Configs[config] = types.ConfigObjConfig{
		External: true,
	}

	replaced := map[string]bool{}
	for _, svcName := range services {
		svc, ok := project.Services[svcName]
		if !ok {
			return fmt.Errorf("service %s not found in project", svcName)
		}

		kept := make([]types.ServiceConfigObjConfig, 0, len(svc.Configs))
		for _, sc := range svc.Configs {
			if sc.Target == name {
				if sc.Source != config {
					replaced[sc.Source] = true
				}
				continue
			}
			kept = append(kept, sc)
		}
		kept = append(kept, types.ServiceConfigObjConfig{
			Source: config,
			Target: name,
		})
		svc.Configs = kept
		project.Services[svcName] = svc
	}

	for old := range replaced {
		if !configReferenced(project, old) {
			delete(project.Configs, old)
		}
	}

	return nil
}

func configReferenced(project *types.Project, source string) bool {
	for _, svc := range project.Services {
		for _, sc := range svc.Configs {
			if sc.Source == source {
				return true
			}
		}
	}
	return false
}

func addComposeSecret(project *types.Project, secret, name string, services []string) error {
	if project.Secrets == nil {
		project.Secrets = make(map[string]types.SecretConfig)
	}

	project.Secrets[secret] = types.SecretConfig{
		External: true,
	}

	replaced := map[string]bool{}
	for _, svcName := range services {
		svc, ok := project.Services[svcName]
		if !ok {
			return fmt.Errorf("service %s not found in project", svcName)
		}

		kept := make([]types.ServiceSecretConfig, 0, len(svc.Secrets))
		for _, sc := range svc.Secrets {
			if sc.Target == name {
				if sc.Source != secret {
					replaced[sc.Source] = true
				}
				continue
			}
			kept = append(kept, sc)
		}
		kept = append(kept, types.ServiceSecretConfig{
			Source: secret,
			Target: name,
		})
		svc.Secrets = kept
		project.Services[svcName] = svc
	}

	for old := range replaced {
		if !secretReferenced(project, old) {
			delete(project.Secrets, old)
		}
	}

	return nil
}

func secretReferenced(project *types.Project, source string) bool {
	for _, svc := range project.Services {
		for _, sc := range svc.Secrets {
			if sc.Source == source {
				return true
			}
		}
	}
	return false
}

var _ usecases.ComposePort = (*ManagerService)(nil)
