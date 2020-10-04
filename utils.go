package i3config

import (
	"encoding/json"
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
	b, err := json.Marshal(str)
	if err != nil {
		panic(err)
	}
	return string(b)
}
