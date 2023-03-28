package stringify

import (
	"fmt"
	"strings"
)

const (
	lineWidth = 100
)

func lpadding(level uint8) string {
	return fmt.Sprintf("%*s", level*4, "")
}

func dumpValues(values []string, level uint8) string {
	rv := lpadding(level) + "{"
	if len(values) > 0 {
		rv += "\n"
	}

	line := lpadding(level + 1)
	for _, value := range values {
		if len(line)+len(value)+2 < lineWidth+1 { // we strip the trailing space
			line += value + ", "
		} else if strings.TrimSpace(line) == "" {
			rv += line + value + ",\n"
			line = lpadding(level + 1)
		} else {
			rv += strings.TrimRight(line, " ") + "\n"
			line = lpadding(level+1) + value + ", "
		}
	}
	if strings.TrimSpace(line) != "" {
		rv += strings.TrimRight(line, " ") + "\n"
	}

	if len(values) > 0 {
		rv += lpadding(level)
	}
	return rv + "}"
}
