package config

import (
	"fmt"
	"reflect"

	"github.com/expr-lang/expr"
	"gopkg.in/yaml.v3"
	"rafaelmartins.com/p/synth-datagen/internal/convert"
	"rafaelmartins.com/p/synth-datagen/internal/ctypes"
)

type Macro struct {
	Identifier string                 `yaml:"-"`
	Type       string                 `yaml:"type"`
	Value      interface{}            `yaml:"value"`
	Hex        bool                   `yaml:"hex"`
	Raw        bool                   `yaml:"raw"`
	Eval       bool                   `yaml:"eval"`
	EvalEnv    map[string]interface{} `yaml:"eval_env"`
}

type Macros []*Macro

func (c *Macros) UnmarshalYAML(value *yaml.Node) error {
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
			m := &Macro{
				Identifier: identifier,
			}

			if cnt.Kind == yaml.AliasNode {
				cnt = cnt.Alias
			}

			if cnt.Kind == yaml.ScalarNode {
				if err := cnt.Decode(&m.Value); err != nil {
					return err
				}
			} else {
				if err := cnt.Decode(m); err != nil {
					return err
				}

				if m.Eval || len(m.EvalEnv) > 0 {
					if input, ok := m.Value.(string); ok {
						out, err := expr.Eval(input, m.EvalEnv)
						if err != nil {
							return err
						}
						m.Value = out
					}
				}

				if !ctypes.TypeIsScalar(reflect.TypeOf(m.Value)) {
					return fmt.Errorf("config: macros: value is not a scalar (line %d, column %d)", cnt.Line, cnt.Column)
				}

				if m.Type != "" {
					value, err := convert.Scalar(m.Value, m.Type)
					if err != nil {
						return err
					}
					m.Value = value
				}
			}
			*c = append(*c, m)
		}
	}

	return nil
}
