package i3config

import (
	"fmt"
	"strings"
)

type WindowType string

const (
	Normal       WindowType = "normal"
	Dialog       WindowType = "dialog"
	Utility      WindowType = "utility"
	Toolbar      WindowType = "toolbar"
	Splash       WindowType = "splash"
	Menu         WindowType = "menu"
	DropdownMenu WindowType = "dropdown_menu"
	PopupMenu    WindowType = "popup_menu"
	Tooltip      WindowType = "tooltip"
	Notification WindowType = "notification"
)

type Urgent string

const (
	Latest Urgent = "latest"
	Oldest Urgent = "oldest"
)

type Criteria struct {
	Class      string     `i3:"class"`       // Compares the window class (the second part of WM_CLASS). Use the special value __focused__ to match all windows having the same window class as the currently focused window.
	Instance   string     `i3:"instance"`    // Compares the window instance (the first part of WM_CLASS). Use the special value __focused__ to match all windows having the same window instance as the currently focused window.
	WindowRole string     `i3:"window_role"` // Compares the window role (WM_WINDOW_ROLE). Use the special value __focused__ to match all windows having the same window role as the currently focused window.
	WindowType WindowType `i3:"window_type"` // Compare the window type (_NET_WM_WINDOW_TYPE). Possible values are normal, dialog, utility, toolbar, splash, menu, dropdown_menu, popup_menu, tooltip and notification.
	ID         string     `i3:"id"`          // Compares the X11 window ID, which you can get via xwininfo for example.
	Title      string     `i3:"title"`       // Compares the X11 window title (_NET_WM_NAME or WM_NAME as fallback). Use the special value __focused__ to match all windows having the same window title as the currently focused window.
	Urgent     Urgent     `i3:"urgent"`      // Compares the urgent state of the window. Can be "latest" or "oldest". Matches the latest or oldest urgent window, respectively. (The following aliases are also available: newest, last, recent, first)
	Workspace  string     `i3:"workspace"`   // Compares the workspace name of the workspace the window belongs to. Use the special value __focused__ to match all windows in the currently focused workspace.
	ConMark    string     `i3:"con_mark"`    // Compares the marks set for this container, see [vim_like_marks]. A match is made if any of the container’s marks matches the specified mark.
	ConID      string     `i3:"con_id"`      // Compares the i3-internal container ID, which you can get via the IPC interface. Handy for scripting. Use the special value __focused__ to match only the currently focused window.
	Floating   string     `i3:"floating"`    // Only matches floating windows. This criterion requires no value.
	Tiling     string     `i3:"tiling"`      // Only matches tiling windows. This criterion requires no value.

}

func (c *Criteria) String() string {
	ret := []string{}

	EachKey(c, func(key, value string) {
		if value != "" {
			ret = append(ret, key+"="+escapeString(value))
		}
	})

	return strings.Join(ret, " ")
}

func (c *Config) ForWindow(criteria Criteria, border *Command) {
	c.raw(fmt.Sprintf("for_window [%s] %s", criteria.String(), border.Generate()))
}

func (c *Config) FocusFollowsMouse(follow bool) {
	strFollow := "no"
	if follow {
		strFollow = "yes"
	}
	c.raw("focus_follows_mouse " + strFollow)
}

func (c *Config) FloatingModifier(key string) {
	c.raw("floating_modifier " + key)
}
