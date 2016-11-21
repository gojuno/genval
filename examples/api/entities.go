package api

type User struct {
	Name string `validate:"min=3"`
	Age  uint   `validate:"min=18"`
}

//aliases
//maps
//maps on maps
//invalid tags
//misspelled tags
//array of array
