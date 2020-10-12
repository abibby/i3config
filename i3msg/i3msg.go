package i3msg

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/abibby/i3config"
	"github.com/abibby/nulls"
	"github.com/pkg/errors"
)

type Error struct {
	ParseError    bool   `json:"parse_error"`
	ErrorMessage  string `json:"error"`
	Input         string `json:"input"`
	ErrorPosition string `json:"errorposition"`
}

func (e *Error) Error() string {
	if e.ParseError {
		return fmt.Sprintf("%s\n%s\n%s", e.Input, e.ErrorPosition, e.ErrorMessage)
	}
	return ""
}

type CommandResult struct {
	Success bool `json:"success"`
	*Error
}

func run(v interface{}, arg ...string) error {
	b, err := exec.Command("i3-msg", arg...).Output()
	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			return err
		}
	}
	return errors.Wrap(json.Unmarshal(b, v), "failed to parse")
}

func Run(commands ...i3config.Command) error {
	r := []*CommandResult{}
	strCommands := []string{}

	for _, cmd := range commands {
		strCommands = append(strCommands, cmd.Generate())
	}
	err := run(&r, strings.Join(strCommands, "; "))
	if err != nil {
		return err
	}
	if len(r) < 0 {
		return fmt.Errorf("no result")
	}
	if r[0].Success == false && r[0].Error != nil {
		return r[0].Error
	}
	return nil
}

type Rect struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Workspace struct {
	ID      int    `json:"id"`
	Num     int    `json:"num"`
	Name    string `json:"name"`
	Visible bool   `json:"visible"`
	Focused bool   `json:"focused"`
	Rect    Rect   `json:"rect"`
	Output  string `json:"output"`
	Urgent  bool   `json:"urgent"`
}

func GetWorkspaces() ([]*Workspace, error) {
	w := []*Workspace{}
	err := run(&w, "-t", "get_workspaces")
	return w, err
}

type Output struct {
	Name             string        `json:"name"`
	Active           bool          `json:"active"`
	Primary          bool          `json:"primary"`
	Rect             Rect          `json:"rect"`
	CurrentWorkspace *nulls.String `json:"current_workspace"`
}

func GetOutputs() ([]*Output, error) {
	o := []*Output{}
	err := run(&o, "-t", "get_outputs")
	return o, err
}

type Node struct {
	ID                 int64         `json:"id"`
	Type               string        `json:"type"`
	Orientation        string        `json:"orientation"`
	ScratchpadState    string        `json:"scratchpad_state"`
	Percent            float64       `json:"percent"`
	Urgent             bool          `json:"urgent"`
	Focused            bool          `json:"focused"`
	Layout             string        `json:"layout"`
	WorkspaceLayout    string        `json:"workspace_layout"`
	LastSplitLayout    string        `json:"last_split_layout"`
	Border             string        `json:"border"`
	CurrentBorderWidth int           `json:"current_border_width"`
	Rect               Rect          `json:"rect"`
	DecoRect           Rect          `json:"deco_rect"`
	WindowRect         Rect          `json:"window_rect"`
	Geometry           Rect          `json:"geometry"`
	Name               string        `json:"name"`
	Window             *nulls.Int    `json:"window"`
	WindowType         *nulls.String `json:"window_type"`
	Nodes              []*Node       `json:"nodes"`
	FloatingNodes      []*Node       `json:"floating_nodes"`
	Focus              []int64       `json:"focus"`
	FullscreenMode     int           `json:"fullscreen_mode"`
	Sticky             bool          `json:"sticky"`
	Floating           string        `json:"floating"`
	Swallows           []interface{} `json:"swallows"`
}

func GetTree() (*Node, error) {
	t := &Node{}
	err := run(&t, "-t", "get_tree")
	return t, errors.Wrap(err, "failed to run get_tree")
}

func (n *Node) Walk(cb func(n *Node) bool) {
	if cb(n) {
		return
	}
	for _, subNode := range n.Nodes {
		subNode.Walk(cb)
	}
	for _, subNode := range n.FloatingNodes {
		subNode.Walk(cb)
	}
}
