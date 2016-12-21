package main

import (
	"fmt"
	"io"
)

func NewTypeNumber() *typeNumber {
	return &typeNumber{}
}

type typeNumber struct {
	Min *string
	Max *string
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

func (t typeNumber) Generate(w io.Writer, cfg GenConfig, suffix, name string) {
	if t.Min != nil {
		cfg.AddImport("fmt")
		fmt.Fprintf(w, "if %s < %s {\n", suffix+name, *t.Min)
		fmt.Fprintf(w, "	return fmt.Errorf(\"field %s is less than %s \" )\n", name, *t.Min)
		fmt.Fprintf(w, "}\n")
	}
	if t.Max != nil {
		cfg.AddImport("fmt")
		fmt.Fprintf(w, "if %s > %s {\n", suffix+name, *t.Max)
		fmt.Fprintf(w, "	return fmt.Errorf(\"field %s is more than %s \" )\n", name, *t.Max)
		fmt.Fprintf(w, "}\n")
	}
}
