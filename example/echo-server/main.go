package main

import (
	"log"

	"github.com/kevinma2010/go-rocket/echo"
)

var (
	cfg  *Config
	name = "go-rocket"
)

func init() {
	cfg = new(Config)
	if err := cfg.Initial(); err != nil {
		log.Fatalf("config load failure: %s", err)
	}
}

func main() {
	e := echo.New(cfg.Server)
	log.Println(name, "start at", cfg.Server.Address)
	log.Fatalln(e.Start(cfg.Server.Address).Error())
}
