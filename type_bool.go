package main

import "io"

func NewBoolType() *typeBool {
	return &typeBool{}
}

type typeBool struct {
}

func (t typeBool) Type() string {
	return "bool"
}

func (t *typeBool) SetTag(tag Tag) error {
	return ErrUnusedTag
}

func (t typeBool) Generate(w io.Writer, cfg GenConfig, name Name) {

}

func (t typeBool) Validate() error {
	return nil
}
