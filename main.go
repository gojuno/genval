package main

//go:generate genval examples/simple
//go:generate genval examples/overriding
//go:generate genval examples/complicated

import (
	"flag"
	"log"
	"os"
)

func main() {
	outputFilePtr := flag.String("outputFile", "validators.go", "output file name")
	needValidatableCheckPtr := flag.Bool("needValidatableCheck", true, "check struct on Validatable before calling Validate()")
	excludeRegexp := flag.String("excludeRegexp", `(client\.go|client_mock\.go)`,
		"comma separated list of file names that generator should exclude, e.g. `client.go;client_mock.go`")

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

	cfg := inspectorConfig{
		dir:              dir,
		outputFile:       *outputFilePtr,
		excludeRegexpStr: *excludeRegexp,
	}

	mainLogic(cfg, *needValidatableCheckPtr)
}

func mainLogic(cfg inspectorConfig, needCheck bool) {
	insp := NewInspector()
	if err := insp.Inspect(cfg); err != nil {
		log.Fatalf("unable to inspect structs for %q: %v", cfg.dir, err)
	}

	g := NewGenerator(insp.Result())

	if err := g.Generate(cfg.dir, cfg.outputFile, needCheck); err != nil {
		log.Fatalf("unable to generate validators for %q: %v", cfg.dir, err)
	}
}
