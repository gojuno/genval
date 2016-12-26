package overriding

import "fmt"

//Just additional validation of request added
type Request4 struct {
	Age  Age4
	Some int `validate:"min=3,max=64"`
}

type Age4 struct {
	Value int `validate:"min=3,max=64"`
}

func (r Request4) Validate() error {
	if r.Age.Value == 10 && r.Some == 10 {
		return fmt.Errorf("fields Age and Some can't be 10 at the same time")
	}
	return r.validate()
}
