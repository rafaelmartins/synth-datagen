package config

import (
	"fmt"
	"reflect"

	"github.com/rafaelmartins/synth-datagen/internal/convert"
	"github.com/rafaelmartins/synth-datagen/internal/ctypes"
	"gopkg.in/yaml.v3"
)

type Variable struct {
	Identifier string      `yaml:"-"`
	Type       string      `yaml:"type"`
	Value      interface{} `yaml:"value"`
	Attributes []string    `yaml:"attributes"`
}

type Variables []*Variable

func (c *Variables) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.AliasNode {
		value = value.Alias
	}

	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("config: variables: not a mapping (line %d, column %d)", value.Line, value.Column)
	}

	identifier := ""
	for i, cnt := range value.Content {
		if i%2 == 0 {
			if err := cnt.Decode(&identifier); err != nil {
				return err
			}
		} else {
			m := &Variable{
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
				if m.Type != "" {
					v := reflect.ValueOf(m.Value)
					if v.Kind() == reflect.Slice {
						value, err := convert.Slice(m.Value, m.Type)
						if err != nil {
							return err
						}
						m.Value = value
					} else if ctypes.TypeIsScalar(v.Type()) {
						value, err := convert.Scalar(m.Value, m.Type)
						if err != nil {
							return err
						}
						m.Value = value
					} else {
						return fmt.Errorf("config: variables: unsupported value (line %d, column %d)", cnt.Line, cnt.Column)
					}
				}
			}
			*c = append(*c, m)
		}
	}

	return nil
}
