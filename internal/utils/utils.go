package utils

import (
	"fmt"
	"math"
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

func FormatStringSliceWidth(s []string, width int) ([]string, error) {
	rv := []string{}
	for _, v := range s {
		f := fmt.Sprintf("%*s", width, v)
		if m := int(math.Abs(float64(width))); len(f) != m {
			return nil, fmt.Errorf("width overflow: %q (%d > %d)", f, len(f), m)
		}
		rv = append(rv, f)
	}
	return rv, nil
}
