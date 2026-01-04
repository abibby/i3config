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

func unescapeString(str string) string {
	if len(str) < 2 || str[0] != '"' || str[len(str)-1] != '"' {
		return str
	}
	str = str[1 : len(str)-1]
	str = strings.ReplaceAll(str, `\\"`, `"`)
	str = strings.ReplaceAll(str, `\\`, `\`)
	return str
}

// example of escaped escaped quotes
// bindsym Mod4+g exec "emacsclient -c -e \\"(find-file \\\\"/tmp\\\\")\\""
