package stringify

import (
	"reflect"

	"github.com/rafaelmartins/synth-datagen/internal/ctypes"
	"github.com/rafaelmartins/synth-datagen/internal/utils"
)

type structSpec struct {
	field reflect.StructField
	ctype string
}

func stringifyStructType(ts *typeSpec) string {
	rv := "struct {"
	if len(ts.stype) > 0 {
		rv += "\n"
	}
	for _, field := range ts.stype {
		rv += lpadding(1) + field.ctype + " " + utils.FieldNameToSnake(field.field.Name) + ";\n"
	}
	return rv + "}"
}

func stringifyStructData(val reflect.Value, level uint8, ts *typeSpec) string {
	if len(ts.stype) == 0 {
		for _, field := range reflect.VisibleFields(val.Type()) {
			if ctype, err := ctypes.FromKind(field.Type.Kind()); err == nil && field.IsExported() {
				ts.stype = append(ts.stype, &structSpec{
					field: field,
					ctype: ctype,
				})
			}
		}
	}

	values := []string{}
	for _, field := range ts.stype {
		if s, err := ctypes.ToString(field.ctype, val.FieldByName(field.field.Name).Interface(), true); err == nil {
			values = append(values, s)
		}
	}
	return dumpValues(values, level)
}
