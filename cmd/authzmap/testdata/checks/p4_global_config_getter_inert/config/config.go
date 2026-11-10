package config

type Config struct{ v bool }

var cur = &Config{}

func Get() *Config              { return cur }
func (c *Config) Verbose() bool { return c.v }
