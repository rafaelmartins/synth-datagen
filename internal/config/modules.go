package config

import (
	"fmt"

	"go.yaml.in/yaml/v3"
)

type Module struct {
	Identifier string         `yaml:"-"`
	Name       string         `yaml:"name"`
	Parameters map[string]any `yaml:"parameters"`
	Selectors  []string       `yaml:"selectors"`
}

type Modules []*Module

func (c *Modules) UnmarshalYAML(value *yaml.Node) error {
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
			m := &Module{
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
