package main

import (
	"fmt"
	"io"

	"github.com/l1va/genval/types"
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

func (s StructDef) GenerateBody(w io.Writer, cfg types.GenConfig, varName string) {
	if len(s.EnumValues) > 0 {
		s.generateEnumValidator(w, cfg, varName)
		return
	}
	for _, field := range s.Fields {
		field.fieldType.Generate(w, cfg, types.NewName("", varName+".", field.fieldName))
	}
}

func (s StructDef) generateEnumValidator(w io.Writer, cfg types.GenConfig, varName string) {
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
	fieldType types.TypeDef
}

func NewField(fieldName string, fieldType types.TypeDef, tags []types.Tag) (*FieldDef, error) {
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
