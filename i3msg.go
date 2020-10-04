package i3config

import (
	"os/exec"
)

func I3Msg(cmd Command) {
	err := exec.Command("i3-msg", cmd.Generate()).Run()
	if err != nil {
		panic(err)
	}
}
