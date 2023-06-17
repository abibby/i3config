package i3config

import (
	"fmt"
	"strings"
)

type Bind struct {
	bindType string
	keys     string
	release  bool
	commands []*Command
	alias    []string
}

func (c *Config) newBind(bindType, keys string, commands []*Command) *Bind {
	b := &Bind{
		bindType: bindType,
		keys:     keys,
		release:  false,
		commands: commands,
		alias:    []string{},
	}

	c.AddLine(b)

	return b
}
func (c *Config) BindSym(keys string, commands ...*Command) *Bind {
	return c.newBind("bindsym", keys, commands)
}

func (c *Config) BindCode(code int, commands ...*Command) *Bind {
	return c.newBind("bindsym", fmt.Sprint(code), commands)
}

func (b *Bind) Release() {
	b.release = true
}

func (b *Bind) Alias(keys string) {
	b.alias = append(b.alias, keys)
}

func (b *Bind) Generate() string {
	release := ""
	if b.release {
		release = "--release "
	}
	strCommands := []string{}

	for _, cmd := range b.commands {
		strCommands = append(strCommands, cmd.Generate())
	}
	src := ""
	for _, keys := range append(b.alias, b.keys) {
		src += b.bindType + " " + release + keys + " " + strings.Join(strCommands, "; ") + "\n"
	}

	return src[:len(src)-1]
}
