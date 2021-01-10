package core

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/structtag"
	"github.com/urfave/cli/v2"
)

const (
	DefaultTpl = `
package main

var (
	// name server name
	name = "{{.Name}}"
	// port server port
	port = 5051
)

type (
	// GreeterApi api group
	GreeterApi interface {
		// SayHello say hello api
		SayHello(int,SayHelloApiArg) SayHelloApiReply
	}
	// SayHelloArg say hello api arg
	SayHelloApiArg struct {
		// Name say hello name
		Name string
	}
	// SayHelloReply say hello api reply
	SayHelloApiReply struct {
		// Message say hello reply message
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

type TplInfo struct {
	Imports    []string
	Structs    []*Struct
	Values     []*Value
	Interfaces []*Interface
}

type Interface struct {
	Name, Doc string
	Methods   []*Function
}

type Function struct {
	Name, Doc string
	Anonymous bool
	Params    []string
	Results   []string
}

type Struct struct {
	Name, Doc string
	Fields    []*StructField
}

type StructField struct {
	Name, Doc string
	Anonymous bool
	Tags      *structtag.Tags
}

type Value struct {
	Name string
	Kind string
	Val  string
	Doc  string
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
		err     error
	)
	var f *ast.File
	if f, err = parser.ParseFile(fileSet, "", ctx.TplSource, parser.ParseComments); err != nil {
		return err
	}

	// get imports code
	if tplInfo.Imports, err = collectImports(ctx, f, fileSet); err != nil {
		return err
	}

	// get structs code
	if tplInfo.Structs, err = collectStruts(ctx, f); err != nil {
		return err
	}

	if tplInfo.Values, err = collectValues(ctx, f); err != nil {
		return err
	}

	if tplInfo.Interfaces, err = collectInterfaces(ctx, f); err != nil {
		return err
	}

	ctx.TplInfo = tplInfo
	return nil
}

func collectImports(ctx *Context, file *ast.File, fileSet *token.FileSet) ([]string, error) {
	var imports []string
	for _, spec := range file.Imports {
		// get code block
		var buffer bytes.Buffer
		if err := format.Node(&buffer, fileSet, spec); err != nil {
			return nil, err
		}
		imports = append(imports, buffer.String())
	}
	return imports, nil
}

func collectStruts(ctx *Context, file *ast.File) ([]*Struct, error) {
	var structs []*Struct

	ast.Inspect(file, func(node ast.Node) bool {
		genDecl, ok := node.(*ast.GenDecl)
		if !ok {
			return true
		}
		for _, spec := range genDecl.Specs {
			var struc = new(Struct)
			switch t := spec.(type) {
			case *ast.TypeSpec:
				if t.Type == nil {
					continue
				}

				// 结构体名称
				struc.Name = t.Name.Name

				// 结构体注释
				if t.Doc != nil {
					struc.Doc = t.Doc.Text()
				}

				// 解析字段列表
				var typ *ast.StructType
				typ, ok = t.Type.(*ast.StructType)
				if !ok {
					continue
				}
				if typ.Fields != nil {
					struc.Fields = collectStructFields(struc.Name, typ.Fields)
				}
			default:
				continue
			}
			structs = append(structs, struc)
		}
		return true
	})

	return structs, nil
}

func collectStructFields(structName string, fields *ast.FieldList) []*StructField {
	var structFields []*StructField
	for _, fieldSpec := range fields.List {
		var structField = new(StructField)
		// 匿名字段
		structField.Anonymous = len(fieldSpec.Names) == 0
		if !structField.Anonymous {
			structField.Name = fieldSpec.Names[0].Name
		}

		// 字段注释
		if fieldSpec.Doc != nil {
			structField.Doc = fieldSpec.Doc.Text()
		}

		// 字段标签
		if fieldSpec.Tag != nil && len(fieldSpec.Tag.Value) > 2 {
			var (
				err      error
				tagValue = fieldSpec.Tag.Value[1 : len(fieldSpec.Tag.Value)-1]
			)
			structField.Tags, err = structtag.Parse(tagValue)
			if err != nil {
				log.Fatalf("parse tag failure, field: [%s %s], tag: [%s], reason: %s",
					structName, structField.Name, tagValue, err)
			}
		}
		structFields = append(structFields, structField)
	}
	return structFields
}

func collectValues(ctx *Context, file *ast.File) ([]*Value, error) {
	var values []*Value

	ast.Inspect(file, func(node ast.Node) bool {
		genDecl, ok := node.(*ast.GenDecl)
		if !ok {
			return true
		}
		for _, spec := range genDecl.Specs {
			var val = new(Value)
			switch t := spec.(type) {
			case *ast.ValueSpec:
				val.Name = t.Names[0].Name

				if t.Doc != nil {
					val.Doc = t.Doc.Text()
				}

				var v *ast.BasicLit
				v, ok = t.Values[0].(*ast.BasicLit)
				if !ok {
					continue
				}
				val.Kind = strings.ToLower(v.Kind.String())
				val.Val = v.Value
				if val.Kind == "string" && len(val.Val) > 2 {
					val.Val = val.Val[1 : len(val.Val)-1]
				}
			default:
				continue
			}
			values = append(values, val)
		}
		return true
	})

	return values, nil
}

func collectInterfaces(ctx *Context, file *ast.File) ([]*Interface, error) {
	var interfaces []*Interface
	ast.Inspect(file, func(node ast.Node) bool {
		genDecl, ok := node.(*ast.GenDecl)
		if !ok {
			return true
		}
		for _, spec := range genDecl.Specs {
			var inf = new(Interface)
			switch t := spec.(type) {
			case *ast.TypeSpec:
				if t.Type == nil {
					continue
				}
				// 接口名称
				inf.Name = t.Name.Name

				// 接口注释
				if t.Doc != nil {
					inf.Doc = t.Doc.Text()
				}

				// 解析函数列表
				var typ *ast.InterfaceType
				typ, ok = t.Type.(*ast.InterfaceType)
				if !ok {
					continue
				}
				if typ.Methods != nil {
					inf.Methods = collectInterfaceMethods(typ.Methods)
				}
			default:
				continue
			}
			interfaces = append(interfaces, inf)
		}
		return true
	})
	return interfaces, nil
}

func collectInterfaceMethods(list *ast.FieldList) []*Function {
	var methods []*Function
	for _, fnSpec := range list.List {
		var fn = new(Function)
		// 匿名函数
		fn.Anonymous = len(fnSpec.Names) == 0
		if !fn.Anonymous {
			fn.Name = fnSpec.Names[0].Name
		}

		// 字段注释
		if fnSpec.Doc != nil {
			fn.Doc = fnSpec.Doc.Text()
		}
		switch f := fnSpec.Type.(type) {
		case *ast.FuncType:
			// 处理 params
			if f.Params != nil {
				fn.Params = collectFuncFields(f.Params)
			}

			// 处理 result
			if f.Results != nil {
				fn.Results = collectFuncFields(f.Results)
			}
		default:
			continue
		}
		methods = append(methods, fn)
	}
	return methods
}

func collectFuncFields(list *ast.FieldList) []string {
	var results []string
	for _, t := range list.List {
		var expr = t.Type
		if e, ok := t.Type.(*ast.StarExpr); ok {
			expr = e.X
		}
		switch ft := expr.(type) {
		case *ast.Ident:
			results = append(results, ft.Name)
		default:
			continue
		}
	}
	return results
}
