package main

//go:generate genval examples/simple
//go:generate genval examples/overriding
//go:generate genval examples/complicated

import (
	"flag"
	"log"
)

func main() {
	outputFilePtr := flag.String("outputFile", "validators.go", "output file name")
	dirPtr := flag.String("d", "api", "directory with files to be validated")
	pkgPtr := flag.String("p", "api", "package with files to be validated")
	needValidatableCheckPtr := flag.Bool("needValidatableCheck", true, "check struct on Validatable before calling Validate()")
	excludeRegexp := flag.String("excludeRegexp", `(client\.go|client_mock\.go)`,
		"regexp file names that generator should exclude, e.g. `(client\\.go|client_mock\\.go)`")

	flag.Parse()

	cfg := config{
		dir:              *dirPtr,
		pkg:              *pkgPtr,
		outputFile:       *outputFilePtr,
		excludeRegexpStr: *excludeRegexp,
	}

	mainLogic(cfg, *needValidatableCheckPtr)
}

func mainLogic(cfg config, needCheck bool) {
	insp := NewInspector()
	if err := insp.Inspect(cfg); err != nil {
		log.Fatalf("unable to inspect structs for %q: %v", cfg.dir, err)
	}

	g := NewGenerator(insp.Result())

	if err := g.Generate(cfg, needCheck); err != nil {
		log.Fatalf("unable to generate validators for %q: %v", cfg.dir, err)
	}
}
