package main

const (
	TemplateFile    = "_tpl.go"
	DefaultTemplate = `package _

var (
	// serverName server name
	serverName = "server-name"
	// httpPort http server port
	httpPort = 5051
	// rpcPort rpc server port
	RpcPort = 6061
	// mapToMysql model mapping to mysql
	mapToMysql = []interface{}{User}
)

type (
	// GreeterApi api group
	// @Path("/api/greeter)
	GreeterApi interface {
		// SayHello say hello api
		SayHello(int, SayHelloApiArg) SayHelloApiReply
	}
	// SayHelloArg say hello api arg
	SayHelloApiArg struct {
		// Name say hello name
		Name string ` + "`" + `binding:"query" validator:"required"` + "`" + `
	}
	// SayHelloReply say hello api reply
	SayHelloApiReply struct {
		// Message say hello reply message
		Message string
	}
)

type (
	UserService interface {
		Get(GetUserServiceArg) GetUserServiceReply
	}
	GetUserServiceArg struct {
		Id int64
	}
	GetUserServiceReply struct {
		Value *User
	}
)

type (
	// User user model
	User struct {
		Id   int64 ` + "`" + `key:"pri"` + "`" + `
		Name string
		Age  int32
	}
)
`
)
