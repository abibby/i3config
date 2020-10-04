package i3config

type GenerateFunc func() string

func (c *Config) raw(line string) {
	c.AddLine(GenerateFunc(func() string {
		return line
	}))
}

func (g GenerateFunc) Generate() string {
	return g()
}
