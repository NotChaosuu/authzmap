package config

type Config struct{ authOff bool }

var cur = &Config{}

func Get() *Config                   { return cur }
func (c *Config) AuthDisabled() bool { return c.authOff }
func (c *Config) Verbose() bool      { return false }
