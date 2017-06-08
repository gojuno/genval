package types

import (
	"fmt"
	"io"
)

const String string = "string"

func NewString() *typeString {
	return &typeString{}
}

type typeString struct {
	minLen *string
	maxLen *string
}

func (t typeString) Type() string {
	return String
}

func (t *typeString) SetTag(tag Tag) error {
	switch tag.Key() {
	case StringMinLenKey:
		st := tag.(SimpleTag)
		t.minLen = &st.Param
	case StringMaxLenKey:
		st := tag.(SimpleTag)
		t.maxLen = &st.Param
	default:
		return ErrUnusedTag
	}
	return nil
}

func (t typeString) Generate(w io.Writer, cfg GenConfig, name Name) {
	if t.minLen != nil {
		if *t.minLen != "0" {
			cfg.AddImport("fmt")
			cfg.AddImport("unicode/utf8")
			fmt.Fprintf(w, "if utf8.RuneCountInString(%s) < %s {\n", name.Full(), *t.minLen)
			fmt.Fprintf(w, "	   errs.Add(fmt.Errorf(\"field %s is shorter than %s chars\"))\n", name.FieldName(), *t.minLen)
			fmt.Fprintf(w, "}\n")
		}
	}
	if t.maxLen != nil {
		cfg.AddImport("fmt")
		cfg.AddImport("unicode/utf8")
		fmt.Fprintf(w, "if utf8.RuneCountInString(%s) > %s {\n", name.Full(), *t.maxLen)
		fmt.Fprintf(w, "	   errs.Add(fmt.Errorf(\"field %s is longer than %s chars\"))\n", name.FieldName(), *t.maxLen)
		fmt.Fprintf(w, "}\n")
	}
}

func (t typeString) Validate() error {
	return validateMinMax(
		t.minLen,
		t.maxLen,
		func(min float64) error {
			if min < 0 {
				return fmt.Errorf("min_len can't be less than 0: %f", min)
			}
			return nil
		},
		func(max float64) error {
			if max < 0 {
				return fmt.Errorf("max_len can't be less than 0: %f", max)
			}
			return nil
		},
	)
}
