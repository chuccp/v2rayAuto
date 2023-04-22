package config

import (
	"github.com/go-acme/lego/v4/log"
	"gopkg.in/ini.v1"
)

type Config struct {
	file *ini.File
}

func (c *Config) ReadString(section string, key string) (string, error) {
	log.Println("读取配置文件 section:", section, "key:", key)
	st, err := c.file.GetSection(section)
	if err != nil {
		return "", err
	}
	return st.Key(key).String(), nil
}

func (c *Config) HasSection(section string) bool {
	return c.file.HasSection(section)
}

func (c *Config) ReadInt(section string, key string) (int, error) {
	log.Println("读取配置文件 section:", section, "key:", key)
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
