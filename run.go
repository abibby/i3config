package i3config

import (
	"fmt"
	"os"
	"os/exec"
)

func (c *Config) Run() {
	if len(os.Args) > 1 && os.Args[1] == "func" {
		err := c.funcs[os.Args[2]]()
		if err != nil {
			exec.Command("notify-send", err.Error()).Run()
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
	} else {
		c.chords.apply(c)
		fmt.Print(c.Generate())
	}
}
