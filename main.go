package main

//go:generate genval examples/simple
//go:generate genval examples/overriding
//go:generate genval examples/complicated

import (
	"flag"
	"log"
	"os"
)

const (
	version = "1.0"
)

func main() {
	outputFilePtr := flag.String("outputFile", "validators.go", "output file name")
	needValidatableCheckPtr := flag.Bool("needValidatableCheck", true, "check struct on Validatable before calling Validate()")

	flag.Parse()

	args := flag.Args()
	if len(args) > 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	dir := "api"
	if len(args) == 1 {
		dir = args[0]
	}

	mainLogic(dir, *outputFilePtr, *needValidatableCheckPtr)

}

func mainLogic(dir string, outputFile string, needCheck bool) {
	insp := NewInspector()
	if err := insp.Inspect(dir, outputFile); err != nil {
		log.Fatalf("unable to inspect structs for %q: %v", dir, err)
	}

	g := NewGenerator(insp.Result())

	if err := g.Generate(dir, outputFile, needCheck); err != nil {
		log.Fatalf("unable to generate validators for %q: %v", dir, err)
	}
}
