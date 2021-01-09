package core

import "github.com/fatih/structtag"

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
