package simple

type User struct {
	Name string `validate:"min_len=3,max_len=64"`
	Age  uint   `validate:"min=18,max=95"`
	Dog  Dog
}

type Dog struct {
	Name string `validate:"min_len=1,max_len=64"`
}
