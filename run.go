package i3config

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	dbus "github.com/godbus/dbus/v5"
)

func (c *Config) Run() {
	arg1 := ""

	err := c.Validate()
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {
		arg1 = os.Args[1]
	}
	if arg1 == "func" {
		err := c.funcs[os.Args[2]]()
		if err != nil {
			notify("i3config", "i3 config exec func error", err.Error(), "")
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

func notify(applicationName, summary, body, icon string) error {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Get the Object identified by path "/org/freedesktop/Notifications"
	// that is defined on peer (could be client or service) "org.freedesktop.Notifications"
	obj := conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")

	// Call method ".Notify"
	// on interface "org.freedesktop.Notifications"
	// which has signature "susssasa{sv}i"
	// and returns "u"
	// Specification: https://specifications.freedesktop.org/notification-spec/notification-spec-latest.html
	call := obj.Call(
		"org.freedesktop.Notifications.Notify",
		// Flags that godbus can pass with the D-Bus message
		0,
		// Optional name of application sending the notification
		applicationName,
		// ID of notification it replaces, or 0 for new notification
		uint(0),
		// Optional program icon
		icon,
		// Summary of the notification
		summary,
		// body
		body,
		// Actions
		make([]string, 0),
		// Hints
		map[string]dbus.Variant{},
		// Expire timeout in milliseconds
		5000,
	)

	return call.Err
}
