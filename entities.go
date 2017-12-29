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
	if s.HasOverridenValidation {
		return
	}

	for _, tag := range cfg.SupportedTags {

		varName := "r"
		fmt.Fprintf(w, "// Validate%s validates %s\n", tag, s.Name)
		fmt.Fprintf(w, "func (%s %s) Validate%s() error {\n", varName, s.Name, tag)

		hasAnythingToValidate := (len(s.EnumValues) > 0) || s.aliasType != nil || len(s.fields) > 0 || s.HasAdditionalValidation
		if !hasAnythingToValidate {
			fmt.Fprintf(w, "	return nil\n")
		} else {
			switch {
			case len(s.EnumValues) > 0:
				s.generateEnumValidator(w, cfg, varName)
			case s.aliasType != nil:
				aliasType := *s.aliasType

				cfg.SeveralErrors = aliasType.Type() == types.Map || aliasType.Type() == types.Array
				if cfg.SeveralErrors {
					cfg.AddImport("github.com/gojuno/genval/errlist")
					fmt.Fprint(w, "	var errs errlist.List\n")
				}

				aliasType.Generate(w, cfg, types.NewSimpleNameWithAliasType(varName, aliasType.Type()))

			default:
				cfg.AddImport("github.com/gojuno/genval/errlist")
				fmt.Fprint(w, "	var errs errlist.List\n")

				cfg.SeveralErrors = true

				for _, field := range s.fields {
					field.fieldType.Generate(
						w, cfg, types.NewName(
							"",
							varName+".",
							field.fieldNames.Get(tag),
							field.fieldNames.GetFromStructDefinition(),
							tag))
				}
			}

			if s.HasAdditionalValidation {
				if cfg.SeveralErrors {
					fmt.Fprintf(w, "	errs.Add(r.validate())\n")
				} else {
					fmt.Fprintf(w, "	return r.validate()\n")
				}
			}

			if cfg.SeveralErrors {
				fmt.Fprintf(w, "	return errs.ErrorOrNil()")
			} else {
				fmt.Fprintf(w, "	return nil\n")
			}
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
	fmt.Fprintf(w, "		return fmt.Errorf(\"invalid value for enum %v: %%v\", %s)\n", s.Name, varName)
	fmt.Fprintf(w, "}\n")
}

type FieldDef struct {
	fieldNames types.FieldTagsNames
	fieldType  types.TypeDef
}

func NewField(fieldNames types.FieldTagsNames, fieldType types.TypeDef, validateTags types.ValidatableTags) (*FieldDef, error) {
	fieldType = toPrimitiveType(fieldType, validateTags)
	for _, t := range validateTags {
		if err := fieldType.SetValidateTag(t); err != nil {
			return nil, fmt.Errorf("set validateTags failed, field %s, tag: %+v, err: %s",
				fieldNames.GetFromStructDefinition(), t, err)
		}
	}
	if err := fieldType.Validate(); err != nil {
		return nil, fmt.Errorf("validateTags validation failed for field %s, err: %s",
			fieldNames.GetFromStructDefinition(), err)
	}
	return &FieldDef{
		fieldNames: fieldNames,
		fieldType:  fieldType,
	}, nil
}

func toPrimitiveType(fieldType types.TypeDef, validateTags types.ValidatableTags) types.TypeDef {
	if validateTags.Empty() {
		return fieldType
	}

	if len(validateTags) == 1 &&
		(validateTags.ContainsTag(types.SimpleTag{Name: types.PointerNotNullKey})) ||
		(validateTags.ContainsTag(types.SimpleTag{Name: types.PointerNullableKey})) ||
		(validateTags.ContainsTag(types.SimpleTag{Name: types.StructFuncKey})) {
		return fieldType
	}

	if fieldType.Expr() == nil {
		return fieldType
	}

	underlyingPrimitive := parsePrimitiveType(fieldType.Expr(), fmt.Sprintf("field type %s", fieldType.Type()))
	if underlyingPrimitive != nil {
		// case: pointer to custom type with not_null: ( "Name *SimpleStringType `validate:not_null,min_len=1`")
		p, isPtr := fieldType.(*types.TypePointer)
		if isPtr {
			return p.SetInnerType(underlyingPrimitive)
		}
		return underlyingPrimitive
	}

	return fieldType
}
