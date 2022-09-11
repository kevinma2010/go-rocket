package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/kevinma2010/astparser"
	"github.com/urfave/cli/v2"
)

func Initial(c *cli.Context) (*Context, error) {
	return initial(c, TemplateFile, DefaultTemplate)
}

func initial(c *cli.Context, tplFile, defaultTpl string) (*Context, error) {
	absPath, err := filepath.Abs("./")
	if err != nil {
		return nil, errors.New("get project absolute path failure, reason: " + err.Error())
	}
	absPath = strings.ReplaceAll(absPath, `\`, `/`)
	absPath = strings.TrimRight(absPath, "/")

	ctx := &Context{
		Name:         absPath[strings.LastIndex(absPath, "/")+1:],
		AbsPath:      absPath,
		Tpl:          defaultTpl,
		TemplateFile: tplFile,
	}
	ctx.Module = ctx.Name
	ctx.GoPath, _ = os.LookupEnv("GOPATH")
	const Src = "/src/"
	var srcIdx = strings.Index(absPath, Src)
	ctx.IsInGoPath = strings.Contains(ctx.GoPath, absPath[:srcIdx])
	if ctx.IsInGoPath {
		ctx.Module = absPath[srcIdx+len(Src):]
	}

	err = loadTemplate(ctx)
	if err != nil {
		return nil, err
	}
	ctx.TemplateInfo, err = astparser.ParseBytes(ctx.TemplateSource)
	if err != nil {
		return nil, err
	}
	return ctx, nil
}

func loadTemplate(ctx *Context) error {
	var isNew bool
	bs, err := ioutil.ReadFile(ctx.TemplateFile)
	if err != nil {
		if os.IsNotExist(err) {
			bs, err = createDefaultTemplate(ctx)
			if err != nil {
				return err
			}
			isNew = true
		} else if os.IsPermission(err) {
			return fmt.Errorf("read %s %s", ctx.TemplateFile, err.Error())
		} else {
			return err
		}
	}
	ctx.IsNew = isNew
	ctx.TemplateSource = bs
	return nil
}

func createDefaultTemplate(ctx *Context) ([]byte, error) {
	t, err := template.New("create_template").Parse(ctx.Tpl)
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

func createTemplateFile(ctx *Context) error {
	// write template file
	f, err := os.OpenFile(ctx.TemplateFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err = f.Write(ctx.TemplateSource); err != nil {
		return err
	}
	log.Println("gen ", ctx.TemplateFile)
	return nil
}
