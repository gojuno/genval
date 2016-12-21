package main

import (
	"fmt"
	"io"
)

func NewTypeString() *typeString {
	defaultMinLen := "1"
	return &typeString{
		MinLen: &defaultMinLen,
	}
}

type typeString struct {
	MinLen *string
	MaxLen *string
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

func (t typeString) Generate(w io.Writer, cfg GenConfig, suffix, name string) {
	if t.MinLen != nil {
		if *t.MinLen != "0" {
			cfg.AddImport("fmt")
			cfg.AddImport("unicode/utf8")
			fmt.Fprintf(w, "if utf8.RuneCountInString(%s) < %s {\n", suffix+name, *t.MinLen)
			fmt.Fprintf(w, "	return fmt.Errorf(\"field %s is less than %s \" )\n", name, *t.MinLen)
			fmt.Fprintf(w, "}\n")
		}
	}
	if t.MaxLen != nil {
		cfg.AddImport("fmt")
		cfg.AddImport("unicode/utf8")
		fmt.Fprintf(w, "if utf8.RuneCountInString(%s) > %s {\n", suffix+name, *t.MaxLen)
		fmt.Fprintf(w, "	return fmt.Errorf(\"field %s is more than %s \" )\n", name, *t.MaxLen)
		fmt.Fprintf(w, "}\n")
	}
}
