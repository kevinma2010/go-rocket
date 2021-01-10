package init

import (
	"github.com/kevinma2010/go-rocket/tools/go-rocket/core"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func Main(c *cli.Context, ctx *core.Context) error {
	// write template file
	f, err := os.OpenFile(ctx.TplFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err = f.Write(ctx.TplSource); err != nil {
		return err
	}
	log.Println("gen ", ctx.TplFile)
	return nil
}
