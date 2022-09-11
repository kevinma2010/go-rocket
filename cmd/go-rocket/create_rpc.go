package main

import "github.com/urfave/cli/v2"

func NewRpcGenerator(c *cli.Context, ctx *Context) *RpcGenerator {
	return &RpcGenerator{c, ctx}
}

type RpcGenerator struct {
	c   *cli.Context
	ctx *Context
}

func (g *RpcGenerator) createDirs() error {
	return nil
}

func (g *RpcGenerator) Gen() error {
	return nil
}
