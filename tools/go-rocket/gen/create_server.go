package gen

import (
	"github.com/kevinma2010/go-rocket/tools/go-rocket/core"
	"github.com/urfave/cli/v2"
)

type ServerGenerator struct {
	c   *cli.Context
	ctx *core.Context
}

func NewServerGenerator(c *cli.Context, ctx *core.Context) *ServerGenerator {
	return &ServerGenerator{c, ctx}
}

func (g *ServerGenerator) Gen() error {
	if err := g.createDirs(); err != nil {
		return err
	}
	return nil
}

func (g *ServerGenerator) createDirs() error {
	return nil
}