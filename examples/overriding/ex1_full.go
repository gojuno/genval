package overriding

import "fmt"

//public validator for Request1 will not be generated because of overriding,
//validator for Age1 will not changed.
type Request1 struct {
	Age  Age1
	Some int `validate:"min=3,max=64"`
}

type Age1 struct {
	Value int `validate:"min=3,max=64"`
}

func (r Request1) Validate() error {
	if r.Age.Value < 10 { // 3 was override on 10
		return fmt.Errorf("field Age is less than 10")
	}
	if r.Age.Value > 64 {
		return fmt.Errorf("field Age is more than 64")
	}
	if r.Some < 3 {
		return fmt.Errorf("field Some is less than 3")
	}
	if r.Some > 64 {
		return fmt.Errorf("field Some is more than 64")
	}
	return nil
}
