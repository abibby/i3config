package i3config

import (
	"fmt"
	"strings"
)

type Gaps struct {
	Inner int
	Outer int
	Smart bool
}

func (g Gaps) Generate() string {
	src := ""
	if g.Inner > 0 {
		src += fmt.Sprintf("gaps inner %d", g.Inner) + "\n"
	}
	if g.Outer > 0 {
		src += fmt.Sprintf("gaps outer %d", g.Outer) + "\n"
	}
	if g.Smart {
		src += "smart_gaps on" + "\n"
	}

	return strings.TrimRight(src, "\n")
}

func (g Gaps) GenerateYabai() string {
	src := ""

	// yabai -m config top_padding    20
	// yabai -m config bottom_padding 20
	// yabai -m config left_padding   20
	// yabai -m config right_padding  20
	// yabai -m config window_gap     20
	if g.Inner > 0 {
		src += fmt.Sprintf("yabai -m config window_gap %d", g.Inner) + "\n"
	}
	if g.Outer > 0 {
		src += fmt.Sprintf("yabai -m config top_padding %d", g.Outer) + "\n"
		src += fmt.Sprintf("yabai -m config bottom_padding %d", g.Outer) + "\n"
		src += fmt.Sprintf("yabai -m config left_padding %d", g.Outer) + "\n"
		src += fmt.Sprintf("yabai -m config right_padding %d", g.Outer) + "\n"
	}
	if g.Smart {
		// src += "smart_gaps on" + "\n"
	}

	return strings.TrimRight(src, "\n")
}

func (c *Config) Gaps(gaps Gaps) {
	c.AddLine(gaps)
}
