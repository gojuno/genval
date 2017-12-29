//This file was automatically generated by the genval generator v1.5
//Please don't modify it manually. Edit your entity tags and then
//run go generate

package aliases

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

// Validate validates FloatType
func (r FloatType) Validate() error {
	return nil
}

// Validate validates IntType
func (r IntType) Validate() error {
	return nil
}

// Validate validates MapType
func (r MapType) Validate() error {
	var errs errlist.List
	return errs.ErrorOrNil()
}

// Validate validates User
func (r User) Validate() error {
	var errs errlist.List
	if utf8.RuneCountInString(string(r.FirstName)) < 2 {
		errs.AddFieldf("FirstName", "shorter than 2 chars")
	}
	if utf8.RuneCountInString(string(r.FirstName)) > 15 {
		errs.AddFieldf("FirstName", "longer than 15 chars")
	}
	if utf8.RuneCountInString(string(r.LastName)) < 1 {
		errs.AddFieldf("LastName", "shorter than 1 chars")
	}
	if utf8.RuneCountInString(string(r.LastName)) > 15 {
		errs.AddFieldf("LastName", "longer than 15 chars")
	}
	if err := r.NonEmptyString.ValidateNotEmpty(); err != nil {
		errs.AddField("NonEmptyString", err)
	}
	if r.FamilyMembers < 1 {
		errs.AddFieldf("FamilyMembers", "less than 1")
	}
	if r.FamilyMembers > 100 {
		errs.AddFieldf("FamilyMembers", "more than 100")
	}
	if r.SomeFloat < 2.55 {
		errs.AddFieldf("SomeFloat", "less than 2.55")
	}
	if r.SomeFloat > 99.99 {
		errs.AddFieldf("SomeFloat", "more than 99.99")
	}
	if len(r.SomeMap) < 2 {
		errs.AddFieldf("SomeMap", "less items than 2")
	}
	for kSomeMap, vSomeMap := range r.SomeMap {
		if utf8.RuneCountInString(string(kSomeMap)) > 64 {
			errs.AddFieldf(fmt.Sprintf("SomeMap"+".key[%v]", kSomeMap), "longer than 64 chars")
		}
		if vSomeMap < -35 {
			errs.AddFieldf(fmt.Sprintf("SomeMap"+".%v", kSomeMap), "less than -35")
		}
		if vSomeMap > 34 {
			errs.AddFieldf(fmt.Sprintf("SomeMap"+".%v", kSomeMap), "more than 34")
		}
	}
	if r.SomePointer == nil {
		errs.AddFieldf("SomePointer", "cannot be nil")
	} else {
		if utf8.RuneCountInString(string(*r.SomePointer)) < 20 {
			errs.AddFieldf("SomePointer", "shorter than 20 chars")
		}
		if utf8.RuneCountInString(string(*r.SomePointer)) > 150 {
			errs.AddFieldf("SomePointer", "longer than 150 chars")
		}
	}
	if r.SomePointerNullable != nil {
		if err := r.SomePointerNullable.Validate(); err != nil {
			errs.AddField("SomePointerNullable", err)
		}
	}
	return errs.ErrorOrNil()
}
