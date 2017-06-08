//This file was automatically generated by the genval generator v1.1
//Please don't modify it manually. Edit your entity tags and then
//run go generate

package complicated

import (
	"fmt"

	"github.com/gojuno/genval/errlist"

	"unicode/utf8"
)

type validatable interface {
	Validate() error
}

func validate(i interface{}) error {
	if v, ok := i.(validatable); ok {
		return v.Validate()
	}
	return nil
}

// Validate validates AliasArray
func (r AliasArray) Validate() error {
	var errs errlist.ErrList
	for _, x := range r {
		_ = x
	}
	return &errs
}

// Validate validates AliasChan
func (r AliasChan) Validate() error {
	var errs errlist.ErrList
	return &errs
}

// Validate validates AliasFunc
func (r AliasFunc) Validate() error {
	var errs errlist.ErrList
	return &errs
}

// Validate validates AliasOnDogsMapAlias
func (r AliasOnDogsMapAlias) Validate() error {
	var errs errlist.ErrList
	if err := DogsMapAlias(r).Validate(); err != nil {
		errs.Add(fmt.Errorf("r is not valid: %v", err))
	}
	return &errs
}

// Validate validates AliasString
func (r AliasString) Validate() error {
	var errs errlist.ErrList
	return &errs
}

// Validate validates Dog
func (r Dog) Validate() error {
	var errs errlist.ErrList
	if utf8.RuneCountInString(r.Name) < 1 {
		errs.Add(fmt.Errorf("field Name is shorter than 1 chars"))
	}
	if utf8.RuneCountInString(r.Name) > 64 {
		errs.Add(fmt.Errorf("field Name is longer than 64 chars"))
	}
	return &errs
}

// Validate validates DogsMapAlias
func (r DogsMapAlias) Validate() error {
	var errs errlist.ErrList
	for k, v := range r {
		_ = k
		_ = v
		if err := v.Validate(); err != nil {
			errs.Add(fmt.Errorf("v is not valid: %v", err))
		}
	}
	return &errs
}

// Validate validates Status
func (r Status) Validate() error {
	var errs errlist.ErrList
	switch r {
	case StatusCreated:
	case StatusPending:
	case StatusActive:
	case StatusFailed:
	default:
		errs.Add(fmt.Errorf("invalid value for enum Status: %v", r))
	}
	return &errs
}

// Validate validates User
func (r User) Validate() error {
	var errs errlist.ErrList
	if utf8.RuneCountInString(r.Name) < 3 {
		errs.Add(fmt.Errorf("field Name is shorter than 3 chars"))
	}
	if utf8.RuneCountInString(r.Name) > 64 {
		errs.Add(fmt.Errorf("field Name is longer than 64 chars"))
	}
	if r.LastName != nil {
		if utf8.RuneCountInString(*r.LastName) < 1 {
			errs.Add(fmt.Errorf("field LastName is shorter than 1 chars"))
		}
		if utf8.RuneCountInString(*r.LastName) > 5 {
			errs.Add(fmt.Errorf("field LastName is longer than 5 chars"))
		}
	}
	if r.Age < 18 {
		errs.Add(fmt.Errorf("field Age is less than 18 "))
	}
	if r.Age > 105 {
		errs.Add(fmt.Errorf("field Age is more than 105 "))
	}
	if r.ChildrenCount == nil {
		errs.Add(fmt.Errorf("field ChildrenCount is required, should not be nil"))
	}
	if *r.ChildrenCount < 0 {
		errs.Add(fmt.Errorf("field ChildrenCount is less than 0 "))
	}
	if *r.ChildrenCount > 15 {
		errs.Add(fmt.Errorf("field ChildrenCount is more than 15 "))
	}
	if r.Float < -4.22 {
		errs.Add(fmt.Errorf("field Float is less than -4.22 "))
	}
	if r.Float > 42.55 {
		errs.Add(fmt.Errorf("field Float is more than 42.55 "))
	}
	if err := r.Dog.Validate(); err != nil {
		errs.Add(fmt.Errorf("Dog is not valid: %v", err))
	}
	if r.DogPointer != nil {
		if err := r.DogPointer.Validate(); err != nil {
			errs.Add(fmt.Errorf("DogPointer is not valid: %v", err))
		}
	}
	if err := r.DogOptional.ValidateOptional(); err != nil {
		errs.Add(fmt.Errorf("DogOptional is not valid: %v", err))
	}
	if len(r.Urls) < 1 {
		errs.Add(fmt.Errorf("array Urls has less items than 1 "))
	}
	for _, x := range r.Urls {
		_ = x
		if utf8.RuneCountInString(x) > 256 {
			errs.Add(fmt.Errorf("field x is longer than 256 chars"))
		}
	}
	if len(r.Dogs) < 1 {
		errs.Add(fmt.Errorf("array Dogs has less items than 1 "))
	}
	for _, x := range r.Dogs {
		_ = x
		if x != nil {
			if err := x.Validate(); err != nil {
				errs.Add(fmt.Errorf("x is not valid: %v", err))
			}
		}
	}
	if r.Test != nil {
		if len(*r.Test) < 1 {
			errs.Add(fmt.Errorf("array Test has less items than 1 "))
		}
		for _, x := range *r.Test {
			_ = x
			if x < 4 {
				errs.Add(fmt.Errorf("field x is less than 4 "))
			}
		}
	}
	if err := validateSome(r.Some); err != nil {
		errs.Add(fmt.Errorf("Some is not valid: %v", err))
	}
	if len(r.SomeArray) < 1 {
		errs.Add(fmt.Errorf("array SomeArray has less items than 1 "))
	}
	for _, x := range r.SomeArray {
		_ = x
		if err := validateSome(x); err != nil {
			errs.Add(fmt.Errorf("x is not valid: %v", err))
		}
	}
	if len(r.Dict) < 2 {
		errs.Add(fmt.Errorf("map Dict has less items than 2 "))
	}
	for k, v := range r.Dict {
		_ = k
		_ = v
		if utf8.RuneCountInString(k) > 64 {
			errs.Add(fmt.Errorf("field k is longer than 64 chars"))
		}
		if v < -35 {
			errs.Add(fmt.Errorf("field v is less than -35 "))
		}
		if v > 34 {
			errs.Add(fmt.Errorf("field v is more than 34 "))
		}
	}
	for k, v := range r.DictDogs {
		_ = k
		_ = v
		if err := v.ValidateOptional(); err != nil {
			errs.Add(fmt.Errorf("v is not valid: %v", err))
		}
		if err := validateMaxDogName(v); err != nil {
			errs.Add(fmt.Errorf("v is not valid: %v", err))
		}
	}
	if err := r.Alias.Validate(); err != nil {
		errs.Add(fmt.Errorf("Alias is not valid: %v", err))
	}
	if err := r.AliasOnAlias.Validate(); err != nil {
		errs.Add(fmt.Errorf("AliasOnAlias is not valid: %v", err))
	}
	if err := r.AliasOnAliasWithCustomValidate.ValidateAlias(); err != nil {
		errs.Add(fmt.Errorf("AliasOnAliasWithCustomValidate is not valid: %v", err))
	}
	for k, v := range r.MapOfMap {
		_ = k
		_ = v
		if len(v) < 1 {
			errs.Add(fmt.Errorf("map v has less items than 1 "))
		}
		for k, v := range v {
			_ = k
			_ = v
			if utf8.RuneCountInString(v) < 3 {
				errs.Add(fmt.Errorf("field v is shorter than 3 chars"))
			}
		}
	}
	for _, x := range r.ByteArray {
		_ = x
	}
	return &errs
}
