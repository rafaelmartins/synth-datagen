package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type ConfigMacro struct {
	Identifier string      `yaml:"identifier"`
	Value      interface{} `yaml:"value"`
	Hex        bool        `yaml:"hex"`
}

type ConfigMacros []*ConfigMacro

func (c *ConfigMacros) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.AliasNode {
		value = value.Alias
	}

	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("config: macros: not a mapping (line %d, column %d)", value.Line, value.Column)
	}

	identifier := ""
	for i, cnt := range value.Content {
		if i%2 == 0 {
			if err := cnt.Decode(&identifier); err != nil {
				return err
			}
		} else {
			m := &ConfigMacro{
				Identifier: identifier,
			}
			if cnt.Kind == yaml.ScalarNode {
				if err := cnt.Decode(&m.Value); err != nil {
					return err
				}
			} else {
				if err := cnt.Decode(m); err != nil {
					return err
				}
			}
			*c = append(*c, m)
		}
	}

	return nil
}