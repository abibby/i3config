package i3config

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/martinlindhe/notify"
)

func (c *Config) Run() {
	arg1 := ""
	if len(os.Args) > 1 {
		arg1 = os.Args[1]
	}
	if arg1 == "func" {
		err := c.funcs[os.Args[2]]()
		if err != nil {
			notify.Notify("i3config", "i3 config exec func error", err.Error(), "")
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
	} else {
		c.chords.apply(c)

		dir := path.Dir(c.path)
		err := exec.Command("go", "build", "-o", path.Join(dir, c.binName)).Run()
		if err != nil {
			log.Fatal(err)
		}

		src := c.Generate()
		if arg1 == "" || arg1 == "-" {
			fmt.Print(src)
		} else {
			err := os.WriteFile(arg1, []byte(src), 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
