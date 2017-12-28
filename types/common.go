package types

import (
	"errors"
	"fmt"
	"go/ast"
	"io"
	"strconv"
	"strings"
)

type TypeDef interface {
	Type() string
	SetValidateTag(ValidatableTag) error
	Validate() error
	NeedGenerate() bool
	Generate(w io.Writer, cfg GenConfig, name Name)
	Expr() ast.Expr
}

var ErrUnusedTag = errors.New("unused tag")

type GenConfig struct {
	NeedValidatableCheck bool
	SeveralErrors        bool
	SupportedTags        []string
	AddImport            func(string)
}

type Name struct {
	aliasType     *string
	pointerPrefix string
	structVar     string
	fieldName     string
	labelName     string
	tagName       string
}

func (n Name) Full() string {
	return n.pointerPrefix + n.structVar + n.fieldName
}

func (n Name) WithoutPointer() string {
	return n.structVar + n.fieldName
}

func (n Name) WithAlias() string {
	if n.aliasType != nil {
		return *n.aliasType + "(" + n.structVar + n.fieldName + ")"
	}
	return n.structVar + n.fieldName
}

func (n Name) FieldName() string {
	return n.fieldName
}

func (n Name) LabelName() string {
	return n.labelName
}

func NewName(pointerPrefix, structVar, labelName, fieldName, tagName string) Name {
	return Name{
		pointerPrefix: pointerPrefix,
		structVar:     structVar,
		fieldName:     fieldName,
		labelName:     fmt.Sprintf("%q", labelName),
		tagName:       tagName,
	}
}

func NewIndexedName(labelName, indexVar, validateVar, tagName string) Name {
	return Name{
		pointerPrefix: "",
		structVar:     "",
		fieldName:     validateVar,
		labelName:     fmt.Sprintf("fmt.Sprintf(%s + \".%%v\", %v)", labelName, indexVar),
		tagName:       tagName,
	}
}
func NewSimpleNameWithAliasType(fieldName, aliasType string) Name {
	return Name{
		aliasType:     &aliasType,
		pointerPrefix: "",
		structVar:     "",
		fieldName:     fieldName,
		labelName:     fmt.Sprintf("%q", fieldName),
	}
}

func (n Name) WithPointer() Name {
	return Name{
		pointerPrefix: "*",
		structVar:     n.structVar,
		fieldName:     n.fieldName,
		labelName:     n.labelName,
	}
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

func parseFuncsParam(p string) []string {
	r := strings.Split(p, ";")
	var res []string
	for _, v := range r {
		if v != "" {
			res = append(res, v)
		}
	}
	return res
}

func validMaxMin(max, min *string) bool {
	return max != nil || (min != nil && *min != "0")
}
