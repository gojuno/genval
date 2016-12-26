package types

import (
	"errors"
	"fmt"
	"io"
	"strconv"
)

type TypeDef interface {
	Type() string
	SetTag(Tag) error
	Validate() error
	Generate(w io.Writer, cfg GenConfig, name Name)
}

var ErrUnusedTag = errors.New("unused tag")

type GenConfig struct {
	NeedValidatableCheck bool
	AddImport            func(string)
}

type Name struct {
	pointerPrefix string
	structVar     string
	fieldName     string
}

func (n Name) Full() string {
	return n.pointerPrefix + n.structVar + n.fieldName
}

func (n Name) WithoutPointer() string {
	return n.structVar + n.fieldName
}
func (n Name) FieldName() string {
	return n.fieldName
}
func NewName(pointerPrefix, structVar, fieldName string) Name {
	return Name{
		pointerPrefix: pointerPrefix,
		structVar:     structVar,
		fieldName:     fieldName,
	}
}
func NewSimpleName(fieldName string) Name {
	return Name{
		pointerPrefix: "",
		structVar:     "",
		fieldName:     fieldName,
	}
}
func (n Name) WithPointer() Name {
	return Name{
		pointerPrefix: "*",
		structVar:     n.structVar,
		fieldName:     n.fieldName,
	}
}

type Tag interface {
	Key() string
}

type SimpleTag struct {
	Name  string
	Param string
}

func (t SimpleTag) Key() string {
	return t.Name
}

type ScopeTag struct {
	Name      string
	InnerTags []Tag
}

func (t ScopeTag) Key() string {
	return t.Name
}

func validateMinMax(minStr, maxStr *string, minValidate, maxValidate func(float64) error) error {
	var min, max float64
	if minStr != nil {
		f, err := strconv.ParseFloat(*minStr, 64)
		if err != nil {
			return fmt.Errorf("failed to parse value for min tag: %s", *minStr)
		}
		if err := minValidate(f); err != nil {
			return err
		}
		min = f
	}
	if maxStr != nil {
		f, err := strconv.ParseFloat(*maxStr, 64)
		if err != nil {
			return fmt.Errorf("failed to parse value for max tag: %s", *maxStr)
		}
		if err := maxValidate(f); err != nil {
			return err
		}
		max = f
	}
	if minStr != nil && maxStr != nil {
		if max < min {
			return fmt.Errorf("max can't be less than min: min=%s, max=%s", *minStr, *maxStr)
		}
	}
	return nil

}
