package main

import "testing"

func TestMain(t *testing.T) {

	mainLogic("examples/simple", "validators_generated.go", true)
	mainLogic("examples/complicated", "validators_generated.go", true)
	mainLogic("examples/complicated_without_check", "validators_generated.go", false)
	mainLogic("examples/overriding", "validators_generated.go", true)
}
