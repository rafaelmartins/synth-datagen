package datareg

import (
	"errors"
	"fmt"
	"reflect"

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

func (p *DataReg) Evaluate(obj interface{}, local map[string]interface{}) error {
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
			return fmt.Errorf("datareg: parameter not defined: %s", utils.FieldNameToSnake(field.Name))
		}

		v := reflect.ValueOf(itf)

		// yaml library returns a slice of interfaces instead of a slice of the underlying type
		if vl, ok := itf.([]interface{}); ok {
			s, err := utils.CastInterfaceSlice(vl)
			if err != nil {
				return err
			}
			v = reflect.ValueOf(s)
		}

		if !v.CanConvert(field.Type) {
			return fmt.Errorf("datareg: invalid parameter value type: %s: parameter is %q, wants %q", utils.FieldNameToSnake(field.Name), v.Type(), field.Type)
		}

		fld := val.FieldByName(field.Name)
		if !fld.CanSet() {
			return fmt.Errorf("datareg: can't set value: %s", utils.FieldNameToSnake(field.Name))
		}
		fld.Set(v.Convert(field.Type))
	}
	return nil
}
