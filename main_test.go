package main

import "testing"

func TestMain(t *testing.T) {

	exc := `(client\.go|client_mock\.go)`
	mainLogic(config{"examples/simple", "simple", "validators.go", exc}, true)
	mainLogic(config{"examples/complicated", "complicated", "validators.go", exc}, true)
	mainLogic(config{"examples/overriding", "overriding", "validators.go", exc}, true)
}
