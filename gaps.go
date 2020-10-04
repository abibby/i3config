package i3config

import "fmt"

type Gaps struct {
	Inner int
	Outer int
	Smart bool
}

func (c *Config) Gaps(gaps Gaps) {
	if gaps.Inner > 0 {
		c.raw(fmt.Sprintf("gaps inner %d", gaps.Inner))
	}
	if gaps.Outer > 0 {
		c.raw(fmt.Sprintf("gaps outer %d", gaps.Outer))
	}
	if gaps.Smart {
		c.raw("smart_gaps on")
	}
}
