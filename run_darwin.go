package i3config

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
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
	// if len(os.Args) > 1 && os.Args[1] == "func" {
	// 	err := c.funcs[os.Args[2]]()
	// 	if err != nil {
	// 		exec.Command("notify-send", err.Error()).Run()
	// 		fmt.Printf("%v\n", err)
	// 		os.Exit(1)
	// 	}
	// } else {
	// 	c.chords.apply(c)
	// 	fmt.Print(c.Generate())
	// }
	mainthread.Init(c.initHotKeys)
}
func (c *Config) initHotKeys() {
	close := make(chan struct{})
	hotkeys := []*hotkey.Hotkey{}

	variables := map[string]string{}

	for _, line := range c.lines {
		if v, ok := line.(*Variable); ok {
			variables[v.Name] = v.Value
		}
		if bind, ok := line.(*Bind); ok {
			hotkeys = append(hotkeys, c.registerHotKey(bind, variables))
		}
	}

	<-close
	for _, hk := range hotkeys {
		hk.Unregister()
	}
}

func (c *Config) registerHotKey(b *Bind, variables map[string]string) *hotkey.Hotkey {
	hk := hotkey.New(c.keys(b, variables))
	hk.Register()

	go func() {
		for range hk.Keydown() {
			for _, command := range b.commands {
				if command.name == "exec" {
					argv := parseArgs(unescapeString(command.value))

					b, err := exec.Command("open", append([]string{"-a", argv[0], "-n", "--args"}, argv[1:]...)...).CombinedOutput()
					if string(b) == fmt.Sprintf("Unable to find application named '%s'\n", argv[0]) {
						err = exec.Command(argv[0], argv[1:]...).Run()
					}
					if err != nil {
						log.Print(err)
					}
				}
			}
		}
	}()

	return hk
}

func (c *Config) keys(b *Bind, variables map[string]string) ([]hotkey.Modifier, hotkey.Key) {
	mods := []hotkey.Modifier{}
	key := hotkey.Key(0)

	keys := strings.Split(b.keys, "+")
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

	return mods, key
}
