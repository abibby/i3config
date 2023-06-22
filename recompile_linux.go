package i3config

func (c *Config) RecompileFunc(configPath string) error {
	b, err := exec.Command("go", "run", c.path).Output()
	if err != nil {
		return err
	}
	err = os.WriteFile(configPath, b, 0644)
	if err != nil {
		return err
	}
	err = I3msg(Restart)
	if err != nil {
		return err
	}
	return nil
}
