package i3config

type Chords map[string][]*BoundCommand
type BoundCommand struct {
	keys     string
	commands []Command
}

func (c *Config) BindChord(key1 string, key2 string, commands ...Command) {

	chords, ok := c.chords[key1]
	if !ok {
		chords = []*BoundCommand{}
	}
	c.chords[key1] = append(chords, &BoundCommand{
		keys:     key2,
		commands: commands,
	})
}

func (ch Chords) apply(c *Config) {
	for key1, commands := range ch {
		chordName := "Chord: " + key1
		c.BindSym(key1, Mode(chordName))
		c.Mode(chordName, func(sub *Config) {
			for _, cmd := range commands {
				sub.BindSym(cmd.keys, append(
					[]Command{Mode("default")},
					cmd.commands...,
				)...)
			}
			sub.BindSym("Escape", Mode("default"))
		})
	}
}
