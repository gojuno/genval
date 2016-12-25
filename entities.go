package main

import (
	"errors"
	"fmt"
	"io"
)

type StructDef struct {
	Name                 string
	Fields               []FieldDef
	PublicValidatorExist bool
	EnumValues           []string
}

func NewStruct(name string) StructDef {
	return StructDef{Name: name}
}

func (s *StructDef) AddField(f FieldDef) {
	s.Fields = append(s.Fields, f)
}

func (s StructDef) GenerateBody(w io.Writer, cfg GenConfig, varName string) {
	if len(s.EnumValues) > 0 {
		s.generateEnumValidator(w, cfg, varName)
		return
	}
	for _, field := range s.Fields {
		field.fieldType.Generate(w, cfg, NewName("", varName+".", field.fieldName))
	}
}

func (s StructDef) generateEnumValidator(w io.Writer, cfg GenConfig, varName string) {
	cfg.AddImport("fmt")

	fmt.Fprintf(w, "switch %s {\n", varName)
	for _, v := range s.EnumValues {
		fmt.Fprintf(w, "case %v: \n", v)
	}
	fmt.Fprintf(w, "	default: \n")
	fmt.Fprintf(w, "		return fmt.Errorf(\"invalid value for enum %v: %%v\", %s) \n", s.Name, varName)
	fmt.Fprintf(w, "}\n")
}

type FieldDef struct {
	fieldName string
	fieldType TypeDef
}

func NewField(fieldName string, fieldType TypeDef, tags []Tag) (*FieldDef, error) {
	for _, t := range tags {
		if err := fieldType.SetTag(t); err != nil {
			return nil, fmt.Errorf("set tags failed, field %s, tag: %+v, err: %s", fieldName, t, err)
		}
	}
	if err := fieldType.Validate(); err != nil {
		return nil, fmt.Errorf("tags validation failed for field %s, err: %s", fieldName, err)
	}
	return &FieldDef{
		fieldName: fieldName,
		fieldType: fieldType,
	}, nil
}

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
