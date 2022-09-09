package main

import "github.com/kevinma2010/go-rocket/echo"

type Config struct {
	Env      string      `yaml:"env" json:"env" xml:"env"`
	LogLevel string      `yaml:"log_level" json:"log_level" xml:"log_level"`
	Server   echo.Config `yaml:"server" json:"server" xml:"server"`
}

func (cfg *Config) Initial() error {
	return nil
}
