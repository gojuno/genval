package main

import (
	"fmt"
	"io"

	"github.com/gojuno/genval/types"
)

type StructDef struct {
	Name                    string
	HasOverridenValidation  bool
	HasAdditionalValidation bool
	EnumValues              []string
	fields                  []FieldDef
	aliasType               *types.TypeDef
}

func NewFieldsStruct(name string) StructDef {
	return StructDef{Name: name}
}
func NewAliasStruct(name string, aliasType types.TypeDef) StructDef {
	return StructDef{Name: name, aliasType: &aliasType}
}

func (s *StructDef) AddField(f FieldDef) {
	s.fields = append(s.fields, f)
}

func (s StructDef) Generate(w io.Writer, cfg types.GenConfig) {
	if !s.HasOverridenValidation {
		varName := "r"
		fmt.Fprintf(w, "// Validate validates %s\n", s.Name)
		fmt.Fprintf(w, "func (%s %s) Validate() error {\n", varName, s.Name)
		switch {
		case len(s.EnumValues) > 0:
			s.generateEnumValidator(w, cfg, varName)
		case s.aliasType != nil:
			aliasType := *s.aliasType
			aliasType.Generate(w, cfg, types.NewSimpleNameWithAliasType(varName, aliasType.Type()))
		default:
			for _, field := range s.fields {
				field.fieldType.Generate(w, cfg, types.NewName("", varName+".", field.fieldName))
			}
		}

		if s.HasAdditionalValidation {
			fmt.Fprint(w, "    return r.validate()")
		} else {
			fmt.Fprintf(w, "	return nil\n")
		}

		fmt.Fprintf(w, "}\n\n")
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
