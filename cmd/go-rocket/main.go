package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/kevinma2010/astparser"
	"github.com/urfave/cli/v2"
)

type Context struct {
	Name, AbsPath, Module, GoPath string
	IsInGoPath, IsNew             bool
	Tpl                           string
	TemplateFile                  string
	TemplateSource                []byte
	TemplateInfo                  *astparser.Parser
}

type Generator interface {
	Gen() error
}

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
				ctx, err := Initial(c)
				if err != nil {
					return err
				}
				return createTemplateFile(ctx)
			},
		},
		{
			Name:  "create",
			Usage: "create project with server",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "type",
					Aliases: []string{"t"},
					Usage:   "create for server type. option: echo-http, rpc",
				},
			},
			Action: func(c *cli.Context) error {
				ctx, err := Initial(c)
				if err != nil {
					return err
				}
				switch c.String("type") {
				case "echo-http":
					return NewEchoServerGenerator(c, ctx).Gen()
				case "rpc":
					return NewRpcServerGenerator(c, ctx).Gen()
				case "":
					return errors.New("type can not be empty")
				default:
					return errors.New("unknown server type")
				}
			},
		},
		{
			Name:  "gen-api",
			Usage: "generate http api code from template",
			Action: func(c *cli.Context) error {
				ctx, err := Initial(c)
				if err != nil {
					return err
				}
				return NewApiGenerator(c, ctx).Gen()
			},
		},
		{
			Name:  "gen-rpc",
			Usage: "generate grpc code from template",
			Action: func(c *cli.Context) error {
				ctx, err := Initial(c)
				if err != nil {
					return err
				}
				return NewRpcGenerator(c, ctx).Gen()
			},
		},
		{
			Name:  "gen-doc",
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
