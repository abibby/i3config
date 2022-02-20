package i3config

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/abibby/nulls"
	"github.com/pkg/errors"
)

type I3msgError struct {
	ParseError    bool   `json:"parse_error"`
	ErrorMessage  string `json:"error"`
	Input         string `json:"input"`
	ErrorPosition string `json:"errorposition"`
}

func (e *I3msgError) Error() string {
	if e.ParseError {
		return fmt.Sprintf("%s\n%s\n%s", e.Input, e.ErrorPosition, e.ErrorMessage)
	}
	return ""
}

type CommandResult struct {
	Success bool `json:"success"`
	*I3msgError
}

func i3msg(v interface{}, arg ...string) error {
	b, err := exec.Command("i3-msg", arg...).Output()
	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			return err
		}
	}
	return errors.Wrap(json.Unmarshal(b, v), "failed to parse")
}

func I3msg(commands ...Command) error {
	r := []*CommandResult{}
	strCommands := []string{}

	for _, cmd := range commands {
		strCommands = append(strCommands, cmd.Generate())
	}
	err := i3msg(&r, strings.Join(strCommands, "; "))
	if err != nil {
		return err
	}
	if len(r) < 0 {
		return fmt.Errorf("no result")
	}
	if r[0].Success == false && r[0].I3msgError != nil {
		return r[0].I3msgError
	}
	return nil
}

type Rect struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type I3msgWorkspace struct {
	ID      int    `json:"id"`
	Num     int    `json:"num"`
	Name    string `json:"name"`
	Visible bool   `json:"visible"`
	Focused bool   `json:"focused"`
	Rect    Rect   `json:"rect"`
	Output  string `json:"output"`
	Urgent  bool   `json:"urgent"`
}

func GetWorkspaces() ([]*I3msgWorkspace, error) {
	w := []*I3msgWorkspace{}
	err := i3msg(&w, "-t", "get_workspaces")
	return w, err
}

type I3msgOutput struct {
	Name             string        `json:"name"`
	Active           bool          `json:"active"`
	Primary          bool          `json:"primary"`
	Rect             Rect          `json:"rect"`
	CurrentWorkspace *nulls.String `json:"current_workspace"`
}

func GetOutputs() ([]*I3msgOutput, error) {
	o := []*I3msgOutput{}
	err := i3msg(&o, "-t", "get_outputs")
	return o, err
}

type I3msgNode struct {
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
	Nodes              []*I3msgNode  `json:"nodes"`
	FloatingNodes      []*I3msgNode  `json:"floating_nodes"`
	Focus              []int64       `json:"focus"`
	FullscreenMode     int           `json:"fullscreen_mode"`
	Sticky             bool          `json:"sticky"`
	Floating           string        `json:"floating"`
	Swallows           []interface{} `json:"swallows"`
}

func GetTree() (*I3msgNode, error) {
	t := &I3msgNode{}
	err := i3msg(&t, "-t", "get_tree")
	return t, errors.Wrap(err, "failed to run get_tree")
}

func (n *I3msgNode) Walk(cb func(n *I3msgNode) bool) {
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
