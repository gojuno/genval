package overriding

import "fmt"

//public validator for Request4 will be generated and will use ValidateMin10 method and validateMax128 func for Age4,
//validator for Age4 will not changed.
type Request4 struct {
	Age  Age4 `validate:"func=.ValidateMin10;validateMax128"`
	Some int  `validate:"min=3,max=64"`
}

type Age4 struct {
	Value int `validate:"min=3,max=64"`
}

func (r Age4) ValidateMin10() error {
	if r.Value < 10 { // 3 was override on 10
		return fmt.Errorf("field Age is less than 10 ")
	}
	return nil
}

func validateMax128(r Age4) error {
	if r.Value > 128 { // 64 override on 128
		return fmt.Errorf("field Age is more than 64 ")
	}
	return nil
}
