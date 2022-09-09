package _

var mongoList = []interface{}{User}

type (
	// User user model
	User struct {
		Id   int64 `key:"pri"`
		Name string
		Age  int32
	}
)
