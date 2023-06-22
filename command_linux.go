package i3config

func (c *Config) ExecFunc(cb func() error) *Command {
	if c.subConfig {
		panic("ExecFunc must be used from a root config")
	}
	key := fmt.Sprint(funcKey)
	funcKey++
	c.funcs[key] = cb
	dir, file := path.Split(c.path)
	return Exec(fmt.Sprintf(`bash -c "cd '%s' && go run '%s' func %s"`, dir, file, key))
}
