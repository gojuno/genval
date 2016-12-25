package main

import (
	"fmt"
	"io"
)

func NewMapType(key, value TypeDef) *typeMap {
	return &typeMap{Key: key, Value: value}
}

type typeMap struct {
	Min   *string
	Max   *string
	Key   TypeDef
	Value TypeDef
}

func (t typeMap) Type() string {
	return "map"
}

func (t *typeMap) SetTag(tag Tag) error {
	switch tag.Key() {
	case MapMinItemsKey:
		st := tag.(SimpleTag)
		t.Min = &st.Param
	case MapMaxItemsKey:
		st := tag.(SimpleTag)
		t.Max = &st.Param
	case MapKeyKey:
		scope := tag.(ScopeTag)
		for _, it := range scope.InnerTags {
			if err := t.Key.SetTag(it); err != nil {
				return fmt.Errorf("set item tags for key failed, tag %+v, err %s", it, err)
			}
		}
	case MapValueKey:
		scope := tag.(ScopeTag)
		for _, it := range scope.InnerTags {
			if err := t.Value.SetTag(it); err != nil {
				return fmt.Errorf("set item tags for value failed, tag %+v, err %s", it, err)
			}
		}
	default:
		return ErrUnusedTag
	}
	return nil
}

func (t typeMap) Generate(w io.Writer, cfg GenConfig, name Name) {
	if t.Min != nil {
		if *t.Min != "0" {
			cfg.AddImport("fmt")
			fmt.Fprintf(w, "if len(%s) < %s {\n", name.Full(), *t.Min)
			fmt.Fprintf(w, "    return fmt.Errorf(\"map %s has less items than %s \" )\n", name.FieldName(), *t.Min)
			fmt.Fprintf(w, "}\n")
		}
	}
	if t.Max != nil {
		cfg.AddImport("fmt")
		fmt.Fprintf(w, "if len(%s) > %s {\n", name.Full(), *t.Max)
		fmt.Fprintf(w, "    return fmt.Errorf(\"map %s has more items than %s \" )\n", name.FieldName(), *t.Max)
		fmt.Fprintf(w, "}\n")
	}
	fmt.Fprintf(w, "for k, v := range %s {\n", name.Full())
	fmt.Fprintf(w, "	_ = k \n")
	fmt.Fprintf(w, "	_ = v \n")
	t.Key.Generate(w, cfg, NewSimpleName("k"))
	t.Value.Generate(w, cfg, NewSimpleName("v"))
	fmt.Fprintf(w, "}\n")
}

func (t typeMap) Validate() error {
	return validateMinMax(
		t.Min,
		t.Max,
		func(min float64) error {
			if min < 0 {
				return fmt.Errorf("min map items can't be less than 0, %f", min)
			}
			return nil
		},
		func(max float64) error {
			if max < 0 {
				return fmt.Errorf("max map items can't be less than 0, %f", max)
			}
			return nil
		},
	)
}
