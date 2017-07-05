//This file was automatically generated by the genval generator v1.3
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
	var errs errlist.List
	return errs.ErrorOrNil()
}

// Validate validates AliasChan
func (r AliasChan) Validate() error {
	return nil
}

// Validate validates AliasFunc
func (r AliasFunc) Validate() error {
	return nil
}

// Validate validates AliasOnDogsMapAlias
func (r AliasOnDogsMapAlias) Validate() error {
	if err := DogsMapAlias(r).Validate(); err != nil {
		return fmt.Errorf("%s %v", "r", err)
	}
	return nil
}

// Validate validates AliasString
func (r AliasString) Validate() error {
	return nil
}

// Validate validates Dog
func (r Dog) Validate() error {
	var errs errlist.List
	if utf8.RuneCountInString(r.Name) < 1 {
		errs.AddFieldf("Name", "shorter than 1 chars")
	}
	if utf8.RuneCountInString(r.Name) > 64 {
		errs.AddFieldf("Name", "longer than 64 chars")
	}
	return errs.ErrorOrNil()
}

// Validate validates DogsMapAlias
func (r DogsMapAlias) Validate() error {
	var errs errlist.List
	for rKey, rValue := range r {
		_ = rKey
		_ = rValue
		if err := rValue.Validate(); err != nil {
			errs.AddField(fmt.Sprintf("r"+".%v", rKey), err)
		}
	}
	return errs.ErrorOrNil()
}

// Validate validates Status
func (r Status) Validate() error {
	switch r {
	case StatusCreated:
	case StatusPending:
	case StatusActive:
	case StatusFailed:
	default:
		return fmt.Errorf("invalid value for enum Status: %v", r)
	}
	return nil
}

// Validate validates User
func (r User) Validate() error {
	var errs errlist.List
	if utf8.RuneCountInString(r.Name) < 3 {
		errs.AddFieldf("Name", "shorter than 3 chars")
	}
	if utf8.RuneCountInString(r.Name) > 64 {
		errs.AddFieldf("Name", "longer than 64 chars")
	}
	if r.LastName != nil {
		if utf8.RuneCountInString(*r.LastName) < 1 {
			errs.AddFieldf("LastName", "shorter than 1 chars")
		}
		if utf8.RuneCountInString(*r.LastName) > 5 {
			errs.AddFieldf("LastName", "longer than 5 chars")
		}
	}
	if r.Age < 18 {
		errs.AddFieldf("Age", "less than 18")
	}
	if r.Age > 105 {
		errs.AddFieldf("Age", "more than 105")
	}
	if r.ChildrenCount == nil {
		errs.AddFieldf("ChildrenCount", "cannot be nil")
	} else {
		if *r.ChildrenCount < 0 {
			errs.AddFieldf("ChildrenCount", "less than 0")
		}
		if *r.ChildrenCount > 15 {
			errs.AddFieldf("ChildrenCount", "more than 15")
		}
	}
	if r.Float < -4.22 {
		errs.AddFieldf("Float", "less than -4.22")
	}
	if r.Float > 42.55 {
		errs.AddFieldf("Float", "more than 42.55")
	}
	if err := r.Dog.Validate(); err != nil {
		errs.AddField("Dog", err)
	}
	if r.DogPointer != nil {
		if err := r.DogPointer.Validate(); err != nil {
			errs.AddField("DogPointer", err)
		}
	}
	if err := r.DogOptional.ValidateOptional(); err != nil {
		errs.AddField("DogOptional", err)
	}
	if len(r.Urls) < 1 {
		errs.AddFieldf("Urls", "less items than 1")
	}
	for i, x := range r.Urls {
		_ = i
		_ = x
		if utf8.RuneCountInString(x) > 256 {
			errs.AddFieldf(fmt.Sprintf("Urls.%v", i), "longer than 256 chars")
		}
	}
	if len(r.Dogs) < 1 {
		errs.AddFieldf("Dogs", "less items than 1")
	}
	for i, x := range r.Dogs {
		_ = i
		_ = x
		if x != nil {
			if err := x.Validate(); err != nil {
				errs.AddField(fmt.Sprintf("Dogs.%v", i), err)
			}
		}
	}
	if r.Test != nil {
		if len(*r.Test) < 1 {
			errs.AddFieldf("Test", "less items than 1")
		}
		for i, x := range *r.Test {
			_ = i
			_ = x
			if x < 4 {
				errs.AddFieldf(fmt.Sprintf("Test.%v", i), "less than 4")
			}
		}
	}
	if err := validateSome(r.Some); err != nil {
		errs.AddField("Some", err)
	}
	if len(r.SomeArray) < 1 {
		errs.AddFieldf("SomeArray", "less items than 1")
	}
	for i, x := range r.SomeArray {
		_ = i
		_ = x
		if err := validateSome(x); err != nil {
			errs.AddField(fmt.Sprintf("SomeArray.%v", i), err)
		}
	}
	if len(r.Dict) < 2 {
		errs.AddFieldf("Dict", "less items than 2")
	}
	for DictKey, DictValue := range r.Dict {
		_ = DictKey
		_ = DictValue
		if utf8.RuneCountInString(DictKey) > 64 {
			errs.AddFieldf(fmt.Sprintf("Dict"+".%v", DictKey), "longer than 64 chars")
		}
		if DictValue < -35 {
			errs.AddFieldf(fmt.Sprintf("Dict"+".%v", DictKey), "less than -35")
		}
		if DictValue > 34 {
			errs.AddFieldf(fmt.Sprintf("Dict"+".%v", DictKey), "more than 34")
		}
	}
	for DictDogsKey, DictDogsValue := range r.DictDogs {
		_ = DictDogsKey
		_ = DictDogsValue
		if err := DictDogsValue.ValidateOptional(); err != nil {
			errs.AddField(fmt.Sprintf("DictDogs"+".%v", DictDogsKey), err)
		}
		if err := validateMaxDogName(DictDogsValue); err != nil {
			errs.AddField(fmt.Sprintf("DictDogs"+".%v", DictDogsKey), err)
		}
	}
	if err := r.Alias.Validate(); err != nil {
		errs.AddField("Alias", err)
	}
	if err := r.AliasOnAlias.Validate(); err != nil {
		errs.AddField("AliasOnAlias", err)
	}
	if err := r.AliasOnAliasWithCustomValidate.ValidateAlias(); err != nil {
		errs.AddField("AliasOnAliasWithCustomValidate", err)
	}
	for MapOfMapKey, MapOfMapValue := range r.MapOfMap {
		_ = MapOfMapKey
		_ = MapOfMapValue
		if len(MapOfMapValue) < 1 {
			errs.AddFieldf(fmt.Sprintf("MapOfMap"+".%v", MapOfMapKey), "less items than 1")
		}
		for MapOfMapValueKey, MapOfMapValueValue := range MapOfMapValue {
			_ = MapOfMapValueKey
			_ = MapOfMapValueValue
			if utf8.RuneCountInString(MapOfMapValueValue) < 3 {
				errs.AddFieldf(fmt.Sprintf(fmt.Sprintf("MapOfMap"+".%v", MapOfMapKey)+".%v", MapOfMapValueKey), "shorter than 3 chars")
			}
		}
	}
	return errs.ErrorOrNil()
}
