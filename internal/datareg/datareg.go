package datareg

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/rafaelmartins/synth-datagen/internal/convert"
	"github.com/rafaelmartins/synth-datagen/internal/selector"
	"github.com/rafaelmartins/synth-datagen/internal/utils"
)

type DataReg struct {
	global map[string]interface{}
}

func New(global map[string]interface{}) *DataReg {
	return &DataReg{
		global: global,
	}
}

func lookup(m map[string]interface{}, key string) (interface{}, bool) {
	var (
		rv    interface{}
		found bool
	)
	key = utils.FieldNameToSnake(key)
	for k, v := range m {
		if utils.FieldNameToSnake(k) == key {
			rv = v
			found = true
		}
	}
	return rv, found
}

func (p *DataReg) Evaluate(obj interface{}, local map[string]interface{}, slt *selector.Selector) error {
	if obj == nil {
		return errors.New("datareg: got nil")
	}

	typ := reflect.TypeOf(obj)
	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return errors.New("datareg: not a struct pointer")
	}

	val := reflect.ValueOf(obj).Elem()
	for _, field := range reflect.VisibleFields(typ.Elem()) {
		if !field.IsExported() {
			continue
		}

		var itf interface{}
		if i, ok := lookup(local, field.Name); ok {
			itf = i
		} else if i, ok := lookup(p.global, field.Name); ok {
			itf = i
		}
		if itf == nil {
			found := ""
			if s, ok := field.Tag.Lookup("selectors"); slt != nil && ok {
				for _, sl := range strings.Split(s, ",") {
					if st := strings.TrimSpace(sl); st != "" && slt.IsSelected(st) {
						found = st
						break
					}
				}
			}
			if found != "" {
				return fmt.Errorf("datareg: parameter not defined: %s (required by selector %q)", utils.FieldNameToSnake(field.Name), found)
			}
			if field.Type.Kind() != reflect.Pointer {
				return fmt.Errorf("datareg: parameter not defined: %s (required, not a pointer)", utils.FieldNameToSnake(field.Name))
			}
			continue
		}

		v := reflect.ValueOf(itf)

		// yaml library returns a slice of interfaces instead of a slice of the underlying type
		if vl, ok := itf.([]interface{}); ok {
			s, err := convert.Slice(vl, "")
			if err != nil {
				return err
			}
			v = reflect.ValueOf(s)
		}

		t := field.Type
		if t.Kind() == reflect.Pointer {
			t = t.Elem()
		}

		if !v.CanConvert(t) {
			return fmt.Errorf("datareg: invalid parameter value type: %s: parameter is %q, wants %q", utils.FieldNameToSnake(field.Name), v.Type(), t)
		}

		fld := val.FieldByName(field.Name)
		if !fld.CanSet() {
			return fmt.Errorf("datareg: can't set value: %s", utils.FieldNameToSnake(field.Name))
		}

		if fld.Kind() == reflect.Pointer {
			fld.Set(reflect.New(t))
			fld.Elem().Set(v.Convert(t))
		} else {
			fld.Set(v.Convert(t))
		}
	}
	return nil
}
