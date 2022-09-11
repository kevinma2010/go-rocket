package main

import (
	"github.com/urfave/cli/v2"
)

type EchoServerGenerator struct {
	c   *cli.Context
	ctx *Context
}

func NewEchoServerGenerator(c *cli.Context, ctx *Context) *EchoServerGenerator {
	return &EchoServerGenerator{c, ctx}
}

func (g *EchoServerGenerator) Gen() error {
	if err := g.createDirs(); err != nil {
		return err
	}
	return nil
}

func (g *EchoServerGenerator) createDirs() error {
	return nil
}

type RpcServerGenerator struct {
	c   *cli.Context
	ctx *Context
}

func NewRpcServerGenerator(c *cli.Context, ctx *Context) *RpcServerGenerator {
	return &RpcServerGenerator{c, ctx}
}

func (g *RpcServerGenerator) Gen() error {
	if err := g.createDirs(); err != nil {
		return err
	}
	return nil
}

func (g *RpcServerGenerator) createDirs() error {
	return nil
}
