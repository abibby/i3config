package i3config

type ModeType struct {
	name   string
	config *Config
}

func (c *Config) Mode(name string, mode func(c *Config)) {
	subConfig := c.newSubConfig()
	mode(subConfig)
	c.AddLine(&ModeType{
		name:   name,
		config: subConfig,
	})
}

func (m ModeType) Generate() string {
	return "mode " + escapeString(m.name) + " {\n" + indent(m.config.Generate()) + "\n}"
}
