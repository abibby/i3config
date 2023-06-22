package i3config

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var ErrUnknownCommand = errors.New("unknown command")

func (c *Config) ExecFunc(cb func() error) *Command {
	if c.subConfig {
		panic("ExecFunc must be used from a root config")
	}
	key := fmt.Sprint(funcKey)
	funcKey++
	c.funcs[key] = cb
	return NewCommand("func", key)
}

func (c *Command) RunYabai(config *Config) error {
	runners := map[string]func(c *Command) error{
		"exec":      runExec,
		"focus":     runFocus,
		"move":      runMove,
		"resize":    runResize,
		"workspace": runWorkspace,
		"func":      config.runFunc,
	}
	runner, ok := runners[c.name]
	if !ok {
		return nil
	}
	err := runner(c)
	if err != nil {
		return fmt.Errorf("%s: %w", c.Generate(), err)
	}
	return nil
}

func runExec(c *Command) error {
	argv := parseArgs(unescapeString(c.value))

	b, err := exec.Command("open", append([]string{"-a", argv[0], "-n", "--args"}, argv[1:]...)...).CombinedOutput()
	if string(b) == fmt.Sprintf("Unable to find application named '%s'\n", argv[0]) {
		err = exec.Command(argv[0], argv[1:]...).Run()
	}
	if err != nil {
		return err
	}

	return nil
}

func runResize(c *Command) error {
	parts := strings.Split(c.value, " ")
	horizontal := 0
	vertial := 0
	direction := "right"
	if parts[0] == "shrink" && parts[1] == "width" {
		direction = "right"
		horizontal = -1
	} else if parts[0] == "grow" && parts[1] == "width" {
		direction = "right"
		horizontal = 1
	} else if parts[0] == "shrink" && parts[1] == "height" {
		direction = "bottom"
		vertial = -1
	} else if parts[0] == "grow" && parts[1] == "height" {
		direction = "bottom"
		vertial = 1
	}

	amount, err := strconv.Atoi(parts[2])
	if err != nil {
		return err
	}

	err = yabai("window", "--resize", fmt.Sprintf("%s:%d:%d", direction, amount*horizontal, amount*vertial))
	if err == nil {
		return nil
	}

	amount = -amount
	if direction == "bottom" {
		direction = "top"
	} else if direction == "right" {
		direction = "left"
	}

	return yabai("window", "--resize", fmt.Sprintf("%s:%d:%d", direction, amount*horizontal, amount*vertial))

}
func runMove(c *Command) error {
	direction := ""
	follow := false
	if c == MoveUp {
		direction = "north"
	}
	if c == MoveDown {
		direction = "south"
	}
	if c == MoveLeft {
		direction = "west"
	}
	if c == MoveRight {
		direction = "east"
	}
	if direction != "" {
		err := yabai("window", "--swap", direction)
		if err == nil {
			return nil
		}
		nextSpace, err := getSpace(direction)
		if err != nil {
			return err
		}

		c.value = fmt.Sprintf("container to workspace \"%d\"", nextSpace.Index)
		follow = true
	}

	workspacePrefix := "container to workspace "
	if strings.HasPrefix(c.value, workspacePrefix) {
		windowID := 0
		if follow {
			w, err := yabaiQueryActiveWindow()
			if err != nil {
				log.Print(err)
			} else {
				windowID = w.ID
			}
		}
		workspace := unescapeString(c.value[len(workspacePrefix):])

		err := yabai("window", "--space", workspace)
		if err != nil {
			return err
		}
		if windowID != 0 {
			return yabai("window", "--focus", fmt.Sprint(windowID))
		}
		return nil
	}

	return ErrUnknownCommand
}

func runFocus(c *Command) error {
	direction := ""
	if c == FocusUp {
		direction = "north"
	}
	if c == FocusDown {
		direction = "south"
	}
	if c == FocusLeft {
		direction = "west"
	}
	if c == FocusRight {
		direction = "east"
	}

	if direction == "" {
		return ErrUnknownCommand
	}
	err := yabai("window", "--focus", direction)
	if err == nil {
		return nil
	}

	space, err := getSpace(direction)
	if err != nil {
		return err
	}

	return yabai("space", "--focus", fmt.Sprint(space.Index))
}

func runWorkspace(c *Command) error {
	return yabai("space", "--focus", unescapeString(c.value))
}

func (c *Config) runFunc(cmd *Command) error {
	fn, ok := c.funcs[cmd.value]
	if !ok {
		return fmt.Errorf("no func %s", cmd.value)
	}
	return fn()
}
func getSpace(direction string) (*yabaiSpace, error) {
	spaces, err := yabaiQuerySpaces()
	if err != nil {
		return nil, err
	}
	displays, err := yabaiQueryDisplays()
	if err != nil {
		return nil, err
	}
	var activeSpace *yabaiSpace
	var nextSpace *yabaiSpace
	var activeDisplay *yabaiDisplay
	var nextDisplay *yabaiDisplay

	for _, s := range spaces {
		if s.HasFocus {
			activeSpace = s
			break
		}
	}

	for _, d := range displays {
		if d.Index == activeSpace.DisplayIndex {
			activeDisplay = d
			break
		}
	}
	for _, d := range displays {
		if d.Frame.Y == activeDisplay.Frame.Y {
			if direction == "east" && d.Frame.X == (activeDisplay.Frame.X+activeDisplay.Frame.Width) {
				nextDisplay = d
				break
			}
			if direction == "west" && d.Frame.X == (activeDisplay.Frame.X-activeDisplay.Frame.Width) {
				nextDisplay = d
				break
			}
		}
		if d.Frame.X == activeDisplay.Frame.X {
			if direction == "north" && d.Frame.Y == (activeDisplay.Frame.Y+activeDisplay.Frame.Height) {
				nextDisplay = d
				break
			}
			if direction == "south" && d.Frame.Y == (activeDisplay.Frame.Y-activeDisplay.Frame.Height) {
				nextDisplay = d
				break
			}
		}
	}

	if nextDisplay == nil {
		return nil, fmt.Errorf("no display %s of current display", direction)
	}

	for _, s := range spaces {
		if s.IsVisible && s.DisplayIndex == nextDisplay.Index {
			nextSpace = s
			break
		}
	}
	return nextSpace, nil
}
