package config

import (
	"os"

	variable "github.com/w-h-a/caus/api/variable/v1alpha1"
	"gopkg.in/yaml.v3"
)

func LoadConfig(path string) (*variable.DiscoveryConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg variable.DiscoveryConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// TODO: add validation

	return &cfg, nil
}
