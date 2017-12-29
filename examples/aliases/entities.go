package aliases

import (
	"errors"
	"unicode/utf8"
)

// StringType validator will be overridden by validation tags in-place
type StringType string

func (t StringType) Validate() error {
	if utf8.RuneCountInString(string(t)) < 3 {
		return errors.New("shorter than 3 chars")
	}
	if utf8.RuneCountInString(string(t)) > 64 {
		return errors.New("longer than 64 chars")
	}
	return nil
}

func (t StringType) ValidateNotEmpty() error {
	if t == "" {
		return errors.New("string is empty")
	}
	return nil
}

type IntType int
type FloatType float64
type MapType map[string]int

type User struct {
	FirstName           StringType  `validate:"min_len=2,max_len=15"`
	LastName            string      `validate:"min_len=1,max_len=15"`
	NonEmptyString      StringType  `validate:"func=.ValidateNotEmpty"`
	FamilyMembers       IntType     `validate:"min=1,max=100"`
	SomeFloat           FloatType   `validate:"min=2.55,max=99.99"`
	SomeMap             MapType     `validate:"min_items=2,key=[max_len=64],value=[min=-35,max=34]"`
	SomePointer         *StringType `validate:"not_null,min_len=20,max_len=150"`
	SomePointerNullable *StringType `validate:"nullable"`
}
