package types

import (
	"fmt"
	"io"
)

func NewString() *typeString {
	return &typeString{}
}

type typeString struct {
	MinLen *string
	MaxLen *string
}

func (t typeString) Type() string {
	return "string"
}

func (t *typeString) SetTag(tag Tag) error {
	switch tag.Key() {
	case StringMinLenKey:
		st := tag.(SimpleTag)
		t.MinLen = &st.Param
	case StringMaxLenKey:
		st := tag.(SimpleTag)
		t.MaxLen = &st.Param
	default:
		return ErrUnusedTag
	}
	return nil
}

func (t typeString) Generate(w io.Writer, cfg GenConfig, name Name) {
	if t.MinLen != nil {
		if *t.MinLen != "0" {
			cfg.AddImport("fmt")
			cfg.AddImport("unicode/utf8")
			fmt.Fprintf(w, "if utf8.RuneCountInString(%s) < %s {\n", name.Full(), *t.MinLen)
			fmt.Fprintf(w, "	return fmt.Errorf(\"field %s is less than %s \" )\n", name.FieldName(), *t.MinLen)
			fmt.Fprintf(w, "}\n")
		}
	}
	if t.MaxLen != nil {
		cfg.AddImport("fmt")
		cfg.AddImport("unicode/utf8")
		fmt.Fprintf(w, "if utf8.RuneCountInString(%s) > %s {\n", name.Full(), *t.MaxLen)
		fmt.Fprintf(w, "	return fmt.Errorf(\"field %s is more than %s \" )\n", name.FieldName(), *t.MaxLen)
		fmt.Fprintf(w, "}\n")
	}
}

func (t typeString) Validate() error {
	return validateMinMax(
		t.MinLen,
		t.MaxLen,
		func(min float64) error {
			if min < 0 {
				return fmt.Errorf("min_len can't be less than 0, %f", min)
			}
			return nil
		},
		func(max float64) error {
			if max < 0 {
				return fmt.Errorf("max_len can't be less than 0, %f", max)
			}
			return nil
		},
	)
}
