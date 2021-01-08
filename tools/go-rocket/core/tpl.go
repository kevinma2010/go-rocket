package core

type TplInfo struct {
	Imports []string
	Structs []*Struct
	Vars    []*Var
}

type Struct struct {
}

type Var struct {
	Name    string
	Comment string
}
