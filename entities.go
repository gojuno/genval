package main

import (
	"errors"
	"io"
	"log"
)

type structDef struct {
	Name   string
	Fields []fieldDef
}

func (sd structDef) GenerateBody(w io.Writer, cfg GenConfig) {
	for _, field := range sd.Fields {
		field.Type.Generate(w, cfg, "r.", field.Name)
	}
}

type fieldDef struct {
	StructName string
	Name       string
	Type       TypeDef
}

func (fd *fieldDef) SetTags(tags []Tag) {
	for _, t := range tags {
		if err := fd.Type.SetTag(t); err != nil {
			log.Fatalf("set tags failed for type %s, field %s, tag: %+v, err: %s", fd.StructName, fd.Name, t, err)
		}
	}
}

type TypeDef interface {
	SetTag(Tag) error
	Generate(w io.Writer, cfg GenConfig, suffix, name string)
}

var ErrUnusedTag = errors.New("unused tag")
