package config

import (
	"fmt"

	"go.yaml.in/yaml/v3"
)

type Include struct {
	Path   string
	System bool
}

type Includes []*Include

func (c *Includes) UnmarshalYAML(value *yaml.Node) error {
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
			*c = append(*c, &Include{
				Path:   path,
				System: system,
			})
		}
	}

	return nil
}
