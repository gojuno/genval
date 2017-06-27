package main

import "testing"

func Test_Main(m *testing.T) {

	exc := `(client\.go|client_mock\.go)`
	generate(config{[]string{""}, "examples/simple", "simple", "validators.go", exc}, true)
	generate(config{[]string{""}, "examples/complicated", "complicated", "validators.go", exc}, true)
	generate(config{[]string{""}, "examples/overriding", "overriding", "validators.go", exc}, true)
}
