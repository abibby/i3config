package i3config

import (
	"fmt"
	"strings"
)

type BarMode string

const (
	BarDock      BarMode = "dock"
	BarHide      BarMode = "hide"
	BarInvisible BarMode = "invisible"
)

type BarHiddenState string

const (
	BarShown  BarHiddenState = "show"
	BarHidden BarHiddenState = "hide"
)

type Bar struct {
	StatusCommand string
	I3BarCommand  string
	Mode          BarMode
	HiddenState   BarHiddenState
	Modifier      string
}

type BarWorkspaceColor struct {
	Border     Color
	Background Color
	Text       Color
}

func (b *BarWorkspaceColor) Generate() string {
	return fmt.Sprintf("%s %s %s",
		b.Border.Generate(),
		b.Background.Generate(),
		b.Text.Generate(),
	)
}

type BarColorConfig struct {
	Background        Color              `i3:"background"`
	StatusLine        Color              `i3:"statusline"`
	Separator         Color              `i3:"separator"`
	FocusedWorkspace  *BarWorkspaceColor `i3:"focused_workspace"`
	ActiveWorkspace   *BarWorkspaceColor `i3:"active_workspace"`
	InactiveWorkspace *BarWorkspaceColor `i3:"inactive_workspace"`
	UrgentWorkspace   *BarWorkspaceColor `i3:"urgent_workspace"`
}

func (b *BarColorConfig) Generate() string {
	lines := []string{}
	EachKey(b, func(key, value string) {
		lines = append(lines, key+" "+value)
	})
	return "colors {\n" + indent(strings.Join(lines, "\n")) + "\n}"
}

type BarConfig struct {
	*Config
}

func (b *BarConfig) Colors(c *BarColorConfig) {
	b.AddLine(c)
}

type BarPosition string

const (
	Top    BarPosition = "top"
	Bottom BarPosition = "Bottom"
)

func (b *BarConfig) Position(p BarPosition) {
	b.raw("position " + string(p))
}

func (b *BarConfig) StatusCommand(command string) {
	b.raw("status_command " + command)
}
func (b *BarConfig) TrayOutput(display string) {
	b.raw("tray_output " + display)
}

func (b *BarConfig) Generate() string {
	return fmt.Sprintf("bar {\n%s\n}", indent(b.Config.Generate()))
}

func (c *Config) Bar(bar func(*BarConfig)) {
	b := &BarConfig{
		Config: c.newSubConfig(),
	}
	bar(b)
	c.AddLine(b)
}
