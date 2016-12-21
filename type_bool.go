package main

import "io"

func NewTypeBool() *typeBool {
	return &typeBool{}
}

type typeBool struct {
}

func (t *typeBool) SetTag(tag Tag) error {
	return ErrUnusedTag
}

func (t typeBool) Generate(w io.Writer, cfg GenConfig, suffix, name string) {

}
