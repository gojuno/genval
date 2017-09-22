package overriding

import "fmt"

//public validator for Request2 will be generated and will use ValidateMin10 method for Age2,
//validator for Age2 will not changed.
type Request2 struct {
	Age  Age2 `validate:"func=.ValidateMin10"`
	Some int  `validate:"min=3,max=64"`
}

type Age2 struct {
	Value int `validate:"min=3,max=64"`
}

func (r Age2) ValidateMin10() error {
	if r.Value < 10 { // 3 was override on 10
		return fmt.Errorf("field Age is less than 10")
	}
	if r.Value > 64 {
		return fmt.Errorf("field Age is more than 64")
	}
	return nil
}
