package main

import (
	"fmt"
	"io"
)

func NewTypeStruct() *typeStruct {
	return &typeStruct{}
}

type typeStruct struct {
	Method *string
}

func (t *typeStruct) SetTag(tag Tag) error {
	switch tag.Key() {
	case StructMethodKey:
		st := tag.(SimpleTag)
		t.Method = &st.Param
	default:
		return ErrUnusedTag
	}
	return nil
}

func (t typeStruct) Generate(w io.Writer, cfg GenConfig, suffix, name string) {
	switch {
	case t.Method != nil:
		fmt.Fprintf(w, "if err := %s.%s(); err != nil {\n", suffix+name, *t.Method)
		fmt.Fprintf(w, "    return err\n")
		fmt.Fprintf(w, "}\n")
	case cfg.NeedValidatableCheck:
		fmt.Fprintf(w, "if err := callValidateIfValidatable(%s); err != nil {\n", suffix+name)
		fmt.Fprintf(w, "    return err\n")
		fmt.Fprintf(w, "}\n")
	default:
		fmt.Fprintf(w, "if err := %s.Validate(); err != nil {\n", suffix+name)
		fmt.Fprintf(w, "    return err\n")
		fmt.Fprintf(w, "}\n")
	}
}
