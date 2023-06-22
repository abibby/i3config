package i3config

import (
	"fmt"
	"strings"
)

type Command struct {
	name   string
	prefix string
	value  string
}

func NewCommand(name, value string) *Command {
	return &Command{
		name:  name,
		value: value,
	}
}

var (
	FocusUp    = NewCommand("focus", "up")
	FocusDown  = NewCommand("focus", "down")
	FocusLeft  = NewCommand("focus", "left")
	FocusRight = NewCommand("focus", "right")

	MoveUp    = NewCommand("move", "up")
	MoveDown  = NewCommand("move", "down")
	MoveLeft  = NewCommand("move", "left")
	MoveRight = NewCommand("move", "right")

	SplitHorizontal = NewCommand("split", "horizontal")
	SplitVertical   = NewCommand("split", "vertical")
	SplitToggle     = NewCommand("split", "toggle")

	LayoutDefault         = NewCommand("layout", "default")
	LayoutTabbed          = NewCommand("layout", "tabbed")
	LayoutStacking        = NewCommand("layout", "stacking")
	LayoutSplitVertical   = NewCommand("layout", "splitv")
	LayoutSplitHorizontal = NewCommand("layout", "splith")

	FullscreenToggle = NewCommand("fullscreen", "toggle")

	FloatingEnabled  = NewCommand("floating", "enabled")
	FloatingDisabled = NewCommand("floating", "disabled")

	Restart = NewCommand("restart", "")
	Reload  = NewCommand("reload", "")

	Kill = NewCommand("kill", "")
)

var funcKey = 0

func Exec(cmd string) *Command {
	return NewCommand("exec", escapeString(cmd))
}

func Mode(name string) *Command {
	return NewCommand("mode", escapeString(name))
}

func Workspace(name string) *Command {
	return NewCommand("workspace", escapeString(name))
}

func (c *Config) WorkspaceOutput(name string, outputs ...string) {
	cmd := NewCommand("workspace", escapeString(name)+" output "+strings.Join(outputs, " "))
	c.AddLine(cmd)
}

func MoveContainer(name string) *Command {
	return NewCommand("move", "container to workspace "+escapeString(name))
}

type Border int

func (b Border) Generate() string {
	return fmt.Sprintf("border pixel %d", b)
}

type Direction string

const (
	Up     Direction = "up"
	Down   Direction = "down"
	Left   Direction = "left"
	Right  Direction = "right"
	Width  Direction = "width"
	Height Direction = "height"
)

func ResizeGrow(direction Direction, amount int) *Command {
	return NewCommand("resize", fmt.Sprintf("grow %s %d px or %d ppt", direction, amount, amount))
}

func ResizeShrink(direction Direction, amount int) *Command {
	return NewCommand("resize", fmt.Sprintf("shrink %s %d px or %d ppt", direction, amount, amount))
}

type Size struct {
	Width  int
	Height int
}

func ResizeSet(size Size) *Command {
	cmd := "set"
	if size.Width != 0 {
		cmd += fmt.Sprintf(" width %d px", size.Width)
	}
	if size.Height != 0 {
		cmd += fmt.Sprintf(" height %d px", size.Height)
	}
	return &Command{
		name:  "resize",
		value: cmd,
	}
}

func (c *Config) Path() string {
	return c.path
}

func (c *Config) Recompile(configPath string) *Command {
	return c.ExecFunc(func() error {
		return c.RecompileFunc(configPath)
	})
}

func (c *Command) Generate() string {
	src := c.name
	if c.prefix != "" {
		src += " " + c.prefix
	}
	if c.value != "" {
		src += " " + c.value
	}
	return src
}

func (c *Command) NoStartupID() *Command {
	c.prefix = "--no-startup-id"
	return c
}

func (c *Config) OnStartup(cmd *Command) {
	c.AddLine(cmd)
}

func (c *Config) AlwaysOnStartup(cmd *Command) {
	cmd.name = "exec_always"
	c.OnStartup(cmd)
}
