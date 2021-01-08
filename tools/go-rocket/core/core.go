package core

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	DefaultTpl = `
package _

var (
	name = "{{.Name}}"
)

type (
	GreeterApi interface {
		SayHello(SayHelloArg) SayHelloReply
	}
	SayHelloArg struct {
		Name string
	}
	SayHelloReply struct {
		Message string
	}
)
`
)

type Context struct {
	Name, AbsPath, Module, GoPath string
	IsInGoPath, IsNew             bool
	TplFile                       string
	TplSource                     string
}

func Initial(c *cli.Context) (*Context, error) {
	absPath, err := filepath.Abs("./")
	if err != nil {
		return nil, errors.New("get project absolute path failure, reason: " + err.Error())
	}
	absPath = strings.ReplaceAll(absPath, `\`, `/`)
	absPath = strings.TrimRight(absPath, "/")

	ctx := &Context{
		Name:    absPath[strings.LastIndex(absPath, "/")+1:],
		AbsPath: absPath,
		TplFile: "_tpl.go",
	}
	ctx.Module = ctx.Name
	ctx.GoPath, _ = os.LookupEnv("GOPATH")
	const Src = "/src/"
	var srcIdx = strings.Index(absPath, Src)
	ctx.IsInGoPath = strings.Contains(ctx.GoPath, absPath[:srcIdx])
	if ctx.IsInGoPath {
		ctx.Module = absPath[srcIdx+len(Src):]
	}

	err = loadTpl(c, ctx)
	if err != nil {
		return nil, err
	}
	return ctx, nil
}

func loadTpl(c *cli.Context, ctx *Context) error {
	var isNew bool
	bs, err := ioutil.ReadFile(ctx.TplFile)
	if err != nil {
		if os.IsNotExist(err) {
			bs, err = createDefaultTpl(ctx)
			if err != nil {
				return err
			}
			isNew = true
		} else if os.IsPermission(err) {
			return fmt.Errorf("read %s %s", ctx.TplFile, err.Error())
		} else {
			return err
		}
	}
	ctx.IsNew = isNew
	ctx.TplSource = string(bs)
	return nil
}

func createDefaultTpl(ctx *Context) ([]byte, error) {
	t, err := template.New("create_tpl").Parse(DefaultTpl)
	if err != nil {
		return nil, err
	}
	var tpl bytes.Buffer
	err = t.Execute(&tpl, ctx)
	if err != nil {
		return nil, err
	}
	return tpl.Bytes(), nil
}
