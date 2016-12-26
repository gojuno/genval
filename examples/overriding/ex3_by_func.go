package overriding

import "fmt"

//public validator for Request3 will be generated and will use validateMin10 func for Name,
//validator for Name3 will not changed.
type Request3 struct {
	Age  Age3 `validate:"func=validateMin10"`
	Some int  `validate:"min=3,max=64"`
}

type Age3 struct {
	Value int `validate:"min=3,max=64"`
}

func validateMin10(r Age3) error {
	if r.Value < 10 { // 3 was override on 10
		return fmt.Errorf("field Age is less than 10 ")
	}
	if r.Value > 64 {
		return fmt.Errorf("field Age is more than 64 ")
	}
	return nil
}
