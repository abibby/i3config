package i3config

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"strings"
)

type Command string

const (
	FocusUp    = Command("focus up")
	FocusDown  = Command("focus down")
	FocusLeft  = Command("focus left")
	FocusRight = Command("focus right")

	MoveUp    = Command("move up")
	MoveDown  = Command("move down")
	MoveLeft  = Command("move left")
	MoveRight = Command("move right")

	SplitHorizontal = Command("split horizontal")
	SplitVertical   = Command("split vertical")
	SplitToggle     = Command("split toggle")

	LayoutDefault         = Command("layout default")
	LayoutTabbed          = Command("layout tabbed")
	LayoutStacking        = Command("layout stacking")
	LayoutSplitVertical   = Command("layout splitv")
	LayoutSplitHorizontal = Command("layout splith")

	FullscreenToggle = Command("fullscreen toggle")

	FloatingEnabled  = Command("floating enabled")
	FloatingDisabled = Command("floating disabled")

	Restart = Command("restart")
	Reload  = Command("reload")

	Kill = Command("kill")
)

var funcKey = 0

func Exec(cmd string) Command {
	return Command("exec " + escapeString(cmd))
}

func Mode(name string) Command {
	return Command("mode " + escapeString(name))
}

func Workspace(name string) Command {
	return Command("workspace " + escapeString(name))
}

func (c *Config) WorkspaceOutput(name string, outputs ...string) {
	cmd := Command("workspace " + escapeString(name) + " output " + strings.Join(outputs, " "))
	c.AddLine(cmd)
}

func MoveContainer(name string) Command {
	return Command("move container to workspace " + escapeString(name))
}

func Border(size int) Command {
	return Command(fmt.Sprintf("border pixel %d", size))
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

func ResizeGrow(direction Direction, amount int) Command {
	return Command(fmt.Sprintf("resize grow %s %d px or %d ppt", direction, amount, amount))
}

func ResizeShrink(direction Direction, amount int) Command {
	return Command(fmt.Sprintf("resize shrink %s %d px or %d ppt", direction, amount, amount))
}

type Size struct {
	Width  int
	Height int
}

func ResizeSet(size Size) Command {
	cmd := "resize set"
	if size.Width != 0 {
		cmd += fmt.Sprintf(" width %d px", size.Width)
	}
	if size.Height != 0 {
		cmd += fmt.Sprintf(" height %d px", size.Height)
	}
	return Command(cmd)
}

func (c *Config) ExecFunc(cb func() error) Command {
	if c.subConfig {
		panic("ExecFunc must be used from a root config")
	}
	key := fmt.Sprint(funcKey)
	funcKey++
	c.funcs[key] = cb
	dir, file := path.Split(c.path)
	return Exec(fmt.Sprintf(`bash -c "cd '%s' && go run '%s' func %s"`, dir, file, key))
}
func (c *Config) Path() string {
	return c.path
}

func (c *Config) Recompile(configPath string) Command {
	return c.ExecFunc(func() error {
		b, err := exec.Command("go", "run", c.path).Output()
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(configPath, b, 0644)
		if err != nil {
			return err
		}
		return nil
	})
}

func (c Command) Generate() string {
	return string(c)
}

func (c Command) replacePrefix(old, new string) Command {
	cmd := string(c)
	if strings.HasPrefix(cmd, old) && !strings.HasPrefix(cmd, new) {
		cmd = new + cmd[len(old):]
	}
	return Command(cmd)
}
func (c Command) NoStartupID() Command {
	return c.replacePrefix("exec", "exec --no-startup-id")
}

func (c *Config) OnStartup(cmd Command) {
	c.AddLine(cmd)
}

func (c *Config) AlwaysOnStartup(cmd Command) {
	c.OnStartup(cmd.replacePrefix("exec", "exec_always"))
}
