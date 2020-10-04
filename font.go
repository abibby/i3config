package i3config

func (c *Config) Font(font string) {
	c.raw("font " + font)
}
