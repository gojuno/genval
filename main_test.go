package main

import "testing"

func Test_Main(m *testing.T) {

	exc := `(client\.go|client_mock\.go)`
	generate(config{"examples/simple", "simple", "validators.go", exc}, true)
	generate(config{"examples/complicated", "complicated", "validators.go", exc}, true)
	generate(config{"examples/overriding", "overriding", "validators.go", exc}, true)
}
