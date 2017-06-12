package simple

type User struct {
	Name   string         `json:"name" validate:"min_len=3,max_len=64"`
	Age    uint           `json:"age" validate:"min=18,max=95"`
	Dog    Dog            `json:"dog_pet"`
	Emails map[int]string `validate:"min_items=1,key=[max=3],value=[min_len=5]"`
}

type Dog struct {
	Name string `validate:"min_len=1,max_len=64" json:"name"`
}
