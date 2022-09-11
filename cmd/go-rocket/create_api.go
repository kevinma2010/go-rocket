package main

import "github.com/urfave/cli/v2"

func NewApiGenerator(c *cli.Context, ctx *Context) *ApiGenerator {
	return &ApiGenerator{c, ctx}
}

type ApiGenerator struct {
	c   *cli.Context
	ctx *Context
}

func (g *ApiGenerator) createDirs() error {
	return nil
}

func (g *ApiGenerator) Gen() error {
	return nil
}
