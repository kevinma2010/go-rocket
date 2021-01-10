package core

const (
	ServerTplFile    = "_server.tpl.go"
	DefaultServerTpl = `package _

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
	ModelTplFile = "_model.tpl.go"
	DefaultModelTpl = `package _
`
)
