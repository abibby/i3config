package i3config

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/caseymrm/menuet"
	"golang.design/x/hotkey"
)

func LinesOfType[T Generator](c *Config) []T {
	lines := []T{}
	for _, line := range c.lines {
		if bind, ok := line.(T); ok {
			lines = append(lines, bind)
		}
	}
	return lines
}
func (c *Config) Run() {
	c.chords.apply(c)

	menuet.App().Name = "i3"
	menuet.App().Label = "i3"
	menuet.App().SetMenuState(&menuet.MenuState{
		Title: "i3",
	})

	go func() {
		for {
			spaces, err := yabaiQuerySpaces()
			if err != nil {
				log.Print(err)
				continue
			}
			displays, err := yabaiQueryDisplays()
			if err != nil {
				log.Print(err)
				continue
			}

			sort.Slice(displays, func(i, j int) bool {
				return displays[i].Frame.X < displays[j].Frame.X
			})

			activeSpaces := ""

			for _, d := range displays {
				for _, s := range spaces {
					if s.DisplayIndex == d.Index && s.IsVisible {
						if activeSpaces != "" {
							activeSpaces += " | "
						}
						if s.Label != "" {
							activeSpaces += s.Label
						} else {
							workspace := s.Index
							if workspace <= 4 {
								workspace += 4
							} else if workspace <= 8 {
								workspace -= 4
							}
							activeSpaces += fmt.Sprint(workspace)
						}
					}
				}
			}
			menuet.App().SetMenuState(&menuet.MenuState{
				Title: "i3: " + activeSpaces,
			})
			time.Sleep(time.Second)
		}
	}()

	go c.initHotKeys()

	menuet.App().RunApplication()
}

func (c *Config) initHotKeys() {
	close := make(chan struct{})
	hotkeys := []*hotkey.Hotkey{}
	variables := map[string]string{}

	yabairc := `yabai -m signal --add event=dock_did_restart action="sudo yabai --load-sa"
sudo yabai --load-sa
yabai -m config layout bsp
`

	for _, line := range c.lines {
		if l, ok := line.(GenerateYabaier); ok {
			line := l.GenerateYabai()
			if line != "" {
				yabairc += line + "\n"
			}
		}
		if v, ok := line.(*Variable); ok {
			variables[v.Name] = v.Value
		}
		if bind, ok := line.(*Bind); ok {
			hotkeys = append(hotkeys, c.registerHotKey(bind, variables)...)
		}
	}

	err := os.WriteFile("/Users/adambibby/.config/yabai/yabairc", []byte(yabairc), 0644)
	if err != nil {
		log.Fatal(err)
	}

	exec.Command("yabai", "--restart-service").Run()

	<-close
	for _, hk := range hotkeys {
		hk.Unregister()
	}
}

func (c *Config) registerHotKey(b *Bind, variables map[string]string) []*hotkey.Hotkey {
	hotkeys := make([]*hotkey.Hotkey, len(b.alias)+1)

	for i, keys := range append(b.alias, b.keys) {
		hk := hotkey.New(c.keys(keys, variables))
		hk.Register()
		hotkeys[i] = hk
		go func() {
			for range hk.Keydown() {
				for _, command := range b.commands {
					err := command.RunYabai(c)
					if err != nil {
						log.Print(err)
					}
				}
			}
		}()
	}

	return hotkeys
}

func (c *Config) keys(keysStr string, variables map[string]string) ([]hotkey.Modifier, hotkey.Key) {
	mods := []hotkey.Modifier{}
	key := hotkey.Key(0)

	keys := strings.Split(keysStr, "+")
	for _, k := range keys {
		if strings.HasPrefix(k, "$") {
			k = variables[k]
		}
		if mod, ok := modMap[strings.ToLower(k)]; ok {
			mods = append(mods, mod)
		}
		if ke, ok := keyMap[strings.ToLower(k)]; ok {
			key = ke
		}
	}

	if len(mods) == 1 && mods[0] == hotkey.ModCtrl && (key == hotkey.KeyC || key == hotkey.KeyD) {
		log.Fatal("cannot use ctrl+c as a hotkey")
	}
	return mods, key
}
