package main

import "testing"

func TestMain(t *testing.T) {

	exc := `(client\.go|client_mock\.go)`
	mainLogic(inspectorConfig{"examples/simple", "validators.go", exc}, true)
	mainLogic(inspectorConfig{"examples/complicated", "validators.go", exc}, true)
	mainLogic(inspectorConfig{"examples/overriding", "validators.go", exc}, true)
}
