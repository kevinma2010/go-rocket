package core

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"go/format"
	"go/parser"
	"go/token"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	DefaultTpl = `
package _

import (
	"fmt"
)

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
	TplSource                     []byte
	TplInfo                       *TplInfo
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
	err = parseTpl(ctx)
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
	ctx.TplSource = bs
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

func parseTpl(ctx *Context) error {
	var (
		tplInfo = new(TplInfo)
		fileSet = token.NewFileSet()
	)
	f, err := parser.ParseFile(fileSet, "", ctx.TplSource, parser.ParseComments)
	if err != nil {
		return err
	}

	// get imports code
	for _, spec := range f.Imports {
		// get code block
		var buffer bytes.Buffer
		if err = format.Node(&buffer, fileSet, spec); err != nil {
			return err
		}
		tplInfo.Imports = append(tplInfo.Imports, buffer.String())
	}

	ctx.TplInfo = tplInfo
	return nil
}
