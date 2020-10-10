package i3config

import (
	"strings"
)

func indent(src string) string {
	lines := strings.Split(src, "\n")
	for i, line := range lines {
		lines[i] = "    " + line
	}
	return strings.Join(lines, "\n")
}

func escapeString(str string) string {
	return `"` + strings.ReplaceAll(str, `"`, `\\"`) + `"`
}
