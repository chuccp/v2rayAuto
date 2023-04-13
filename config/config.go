package config

import "gopkg.in/ini.v1"

type Config struct {
	file *ini.File
}

func (c *Config) ReadString(section string, key string) (string, error) {
	st, err := c.file.GetSection(section)
	if err != nil {
		return "", err
	}
	return st.Key(key).String(), nil
}
func (c *Config) ReadInt(section string, key string) (int, error) {
	st, err := c.file.GetSection(section)
	if err != nil {
		return 0, err
	}
	return st.Key(key).Int()
}

func ReadConfig(path string) (*Config, error) {
	cfg, err := ini.Load(path)
	if err != nil {
		return nil, err
	}
	return &Config{file: cfg}, nil
}
