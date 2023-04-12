package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Output struct {
	HeaderOutput string    `yaml:"-"`
	ChartsOutput string    `yaml:"charts_output"`
	Includes     Includes  `yaml:"includes"`
	Macros       Macros    `yaml:"macros"`
	Variables    Variables `yaml:"variables"`
	Modules      Modules   `yaml:"modules"`
}

type Outputs []*Output

func (c *Outputs) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.AliasNode {
		value = value.Alias
	}

	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("config: outputs: not a mapping (line %d, column %d)", value.Line, value.Column)
	}

	header := ""
	for i, cnt := range value.Content {
		if i%2 == 0 {
			if err := cnt.Decode(&header); err != nil {
				return err
			}
		} else {
			m := &Output{
				HeaderOutput: header,
			}
			if err := cnt.Decode(m); err != nil {
				return err
			}
			*c = append(*c, m)
		}
	}

	return nil
}
