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

func unescapeString(str string) string {
	return strings.ReplaceAll(str[1:len(str)-1], `\\"`, `"`)
}

func parseArgs(str string) []string {
	argv := []string{""}
	i := 0
	var quote rune
	escape := false
	for _, v := range str {
		if escape {
			argv[i] += string(v)
			continue
		}
		if v == '\\' {
			escape = true
			continue
		}
		if quote != 0 {
			if v == quote {
				quote = 0
				continue
			}
			argv[i] += string(v)
			continue
		}

		if v == '"' || v == '\'' {
			quote = v
			continue
		}
		if v == ' ' || v == '\t' {
			argv = append(argv, "")
			i++
			continue
		}

		argv[i] += string(v)
	}
	return argv
}
