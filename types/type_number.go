package types

import (
	"fmt"
	"io"
)

func NewNumber(typeName string) *typeNumber {
	return &typeNumber{typeName: typeName}
}

type typeNumber struct {
	typeName string
	Min      *string
	Max      *string
}

func (t typeNumber) Type() string {
	return t.typeName
}
func (t *typeNumber) SetTag(tag Tag) error {
	switch tag.Key() {
	case NumberMinKey:
		st := tag.(SimpleTag)
		t.Min = &st.Param
	case NumberMaxKey:
		st := tag.(SimpleTag)
		t.Max = &st.Param
	default:
		return ErrUnusedTag
	}
	return nil
}

func (t typeNumber) Generate(w io.Writer, cfg GenConfig, name Name) {
	if t.Min != nil {
		cfg.AddImport("fmt")
		fmt.Fprintf(w, "if %s < %s {\n", name.Full(), *t.Min)
		fmt.Fprintf(w, "	return fmt.Errorf(\"field %s is less than %s \" )\n", name.FieldName(), *t.Min)
		fmt.Fprintf(w, "}\n")
	}
	if t.Max != nil {
		cfg.AddImport("fmt")
		fmt.Fprintf(w, "if %s > %s {\n", name.Full(), *t.Max)
		fmt.Fprintf(w, "	return fmt.Errorf(\"field %s is more than %s \" )\n", name.FieldName(), *t.Max)
		fmt.Fprintf(w, "}\n")
	}
}

func (t typeNumber) Validate() error {
	return validateMinMax(
		t.Min,
		t.Max,
		func(min float64) error {
			return nil
		},
		func(max float64) error {
			return nil
		},
	)
}
