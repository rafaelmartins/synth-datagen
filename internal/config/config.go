package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	GlobalParameters map[string]interface{} `yaml:"global_parameters"`
	Output           map[string]struct {
		Includes ConfigIncludes `yaml:"includes"`
		Macros   ConfigMacros   `yaml:"macros"`
		Modules  []struct {
			Name       string                 `yaml:"name"`
			Parameters map[string]interface{} `yaml:"parameters"`
		} `yaml:"modules"`
	} `yaml:"output"`
}

func New(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rv := &Config{}
	if err := yaml.NewDecoder(f).Decode(rv); err != nil {
		return nil, err
	}

	return rv, nil
}
