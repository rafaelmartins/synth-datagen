package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type ConfigModule struct {
	Identifier string                 `yaml:"-"`
	Name       string                 `yaml:"name"`
	Parameters map[string]interface{} `yaml:"parameters"`
}

type ConfigModules []*ConfigModule

func (c *ConfigModules) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.AliasNode {
		value = value.Alias
	}

	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("config: modules: not a mapping (line %d, column %d)", value.Line, value.Column)
	}

	identifier := ""
	for i, cnt := range value.Content {
		if i%2 == 0 {
			if err := cnt.Decode(&identifier); err != nil {
				return err
			}
		} else {
			m := &ConfigModule{
				Identifier: identifier,
			}
			if err := cnt.Decode(m); err != nil {
				return err
			}
			*c = append(*c, m)
		}
	}

	return nil
}
