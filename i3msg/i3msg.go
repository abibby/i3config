package i3msg

import (
	"encoding/json"
	"fmt"
	"os/exec"

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

func Run(cmd i3config.Command) error {
	r := []*CommandResult{}
	err := run(&r, cmd.Generate())
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

/*{
  "id": 94681385583136,
  "type": "root",
  "orientation": "horizontal",
  "scratchpad_state": "none",
  "percent": 1,
  "urgent": false,
  "focused": false,
  "layout": "splith",
  "workspace_layout": "default",
  "last_split_layout": "splith",
  "border": "normal",
  "current_border_width": -1,
  "rect": {
    "x": 0,
    "y": 0,
    "width": 6400,
    "height": 1440
  },
  "deco_rect": {
    "x": 0,
    "y": 0,
    "width": 0,
    "height": 0
  },
  "window_rect": {
    "x": 0,
    "y": 0,
    "width": 0,
    "height": 0
  },
  "geometry": {
    "x": 0,
    "y": 0,
    "width": 0,
    "height": 0
  },
  "name": "root",
  "window": null,
  "window_type": null,
  "nodes": [
    {
      "id": 94681385585680,
      "type": "output",
      "orientation": "none",
      "scratchpad_state": "none",
      "percent": 0.25,
      "urgent": false,
      "focused": false,
      "layout": "output",
      "workspace_layout": "default",
      "last_split_layout": "splith",
      "border": "normal",
      "current_border_width": -1,
      "rect": {
        "x": 0,
        "y": 0,
        "width": 7680,
        "height": 4320
      },
      "deco_rect": {
        "x": 0,
        "y": 0,
        "width": 0,
        "height": 0
      },
      "window_rect": {
        "x": 0,
        "y": 0,
        "width": 0,
        "height": 0
      },
      "geometry": {
        "x": 0,
        "y": 0,
        "width": 0,
        "height": 0
      },
      "name": "__i3",
      "window": null,
      "window_type": null,
      "nodes": [
        {
          "id": 94681385586160,
          "type": "con",
          "orientation": "horizontal",
          "scratchpad_state": "none",
          "percent": 1,
          "urgent": false,
          "focused": false,
          "output": "__i3",
          "layout": "splith",
          "workspace_layout": "default",
          "last_split_layout": "splith",
          "border": "normal",
          "current_border_width": -1,
          "rect": {
            "x": 0,
            "y": 0,
            "width": 0,
            "height": 0
          },
          "deco_rect": {
            "x": 0,
            "y": 0,
            "width": 0,
            "height": 0
          },
          "window_rect": {
            "x": 0,
            "y": 0,
            "width": 0,
            "height": 0
          },
          "geometry": {
            "x": 0,
            "y": 0,
            "width": 0,
            "height": 0
          },
          "name": "content",
          "window": null,
          "window_type": null,
          "nodes": [
            {
              "id": 94681385586640,
              "type": "workspace",
              "orientation": "none",
              "scratchpad_state": "none",
              "percent": 1,
              "urgent": false,
              "focused": false,
              "output": "__i3",
              "layout": "splith",
              "workspace_layout": "default",
              "last_split_layout": "splith",
              "border": "normal",
              "current_border_width": -1,
              "rect": {
                "x": 0,
                "y": 0,
                "width": 0,
                "height": 0
              },
              "deco_rect": {
                "x": 0,
                "y": 0,
                "width": 0,
                "height": 0
              },
              "window_rect": {
                "x": 0,
                "y": 0,
                "width": 0,
                "height": 0
              },
              "geometry": {
                "x": 0,
                "y": 0,
                "width": 0,
                "height": 0
              },
              "name": "__i3_scratch",
              "num": -1,
              "gaps": {
                "inner": 0,
                "outer": 0,
                "top": 0,
                "right": 0,
                "bottom": 0,
                "left": 0
              },
              "window": null,
              "window_type": null,
              "nodes": [],
              "floating_nodes": [],
              "focus": [],
              "fullscreen_mode": 1,
              "sticky": false,
              "floating": "auto_off",
              "swallows": []
            }
          ],
          "floating_nodes": [],
          "focus": [
            94681385586640
          ],
          "fullscreen_mode": 0,
          "sticky": false,
          "floating": "auto_off",
          "swallows": []
        }
      ],
      "floating_nodes": [],
      "focus": [
        94681385586160
      ],
      "fullscreen_mode": 0,
      "sticky": false,
      "floating": "auto_off",
      "swallows": []
    },
    {
      "id": 94681385595040,
      "type": "output",
      "orientation": "none",
      "scratchpad_state": "none",
      "percent": 0.25,
      "urgent": false,
      "focused": false,
      "layout": "output",
      "workspace_layout": "default",
      "last_split_layout": "splith",
      "border": "normal",
      "current_border_width": -1,
      "rect": {
        "x": 0,
        "y": 0,
        "width": 1920,
        "height": 1080
      },
      "deco_rect": {
        "x": 0,
        "y": 0,
        "width": 0,
        "height": 0
      },
      "window_rect": {
        "x": 0,
        "y": 0,
        "width": 0,
        "height": 0
      },
      "geometry": {
        "x": 0,
        "y": 0,
        "width": 0,
        "height": 0
      },
      "name": "HDMI1",
      "window": null,
      "window_type": null,
      "nodes": [
        {
          "id": 94681385595552,
          "type": "dockarea",
          "orientation": "none",
          "scratchpad_state": "none",
          "percent": 0.3333333333333333,
          "urgent": false,
          "focused": false,
          "output": "HDMI1",
          "layout": "dockarea",
          "workspace_layout": "default",
          "last_split_layout": "splith",
          "border": "normal",
          "current_border_width": -1,
          "rect": {
            "x": 0,
            "y": 0,
            "width": 1920,
            "height": 24
          },
          "deco_rect": {
            "x": 0,
            "y": 0,
            "width": 0,
            "height": 0
          },
          "window_rect": {
            "x": 0,
            "y": 0,
            "width": 0,
            "height": 0
          },
          "geometry": {
            "x": 0,
            "y": 0,
            "width": 0,
            "height": 0
          },
          "name": "topdock",
          "window": null,
          "window_type": null,
          "nodes": [
            {
              "id": 94681385791104,
              "type": "con",
              "orientation": "none",
              "scratchpad_state": "none",
              "percent": 1,
              "urgent": false,
              "focused": false,
              "output": "HDMI1",
              "layout": "splith",
              "workspace_layout": "default",
              "last_split_layout": "splith",
              "border": "normal",
              "current_border_width": 2,
              "rect": {
                "x": 0,
                "y": 0,
                "width": 1920,
                "height": 24
              },
              "deco_rect": {
                "x": 0,
                "y": 0,
                "width": 0,
                "height": 0
              },
              "window_rect": {
                "x": 0,
                "y": 0,
                "width": 1920,
                "height": 24
              },
              "geometry": {
                "x": 0,
                "y": 1056,
                "width": 1920,
                "height": 24
              },
              "name": "i3bar for output HDMI1",
              "window": 85983244,
              "window_type": "unknown",
              "window_properties": {
                "class": "i3bar",
                "instance": "i3bar",
                "title": "i3bar for output HDMI1",
                "transient_for": null
              },
              "nodes": [],
              "floating_nodes": [],
              "focus": [],
              "fullscreen_mode": 0,
              "sticky": false,
              "floating": "auto_off",
              "swallows": []
            }
          ],
          "floating_nodes": [],
          "focus": [
            94681385791104
          ],
          "fullscreen_mode": 0,
          "sticky": false,
          "floating": "auto_off",
          "swallows": [
            {
              "dock": 2,
              "insert_where": 2
            }
          ]
        },
        {
          "id": 94681385598704,
          "type": "con",
          "orientation": "horizontal",
          "scratchpad_state": "none",
          "percent": 0.3333333333333333,
          "urgent": false,
          "focused": false,
          "output": "HDMI1",
          "layout": "splith",
          "workspace_layout": "default",
          "last_split_layout": "splith",
          "border": "normal",
          "current_border_width": -1,
          "rect": {
            "x": 0,
            "y": 24,
            "width": 1920,
            "height": 1056
          },
          "deco_rect": {
            "x": 0,
            "y": 0,
            "width": 0,
            "height": 0
          },
          "window_rect": {
            "x": 0,
            "y": 0,
            "width": 0,
            "height": 0
          },
          "geometry": {
            "x": 0,
            "y": 0,
            "width": 0,
            "height": 0
          },
          "name": "content",
          "window": null,
          "window_type": null,
          "nodes": [
            {
              "id": 94681385599216,
              "type": "workspace",
              "orientation": "horizontal",
              "scratchpad_state": "none",
              "percent": 0.5,
              "urgent": false,
              "focused": false,
              "output": "HDMI1",
              "layout": "splith",
              "workspace_layout": "default",
              "last_split_layout": "splith",
              "border": "normal",
              "current_border_width": -1,
              "rect": {
                "x": 0,
                "y": 0,
                "width": 1920,
                "height": 1080
              },
              "deco_rect": {
                "x": 0,
                "y": 0,
                "width": 0,
                "height": 0
              },
              "window_rect": {
                "x": 0,
                "y": 0,
                "width": 0,
                "height": 0
              },
              "geometry": {
                "x": 0,
                "y": 0,
                "width": 0,
                "height": 0
              },
              "name": "1",
              "num": 1,
              "gaps": {
                "inner": 0,
                "outer": 0,
                "top": 0,
                "right": 0,
                "bottom": 0,
                "left": 0
              },
              "window": null,
              "window_type": null,
              "nodes": [
                {
                  "id": 94681385599728,
                  "type": "con",
                  "orientation": "none",
                  "scratchpad_state": "none",
                  "percent": 1,
                  "urgent": false,
                  "focused": false,
                  "output": "HDMI1",
                  "layout": "splith",
                  "workspace_layout": "default",
                  "last_split_layout": "splith",
                  "border": "pixel",
                  "current_border_width": 4,
                  "rect": {
                    "x": 0,
                    "y": 0,
                    "width": 1920,
                    "height": 1080
                  },
                  "deco_rect": {
                    "x": 0,
                    "y": 0,
                    "width": 0,
                    "height": 0
                  },
                  "window_rect": {
                    "x": 0,
                    "y": 0,
                    "width": 1920,
                    "height": 1080
                  },
                  "geometry": {
                    "x": 2560,
                    "y": 26,
                    "width": 3840,
                    "height": 2134
                  },
                  "name": "Unknown Song - Unknown Artist",
                  "window": 79691777,
                  "window_type": "normal",
                  "window_properties": {
                    "class": "Google Play Music Desktop Player",
                    "instance": "google play music desktop player",
                    "window_role": "browser-window",
                    "title": "Unknown Song - Unknown Artist",
                    "transient_for": null
                  },
                  "nodes": [],
                  "floating_nodes": [],
                  "focus": [],
                  "fullscreen_mode": 0,
                  "sticky": false,
                  "floating": "auto_off",
                  "swallows": []
                }
              ],
              "floating_nodes": [],
              "focus": [
                94681385599728
              ],
              "fullscreen_mode": 0,
              "sticky": false,
              "floating": "auto_off",
              "swallows": []
            },
            {
              "id": 94681385605536,
              "type": "workspace",
              "orientation": "none",
              "scratchpad_state": "none",
              "percent": 0.5,
              "urgent": false,
              "focused": false,
              "output": "HDMI1",
              "layout": "splith",
              "workspace_layout": "default",
              "last_split_layout": "splith",
              "border": "normal",
              "current_border_width": -1,
              "rect": {
                "x": 0,
                "y": 24,
                "width": 1920,
                "height": 1056
              },
              "deco_rect": {
                "x": 0,
                "y": 0,
                "width": 0,
                "height": 0
              },
              "window_rect": {
                "x": 0,
                "y": 0,
                "width": 0,
                "height": 0
              },
              "geometry": {
                "x": 0,
                "y": 0,
                "width": 0,
                "height": 0
              },
              "name": "2",
              "num": 2,
              "gaps": {
                "inner": 0,
                "outer": 0,
                "top": 0,
                "right": 0,
                "bottom": 0,
                "left": 0
              },
              "window": null,
              "window_type": null,
              "nodes": [],
              "floating_nodes": [],
              "focus": [],
              "fullscreen_mode": 1,
              "sticky": false,
              "floating": "auto_off",
              "swallows": []
            }
          ],
          "floating_nodes": [],
          "focus": [
            94681385605536,
            94681385599216
          ],
          "fullscreen_mode": 0,
          "sticky": false,
          "floating": "auto_off",
          "swallows": []
        },
        {
          "id": 94681385611328,
          "type": "dockarea",
          "orientation": "none",
          "scratchpad_state": "none",
          "percent": 0.3333333333333333,
          "urgent": false,
          "focused": false,
          "output": "HDMI1",
          "layout": "dockarea",
          "workspace_layout": "default",
          "last_split_layout": "splith",
          "border": "normal",
          "current_border_width": -1,
          "rect": {
            "x": 0,
            "y": 1080,
            "width": 1920,
            "height": 0
          },
          "deco_rect": {
            "x": 0,
            "y": 0,
            "width": 0,
            "height": 0
          },
          "window_rect": {
            "x": 0,
            "y": 0,
            "width": 0,
            "height": 0
          },
          "geometry": {
            "x": 0,
            "y": 0,
            "width": 0,
            "height": 0
          },
          "name": "bottomdock",
          "window": null,
          "window_type": null,
          "nodes": [],
          "floating_nodes": [],
          "focus": [],
          "fullscreen_mode": 0,
          "sticky": false,
          "floating": "auto_off",
          "swallows": [
            {
              "dock": 3,
              "insert_where": 2
            }
          ]
        }
      ],
      "floating_nodes": [],
      "focus": [
        94681385598704,
        94681385595552,
        94681385611328
      ],
      "fullscreen_mode": 0,
      "sticky": false,
      "floating": "auto_off",
      "swallows": []
    },
    {
      "id": 94681385617184,
      "type": "output",
      "orientation": "none",
      "scratchpad_state": "none",
      "percent": 0.25,
      "urgent": false,
      "focused": false,
      "layout": "output",
      "workspace_layout": "default",
      "last_split_layout": "splith",
      "border": "normal",
      "current_border_width": -1,
      "rect": {
        "x": 4480,
        "y": 0,
        "width": 1920,
        "height": 1080
      },
      "deco_rect": {
        "x": 0,
        "y": 0,
        "width": 0,
        "height": 0
      },
      "window_rect": {
        "x": 0,
        "y": 0,
        "width": 0,
        "height": 0
      },
      "geometry": {
        "x": 0,
        "y": 0,
        "width": 0,
        "height": 0
      },
      "name": "HDMI-1-3",
      "window": null,
      "window_type": null,
      "nodes": [
        {
          "id": 94681385617696,
          "type": "dockarea",
          "orientation": "none",
          "scratchpad_state": "none",
          "percent": 0.3333333333333333,
          "urgent": false,
          "focused": false,
          "output": "HDMI-1-3",
          "layout": "dockarea",
          "workspace_layout": "default",
          "last_split_layout": "splith",
          "border": "normal",
          "current_border_width": -1,
          "rect": {
            "x": 4480,
            "y": 0,
            "width": 1920,
            "height": 24
          },
          "deco_rect": {
            "x": 0,
            "y": 0,
            "width": 0,
            "height": 0
          },
          "window_rect": {
            "x": 0,
            "y": 0,
            "width": 0,
            "height": 0
          },
          "geometry": {
            "x": 0,
            "y": 0,
            "width": 0,
            "height": 0
          },
          "name": "topdock",
          "window": null,
          "window_type": null,
          "nodes": [
            {
              "id": 94681385979312,
              "type": "con",
              "orientation": "none",
              "scratchpad_state": "none",
              "percent": 1,
              "urgent": false,
              "focused": false,
              "output": "HDMI-1-3",
              "layout": "splith",
              "workspace_layout": "default",
              "last_split_layout": "splith",
              "border": "normal",
              "current_border_width": 2,
              "rect": {
                "x": 4480,
                "y": 0,
                "width": 1920,
                "height": 24
              },
              "deco_rect": {
                "x": 0,
                "y": 0,
                "width": 0,
                "height": 0
              },
              "window_rect": {
                "x": 0,
                "y": 0,
                "width": 1920,
                "height": 24
              },
              "geometry": {
                "x": 4480,
                "y": 1056,
                "width": 1920,
                "height": 24
              },
              "name": "i3bar for output HDMI-1-3",
              "window": 85983238,
              "window_type": "unknown",
              "window_properties": {
                "class": "i3bar",
                "instance": "i3bar",
                "title": "i3bar for output HDMI-1-3",
                "transient_for": null
              },
              "nodes": [],
              "floating_nodes": [],
              "focus": [],
              "fullscreen_mode": 0,
              "sticky": false,
              "floating": "auto_off",
              "swallows": []
            }
          ],
          "floating_nodes": [],
          "focus": [
            94681385979312
          ],
          "fullscreen_mode": 0,
          "sticky": false,
          "floating": "auto_off",
          "swallows": [
            {
              "dock": 2,
              "insert_where": 2
            }
          ]
        },
        {
          "id": 94681385620848,
          "type": "con",
          "orientation": "horizontal",
          "scratchpad_state": "none",
          "percent": 0.3333333333333333,
          "urgent": false,
          "focused": false,
          "output": "HDMI-1-3",
          "layout": "splith",
          "workspace_layout": "default",
          "last_split_layout": "splith",
          "border": "normal",
          "current_border_width": -1,
          "rect": {
            "x": 4480,
            "y": 24,
            "width": 1920,
            "height": 1056
          },
          "deco_rect": {
            "x": 0,
            "y": 0,
            "width": 0,
            "height": 0
          },
          "window_rect": {
            "x": 0,
            "y": 0,
            "width": 0,
            "height": 0
          },
          "geometry": {
            "x": 0,
            "y": 0,
            "width": 0,
            "height": 0
          },
          "name": "content",
          "window": null,
          "window_type": null,
          "nodes": [
            {
              "id": 94681385650768,
              "type": "workspace",
              "orientation": "horizontal",
              "scratchpad_state": "none",
              "percent": 1,
              "urgent": false,
              "focused": false,
              "output": "HDMI-1-3",
              "layout": "splith",
              "workspace_layout": "default",
              "last_split_layout": "splith",
              "border": "normal",
              "current_border_width": -1,
              "rect": {
                "x": 4480,
                "y": 24,
                "width": 1920,
                "height": 1056
              },
              "deco_rect": {
                "x": 0,
                "y": 0,
                "width": 0,
                "height": 0
              },
              "window_rect": {
                "x": 0,
                "y": 0,
                "width": 0,
                "height": 0
              },
              "geometry": {
                "x": 0,
                "y": 0,
                "width": 0,
                "height": 0
              },
              "name": "9",
              "num": 9,
              "gaps": {
                "inner": 0,
                "outer": 0,
                "top": 0,
                "right": 0,
                "bottom": 0,
                "left": 0
              },
              "window": null,
              "window_type": null,
              "nodes": [
                {
                  "id": 94681385647424,
                  "type": "con",
                  "orientation": "none",
                  "scratchpad_state": "none",
                  "percent": 1,
                  "urgent": false,
                  "focused": false,
                  "output": "HDMI-1-3",
                  "layout": "splith",
                  "workspace_layout": "default",
                  "last_split_layout": "splith",
                  "border": "pixel",
                  "current_border_width": 4,
                  "rect": {
                    "x": 4480,
                    "y": 24,
                    "width": 1920,
                    "height": 1056
                  },
                  "deco_rect": {
                    "x": 0,
                    "y": 0,
                    "width": 0,
                    "height": 0
                  },
                  "window_rect": {
                    "x": 0,
                    "y": 0,
                    "width": 1920,
                    "height": 1056
                  },
                  "geometry": {
                    "x": 0,
                    "y": 0,
                    "width": 1257,
                    "height": 1080
                  },
                  "name": "i3: i3-msg(1) - Mozilla Firefox",
                  "window": 27262979,
                  "window_type": "normal",
                  "window_properties": {
                    "class": "firefox",
                    "instance": "Navigator",
                    "window_role": "browser",
                    "title": "i3: i3-msg(1) - Mozilla Firefox",
                    "transient_for": null
                  },
                  "nodes": [],
                  "floating_nodes": [],
                  "focus": [],
                  "fullscreen_mode": 0,
                  "sticky": false,
                  "floating": "auto_off",
                  "swallows": []
                }
              ],
              "floating_nodes": [],
              "focus": [
                94681385647424
              ],
              "fullscreen_mode": 1,
              "sticky": false,
              "floating": "auto_off",
              "swallows": []
            }
          ],
          "floating_nodes": [],
          "focus": [
            94681385650768
          ],
          "fullscreen_mode": 0,
          "sticky": false,
          "floating": "auto_off",
          "swallows": []
        },
        {
          "id": 94681385627152,
          "type": "dockarea",
          "orientation": "none",
          "scratchpad_state": "none",
          "percent": 0.3333333333333333,
          "urgent": false,
          "focused": false,
          "output": "HDMI-1-3",
          "layout": "dockarea",
          "workspace_layout": "default",
          "last_split_layout": "splith",
          "border": "normal",
          "current_border_width": -1,
          "rect": {
            "x": 4480,
            "y": 1080,
            "width": 1920,
            "height": 0
          },
          "deco_rect": {
            "x": 0,
            "y": 0,
            "width": 0,
            "height": 0
          },
          "window_rect": {
            "x": 0,
            "y": 0,
            "width": 0,
            "height": 0
          },
          "geometry": {
            "x": 0,
            "y": 0,
            "width": 0,
            "height": 0
          },
          "name": "bottomdock",
          "window": null,
          "window_type": null,
          "nodes": [],
          "floating_nodes": [],
          "focus": [],
          "fullscreen_mode": 0,
          "sticky": false,
          "floating": "auto_off",
          "swallows": [
            {
              "dock": 3,
              "insert_where": 2
            }
          ]
        }
      ],
      "floating_nodes": [],
      "focus": [
        94681385620848,
        94681385617696,
        94681385627152
      ],
      "fullscreen_mode": 0,
      "sticky": false,
      "floating": "auto_off",
      "swallows": []
    },
    {
      "id": 94681385632944,
      "type": "output",
      "orientation": "none",
      "scratchpad_state": "none",
      "percent": 0.25,
      "urgent": false,
      "focused": false,
      "layout": "output",
      "workspace_layout": "default",
      "last_split_layout": "splith",
      "border": "normal",
      "current_border_width": -1,
      "rect": {
        "x": 1920,
        "y": 0,
        "width": 2560,
        "height": 1440
      },
      "deco_rect": {
        "x": 0,
        "y": 0,
        "width": 0,
        "height": 0
      },
      "window_rect": {
        "x": 0,
        "y": 0,
        "width": 0,
        "height": 0
      },
      "geometry": {
        "x": 0,
        "y": 0,
        "width": 0,
        "height": 0
      },
      "name": "DVI-D-1-1",
      "window": null,
      "window_type": null,
      "nodes": [
        {
          "id": 94681385633456,
          "type": "dockarea",
          "orientation": "none",
          "scratchpad_state": "none",
          "percent": 0.3333333333333333,
          "urgent": false,
          "focused": false,
          "output": "DVI-D-1-1",
          "layout": "dockarea",
          "workspace_layout": "default",
          "last_split_layout": "splith",
          "border": "normal",
          "current_border_width": -1,
          "rect": {
            "x": 1920,
            "y": 0,
            "width": 2560,
            "height": 24
          },
          "deco_rect": {
            "x": 0,
            "y": 0,
            "width": 0,
            "height": 0
          },
          "window_rect": {
            "x": 0,
            "y": 0,
            "width": 0,
            "height": 0
          },
          "geometry": {
            "x": 0,
            "y": 0,
            "width": 0,
            "height": 0
          },
          "name": "topdock",
          "window": null,
          "window_type": null,
          "nodes": [
            {
              "id": 94681385981888,
              "type": "con",
              "orientation": "none",
              "scratchpad_state": "none",
              "percent": 1,
              "urgent": false,
              "focused": false,
              "output": "DVI-D-1-1",
              "layout": "splith",
              "workspace_layout": "default",
              "last_split_layout": "splith",
              "border": "normal",
              "current_border_width": 2,
              "rect": {
                "x": 1920,
                "y": 0,
                "width": 2560,
                "height": 24
              },
              "deco_rect": {
                "x": 0,
                "y": 0,
                "width": 0,
                "height": 0
              },
              "window_rect": {
                "x": 0,
                "y": 0,
                "width": 2560,
                "height": 24
              },
              "geometry": {
                "x": 1920,
                "y": 1416,
                "width": 2560,
                "height": 24
              },
              "name": "i3bar for output DVI-D-1-1",
              "window": 85983250,
              "window_type": "unknown",
              "window_properties": {
                "class": "i3bar",
                "instance": "i3bar",
                "title": "i3bar for output DVI-D-1-1",
                "transient_for": null
              },
              "nodes": [],
              "floating_nodes": [],
              "focus": [],
              "fullscreen_mode": 0,
              "sticky": false,
              "floating": "auto_off",
              "swallows": []
            }
          ],
          "floating_nodes": [],
          "focus": [
            94681385981888
          ],
          "fullscreen_mode": 0,
          "sticky": false,
          "floating": "auto_off",
          "swallows": [
            {
              "dock": 2,
              "insert_where": 2
            }
          ]
        },
        {
          "id": 94681385636608,
          "type": "con",
          "orientation": "horizontal",
          "scratchpad_state": "none",
          "percent": 0.3333333333333333,
          "urgent": false,
          "focused": false,
          "output": "DVI-D-1-1",
          "layout": "splith",
          "workspace_layout": "default",
          "last_split_layout": "splith",
          "border": "normal",
          "current_border_width": -1,
          "rect": {
            "x": 1920,
            "y": 24,
            "width": 2560,
            "height": 1416
          },
          "deco_rect": {
            "x": 0,
            "y": 0,
            "width": 0,
            "height": 0
          },
          "window_rect": {
            "x": 0,
            "y": 0,
            "width": 0,
            "height": 0
          },
          "geometry": {
            "x": 0,
            "y": 0,
            "width": 0,
            "height": 0
          },
          "name": "content",
          "window": null,
          "window_type": null,
          "nodes": [
            {
              "id": 94681385637120,
              "type": "workspace",
              "orientation": "horizontal",
              "scratchpad_state": "none",
              "percent": 1,
              "urgent": false,
              "focused": false,
              "output": "DVI-D-1-1",
              "layout": "splith",
              "workspace_layout": "default",
              "last_split_layout": "splith",
              "border": "normal",
              "current_border_width": -1,
              "rect": {
                "x": 1920,
                "y": 24,
                "width": 2560,
                "height": 1416
              },
              "deco_rect": {
                "x": 0,
                "y": 0,
                "width": 0,
                "height": 0
              },
              "window_rect": {
                "x": 0,
                "y": 0,
                "width": 0,
                "height": 0
              },
              "geometry": {
                "x": 0,
                "y": 0,
                "width": 0,
                "height": 0
              },
              "name": "5",
              "num": 5,
              "gaps": {
                "inner": 0,
                "outer": 0,
                "top": 0,
                "right": 0,
                "bottom": 0,
                "left": 0
              },
              "window": null,
              "window_type": null,
              "nodes": [
                {
                  "id": 94681385637632,
                  "type": "con",
                  "orientation": "none",
                  "scratchpad_state": "none",
                  "percent": 1,
                  "urgent": false,
                  "focused": true,
                  "output": "DVI-D-1-1",
                  "layout": "splith",
                  "workspace_layout": "default",
                  "last_split_layout": "splith",
                  "border": "pixel",
                  "current_border_width": 4,
                  "rect": {
                    "x": 1920,
                    "y": 24,
                    "width": 2560,
                    "height": 1416
                  },
                  "deco_rect": {
                    "x": 0,
                    "y": 0,
                    "width": 0,
                    "height": 0
                  },
                  "window_rect": {
                    "x": 0,
                    "y": 0,
                    "width": 2560,
                    "height": 1416
                  },
                  "geometry": {
                    "x": 1289,
                    "y": 40,
                    "width": 622,
                    "height": 1080
                  },
                  "name": "main.go - i3config - Code - OSS",
                  "window": 60817409,
                  "window_type": "normal",
                  "window_properties": {
                    "class": "code-oss",
                    "instance": "code-oss",
                    "window_role": "browser-window",
                    "title": "main.go - i3config - Code - OSS",
                    "transient_for": null
                  },
                  "nodes": [],
                  "floating_nodes": [],
                  "focus": [],
                  "fullscreen_mode": 0,
                  "sticky": false,
                  "floating": "auto_off",
                  "swallows": []
                }
              ],
              "floating_nodes": [],
              "focus": [
                94681385637632
              ],
              "fullscreen_mode": 1,
              "sticky": false,
              "floating": "auto_off",
              "swallows": []
            },
            {
              "id": 94681385640944,
              "type": "workspace",
              "orientation": "horizontal",
              "scratchpad_state": "none",
              "percent": null,
              "urgent": false,
              "focused": false,
              "output": "DVI-D-1-1",
              "layout": "splith",
              "workspace_layout": "default",
              "last_split_layout": "splith",
              "border": "normal",
              "current_border_width": -1,
              "rect": {
                "x": 1920,
                "y": 24,
                "width": 2560,
                "height": 1416
              },
              "deco_rect": {
                "x": 0,
                "y": 0,
                "width": 0,
                "height": 0
              },
              "window_rect": {
                "x": 0,
                "y": 0,
                "width": 0,
                "height": 0
              },
              "geometry": {
                "x": 0,
                "y": 0,
                "width": 0,
                "height": 0
              },
              "name": "6",
              "num": 6,
              "gaps": {
                "inner": 0,
                "outer": 0,
                "top": 0,
                "right": 0,
                "bottom": 0,
                "left": 0
              },
              "window": null,
              "window_type": null,
              "nodes": [
                {
                  "id": 94681385750288,
                  "type": "con",
                  "orientation": "none",
                  "scratchpad_state": "none",
                  "percent": 1,
                  "urgent": false,
                  "focused": false,
                  "output": "DVI-D-1-1",
                  "layout": "splith",
                  "workspace_layout": "default",
                  "last_split_layout": "splith",
                  "border": "pixel",
                  "current_border_width": 4,
                  "rect": {
                    "x": 1920,
                    "y": 24,
                    "width": 2560,
                    "height": 1416
                  },
                  "deco_rect": {
                    "x": 0,
                    "y": 0,
                    "width": 0,
                    "height": 0
                  },
                  "window_rect": {
                    "x": 0,
                    "y": 0,
                    "width": 2560,
                    "height": 1416
                  },
                  "geometry": {
                    "x": 2688,
                    "y": 336,
                    "width": 1024,
                    "height": 768
                  },
                  "name": "generated.go - nulls - Code - OSS",
                  "window": 60817542,
                  "window_type": "normal",
                  "window_properties": {
                    "class": "code-oss",
                    "instance": "code-oss",
                    "window_role": "browser-window",
                    "title": "generated.go - nulls - Code - OSS",
                    "transient_for": null
                  },
                  "nodes": [],
                  "floating_nodes": [],
                  "focus": [],
                  "fullscreen_mode": 0,
                  "sticky": false,
                  "floating": "auto_off",
                  "swallows": []
                }
              ],
              "floating_nodes": [],
              "focus": [
                94681385750288
              ],
              "fullscreen_mode": 0,
              "sticky": false,
              "floating": "auto_off",
              "swallows": []
            }
          ],
          "floating_nodes": [],
          "focus": [
            94681385637120,
            94681385640944
          ],
          "fullscreen_mode": 0,
          "sticky": false,
          "floating": "auto_off",
          "swallows": []
        },
        {
          "id": 94681385662496,
          "type": "dockarea",
          "orientation": "none",
          "scratchpad_state": "none",
          "percent": 0.3333333333333333,
          "urgent": false,
          "focused": false,
          "output": "DVI-D-1-1",
          "layout": "dockarea",
          "workspace_layout": "default",
          "last_split_layout": "splith",
          "border": "normal",
          "current_border_width": -1,
          "rect": {
            "x": 1920,
            "y": 1440,
            "width": 2560,
            "height": 0
          },
          "deco_rect": {
            "x": 0,
            "y": 0,
            "width": 0,
            "height": 0
          },
          "window_rect": {
            "x": 0,
            "y": 0,
            "width": 0,
            "height": 0
          },
          "geometry": {
            "x": 0,
            "y": 0,
            "width": 0,
            "height": 0
          },
          "name": "bottomdock",
          "window": null,
          "window_type": null,
          "nodes": [],
          "floating_nodes": [],
          "focus": [],
          "fullscreen_mode": 0,
          "sticky": false,
          "floating": "auto_off",
          "swallows": [
            {
              "dock": 3,
              "insert_where": 2
            }
          ]
        }
      ],
      "floating_nodes": [],
      "focus": [
        94681385636608,
        94681385633456,
        94681385662496
      ],
      "fullscreen_mode": 0,
      "sticky": false,
      "floating": "auto_off",
      "swallows": []
    }
  ],
  "floating_nodes": [],
  "focus": [
    94681385632944,
    94681385617184,
    94681385595040,
    94681385585680
  ],
  "fullscreen_mode": 0,
  "sticky": false,
  "floating": "auto_off",
  "swallows": []
}

*/
