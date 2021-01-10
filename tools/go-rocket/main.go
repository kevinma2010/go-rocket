package main

import (
	"fmt"
	"github.com/kevinma2010/go-rocket/tools/go-rocket/gen"
	"log"
	"os"
	"runtime"

	"github.com/kevinma2010/go-rocket/tools/go-rocket/core"
	initial "github.com/kevinma2010/go-rocket/tools/go-rocket/init"
	"github.com/urfave/cli/v2"
)

func main() {
	var buildVersion = "0.1.0"
	app := &cli.App{
		Name:    "go-rocket",
		Usage:   "a cli tool to generate code",
		Version: fmt.Sprintf("%s %s/%s", buildVersion, runtime.GOOS, runtime.GOARCH),
		Authors: []*cli.Author{{
			Name:  "Kevin Ma",
			Email: "mlongbo@gmail.com",
		}},
		Action: func(c *cli.Context) error {
			fmt.Println("Hello friend!")
			return nil
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:  "init",
			Usage: "initial template file",
			Action: func(c *cli.Context) error {
				var (
					serverCtx *core.Context
					modelCtx  *core.Context
					err       error
				)
				if serverCtx, err = core.InitialServer(c); err != nil {
					return err
				}
				if modelCtx, err = core.InitialModel(c); err != nil {
					return err
				}
				return initial.Main(c, serverCtx, modelCtx)
			},
		},
		{
			Name:  "gen",
			Usage: "generate code from template",
			Action: func(c *cli.Context) (err error) {
				var (
					serverCtx *core.Context
					modelCtx  *core.Context
				)
				if serverCtx, err = core.InitialServer(c); err != nil {
					return
				}
				if err = gen.NewServerGenerator(c, serverCtx).Gen(); err != nil {
					return
				}
				if err = gen.NewModelGenerator(c, modelCtx).Gen(); err != nil {
					return
				}
				return
			},
		},
		{
			Name:  "build",
			Usage: "build docker image",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:  "doc",
			Usage: "generate doc files",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
