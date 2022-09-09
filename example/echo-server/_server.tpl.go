package _

var (
	// name server name
	name = "micro-server"
	// port server port
	port = 5051
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
		Name string `binding:"query" validator:"required"`
	}
	// SayHelloReply say hello api reply
	SayHelloApiReply struct {
		// Message say hello reply message
		Message string
	}
)

type (
	UserService interface {
		GetInfo(GetUserInfoServiceArg) GetUserInfoServiceReply
	}
	GetUserInfoServiceArg struct {
		Id int64
	}
	GetUserInfoServiceReply struct {
		Name string
		Age  int32
	}
)
