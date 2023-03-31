package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	GlobalParameters map[string]interface{} `yaml:"global_parameters"`
	Output           map[string]struct {
		Includes  Includes  `yaml:"includes"`
		Macros    Macros    `yaml:"macros"`
		Variables Variables `yaml:"variables"`
		Modules   Modules   `yaml:"modules"`
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
