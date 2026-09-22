package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/entities"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/usecases"
)

type StoreService struct{}

func NewStoreService() *StoreService {
	return &StoreService{}
}

func (s *StoreService) Load(env string) (entities.MinipaasConfig, error) {
	fn := filepath.Join(env, entities.ConfigFile)
	data, err := os.ReadFile(fn)
	if err != nil {
		return entities.MinipaasConfig{}, err
	}

	var cfg entities.MinipaasConfig
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return entities.MinipaasConfig{}, err
	}

	return cfg, nil
}

func (s *StoreService) Save(env string, cfg entities.MinipaasConfig) (string, error) {
	fn := filepath.Join(env, entities.ConfigFile)

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fn, err
	}

	if err = os.MkdirAll(env, 0755); err != nil {
		return fn, err
	}

	return fn, os.WriteFile(fn, data, 0644)
}

var _ usecases.ConfigStorePort = (*StoreService)(nil)
