package gen

import (
	"github.com/kevinma2010/go-rocket/tools/go-rocket/core"
	"github.com/urfave/cli/v2"
)

type ModelGenerator struct {
	c   *cli.Context
	ctx *core.Context
}

func NewModelGenerator(c *cli.Context, ctx *core.Context) *ModelGenerator {
	return &ModelGenerator{c, ctx}
}

func (g *ModelGenerator) Gen() error {
	if err := g.createDirs(); err != nil {
		return err
	}
	return nil
}

func (g *ModelGenerator) createDirs() error {
	return nil
}
