package i3config

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/abibby/salusa/extra/sets"
)

func (c *Config) Validate() error {
	apps := sets.New[string]()
	for _, line := range c.lines {
		switch line := line.(type) {
		case *Command:
			app := getApplication(line)
			if app != "" {
				apps.Add(app)
			}
		case *Bind:
			for _, c := range line.commands {
				app := getApplication(c)
				if app != "" {
					apps.Add(app)
				}
			}
		}
	}

	missing := []string{}
	for app := range apps.All() {
		err := checkApp(app)
		if err != nil {
			missing = append(missing, app)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing applications: %s", strings.Join(missing, ", "))
	}
	return nil
}

func getApplication(c *Command) string {
	if c.name != "exec" {
		return ""
	}
	command := unescapeString(c.value)
	parts := strings.SplitN(command, " ", 2)
	return parts[0]
}

func checkApp(path string) error {
	if strings.HasPrefix(path, "/") {
		_, err := os.Stat(path)
		return err
	}
	_, err := exec.LookPath(path)
	return err
}
