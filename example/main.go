package main

import (
	"fmt"

	. "github.com/abibby/i3config"
)

var (
	Nord0  = HexColor("2e3440")
	Nord1  = HexColor("3b4252")
	Nord2  = HexColor("434c5e")
	Nord3  = HexColor("4c566a")
	Nord4  = HexColor("d8dee9")
	Nord5  = HexColor("e5e9f0")
	Nord6  = HexColor("eceff4")
	Nord7  = HexColor("8fbcbb")
	Nord8  = HexColor("88c0d0")
	Nord9  = HexColor("81a1c1")
	Nord10 = HexColor("5e81ac")
	Nord11 = HexColor("bf616a")
	Nord12 = HexColor("d08770")
	Nord13 = HexColor("ebcb8b")
	Nord14 = HexColor("a3be8c")
	Nord15 = HexColor("b48ead")
)

var (
	Background = Nord0
	Foreground = Nord4

	Black   = Nord1
	Red     = Nord11
	Green   = Nord14
	Yellow  = Nord13
	Blue    = Nord9
	Magenta = Nord15
	Cyan    = Nord8
	White   = Nord15
)

func main() {
	// "~/.config/i3/config"
	c := New("/home/adam/Documents/code/i3config/example/")

	term := "alacritty"
	editor := "code"
	// editor := term + " -e nvim"

	c.Set("$mod", "Mod4")

	c.Gaps(Gaps{
		Inner: 10,
		Smart: true,
	})
	c.ForWindow(Criteria{Class: ".*"}, Border(4))
	c.HideEdgeBorders(Both)

	c.Colors(&ColorConfig{
		Focused:         ConstantColorClass(Blue),
		FocusedInactive: ConstantColorClass(Background),
		Unfocused:       ConstantColorClass(Background),
		Urgent:          ConstantColorClass(Red),
		Placeholder:     ConstantColorClass(Green),
		Background:      Background,
	})

	c.FocusFollowsMouse(false)
	c.Font("pango:DejaVu Sans Mono 11")

	c.FloatingModifier("$mod")

	c.BindSym("$mod+Shift+r", c.Recompile(), Restart)
	// c.BindSym("$mod+Shift+r", Exec("make -C ~/.config/i3"), Restart)

	c.BindSym("$mod+Return", Exec(term))

	c.BindSym("$mod+Shift+q", Kill)

	c.BindSym("$mod+r", Exec("rofi -show drun"))

	c.BindSym("$mod+Left", FocusLeft).Alias("$mod+a")
	c.BindSym("$mod+Right", FocusRight).Alias("$mod+d")
	c.BindSym("$mod+Up", FocusUp).Alias("$mod+w")
	c.BindSym("$mod+Down", FocusDown).Alias("$mod+s")

	c.BindSym("$mod+Shift+Left", MoveLeft).Alias("$mod+Shift+a")
	c.BindSym("$mod+Shift+Right", MoveRight).Alias("$mod+Shift+d")
	c.BindSym("$mod+Shift+Up", MoveUp).Alias("$mod+Shift+w")
	c.BindSym("$mod+Shift+Down", MoveDown).Alias("$mod+Shift+s")

	c.BindSym("$mod+h", SplitHorizontal)
	c.BindSym("$mod+v", SplitVertical)

	monitors := [][]string{
		{"eDP-1", "DP-0", "eDP1", "eDP-1"},
		{"DP-1", "DVI-D-0", "eDP1", "eDP-1"},
		{"DP-2", "HDMI-0", "eDP1", "eDP-1"},
	}
	for i := 1; i <= 12; i++ {
		workspaceName := fmt.Sprintf("%d", i)

		c.WorkspaceOutput(workspaceName, monitors[(i-1)/4]...)

		if i <= 10 {
			c.BindSym(fmt.Sprintf("$mod+%d", i%10), Workspace(workspaceName))
			c.BindSym(fmt.Sprintf("$mod+Shift+%d", i%10), MoveContainer(workspaceName))
		}
		c.BindSym(fmt.Sprintf("$mod+F%d", i), Workspace(workspaceName))
		c.BindSym(fmt.Sprintf("$mod+Shift+F%d", i), MoveContainer(workspaceName))
	}

	c.BindSym("$mod+Shift+e", Exec("i3-nagbar -t warning -m 'You pressed the exit shortcut. Do you really want to exit i3? This will end your X session.' -B 'Yes, exit i3' 'i3-msg exit'"))

	c.BindSym("$mod+Ctrl+Up", ResizeGrow(Up, 10), ResizeShrink(Down, 10)).Alias("$mod+Ctrl+w")
	c.BindSym("$mod+Ctrl+Down", ResizeGrow(Down, 10), ResizeShrink(Up, 10)).Alias("$mod+Ctrl+s")
	c.BindSym("$mod+Ctrl+Left", ResizeGrow(Left, 10), ResizeShrink(Right, 10)).Alias("$mod+Ctrl+a")
	c.BindSym("$mod+Ctrl+Right", ResizeGrow(Right, 10), ResizeShrink(Left, 10)).Alias("$mod+Ctrl+d")

	c.Bar(func(bc *BarConfig) {
		bc.Position(Top)
		bc.StatusCommand("$HOME/go/bin/i3gobar")
		bc.TrayOutput("primary")
		bc.Colors(&BarColorConfig{
			Background: Background,
			StatusLine: Foreground,
			Separator:  Foreground,
			FocusedWorkspace: &BarWorkspaceColor{
				Border:     Blue,
				Background: Blue,
				Text:       Background,
			},
			ActiveWorkspace: &BarWorkspaceColor{
				Border:     Blue,
				Background: Background,
				Text:       Foreground,
			},
			InactiveWorkspace: &BarWorkspaceColor{
				Border:     Background,
				Background: Background,
				Text:       Foreground,
			},
			UrgentWorkspace: &BarWorkspaceColor{
				Border:     Red,
				Background: Red,
				Text:       Background,
			},
		})
	})

	c.BindSym("$mod+e", Exec("emoji"))
	c.BindSym("$mod+c", Exec(editor))
	c.Chord("$mod+b", "b", Exec("firefox -P Home"))
	c.Chord("$mod+b", "g", Exec("firefox -P Work"))

	c.BindSym("$mod+x", Exec("nautilus"))
	c.BindSym("$mod+p", Exec("passmenu"))
	c.BindSym("$mod+Shift+p", Exec("maim -s --format=png /dev/stdout | xclip -selection clipboard -t image/png -i"))

	c.BindSym("$mod+u", Exec("cat ~/.config/adam/bookmarks | sort | rofi -dmenu -i -p sites | xargs -r surf"))
	c.BindSym("$mod+Shift+b", Exec("find ~/Pictures/wallpapers -type f | rofi -dmenu -i -p Wallpaper > ~/.config/adam/wallpaper && feh --bg-fill \"$(cat ~/.config/adam/wallpaper)\""))

	c.BindSym("$mod+Shift+l", Exec("find ~/.screenlayout -type f | rofi -dmenu -i -p Layout | xargs -r sh && i3-msg restart"))

	// Music Control
	c.BindSym("$mod+space", Exec("playerctl play-pause"))
	c.BindSym("$mod+comma", Exec("playerctl previous"))
	c.BindSym("$mod+period", Exec("playerctl next"))
	c.BindSym("$mod+minus", Exec("changeVolume 2dB- unmute"))
	c.BindSym("$mod+equal", Exec("changeVolume 2dB+ unmute"))
	c.BindSym("$mod+l", Exec("listAudioOutputs > ~/lao.log 2>&1"))

	c.BindSym("$mod+Shift+minus", Exec("changeBrightness -dec 10"))
	c.BindSym("$mod+Shift+equal", Exec("changeBrightness -inc 10"))

	c.BindSym("$mod+Shift+m", Exec("xmodmap /home/adam/.Xmodmap"))

	c.Run()
}

func quake(c *Config, name, keys, command string) {
	modeName := "quake " + name

	termCommand := "alacritty --class quake_term"

	c.ForWindow(Criteria{Instance: "quake_term"}, FloatingEnabled)
	c.BindSym(keys, Mode(modeName), Exec(termCommand+" -e '"+command+"'"), Mode("default"))
	c.Mode(modeName, func(c *Config) {
		c.BindSym("Escape", Mode("default"), Exec("pkill -f "+termCommand))
	})
}
