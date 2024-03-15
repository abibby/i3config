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
	str = strings.ReplaceAll(str, `\`, `\\`)
	str = strings.ReplaceAll(str, `"`, `\\"`)
	return `"` + str + `"`
}

// example of escaped escaped quotes
// bindsym Mod4+g exec "emacsclient -c -e \\"(find-file \\\\"/tmp\\\\")\\""
