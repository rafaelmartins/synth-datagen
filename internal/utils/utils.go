package utils

import (
	"unicode"
)

func FieldNameToSnake(name string) string {
	rv := ""
	for idx, c := range name {
		if idx == 0 {
			rv += string(unicode.ToLower(c))
			continue
		}

		if unicode.IsUpper(c) {
			rv += "_" + string(unicode.ToLower(c))
		} else if unicode.IsLower(c) {
			rv += string(c)
		} else {
			rv += "_"
		}
	}
	return rv
}
