package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type ConfigInclude struct {
	Path   string
	System bool
}

type ConfigIncludes []*ConfigInclude

func (c *ConfigIncludes) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.AliasNode {
		value = value.Alias
	}

	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("config: includes: not a mapping (line %d, column %d)", value.Line, value.Column)
	}

	path := ""
	for i, cnt := range value.Content {
		if i%2 == 0 {
			if err := cnt.Decode(&path); err != nil {
				return err
			}
		} else {
			system := false
			if err := cnt.Decode(&system); err != nil {
				return err
			}
			*c = append(*c, &ConfigInclude{
				Path:   path,
				System: system,
			})
		}
	}

	return nil
}
