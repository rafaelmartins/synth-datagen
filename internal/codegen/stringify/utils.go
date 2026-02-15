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
	rv := strings.Builder{}
	rv.WriteString(lpadding(level) + "{")
	if len(values) > 0 {
		rv.WriteString("\n")
	}

	line := lpadding(level + 1)
	for _, value := range values {
		if len(line)+len(value)+2 < lineWidth+1 { // we strip the trailing space
			line += value + ", "
		} else if strings.TrimSpace(line) == "" {
			rv.WriteString(line + value + ",\n")
			line = lpadding(level + 1)
		} else {
			rv.WriteString(strings.TrimRight(line, " ") + "\n")
			line = lpadding(level+1) + value + ", "
		}
	}
	if strings.TrimSpace(line) != "" {
		rv.WriteString(strings.TrimRight(line, " ") + "\n")
	}

	if len(values) > 0 {
		rv.WriteString(lpadding(level))
	}
	return rv.String() + "}"
}
