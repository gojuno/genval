package main

import (
	"fmt"
	"io"
)

func NewArrayType(inner TypeDef) *typeArray {
	return &typeArray{InnerType: inner}
}

type typeArray struct {
	Min       *string
	Max       *string
	InnerType TypeDef
}

func (t typeArray) Type() string {
	return t.InnerType.Type()
}

func (t *typeArray) SetTag(tag Tag) error {
	switch tag.Key() {
	case ArrayMinItemsKey:
		st := tag.(SimpleTag)
		t.Min = &st.Param
	case ArrayMaxItemsKey:
		st := tag.(SimpleTag)
		t.Max = &st.Param
	case ArrayItemKey:
		scope := tag.(ScopeTag)
		for _, it := range scope.InnerTags {
			if err := t.InnerType.SetTag(it); err != nil {
				return fmt.Errorf("set item tags failed for %+v, err: %s", it, err)
			}
		}
	default:
		return ErrUnusedTag
	}
	return nil
}

func (t typeArray) Generate(w io.Writer, cfg GenConfig, name Name) {
	if t.Min != nil {
		if *t.Min != "0" {
			cfg.AddImport("fmt")
			fmt.Fprintf(w, "if len(%s) < %s {\n", name.Full(), *t.Min)
			fmt.Fprintf(w, "    return fmt.Errorf(\"array %s has less items than %s \" )\n", name.FieldName(), *t.Min)
			fmt.Fprintf(w, "}\n")
		}
	}
	if t.Max != nil {
		cfg.AddImport("fmt")
		fmt.Fprintf(w, "if len(%s) > %s {\n", name.Full(), *t.Max)
		fmt.Fprintf(w, "    return fmt.Errorf(\"array %s has more items than %s \" )\n", name.FieldName(), *t.Max)
		fmt.Fprintf(w, "}\n")
	}
	fmt.Fprintf(w, "for _, x := range %s {\n", name.Full())
	fmt.Fprintf(w, "	_ = x \n")
	t.InnerType.Generate(w, cfg, NewSimpleName("x"))
	fmt.Fprintf(w, "}\n")
}

func (t typeArray) Validate() error {
	return validateMinMax(
		t.Min,
		t.Max,
		func(min float64) error {
			if min < 0 {
				return fmt.Errorf("min items can't be less than 0, %f", min)
			}
			return nil
		},
		func(max float64) error {
			if max < 0 {
				return fmt.Errorf("max items can't be less than 0, %f", max)
			}
			return nil
		},
	)
}
