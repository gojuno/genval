package main

//go:generate ./genval -d examples/aliases -p aliases
//go:generate ./genval -d examples/simple -p simple
//go:generate ./genval -d examples/overriding -p overriding
//go:generate ./genval -d examples/complicated -p complicated
//go:generate ./genval -d examples/empty -p empty

import (
	"flag"
	"log"
	"os"
	"strings"
)

const (
	version = "1.4"
)

var (
	supportedTags        = flag.String("tags", "", "supported tags as source for field naming")
	outputFile           = flag.String("outputFile", "validators.go", "output file name")
	dir                  = flag.String("d", "", "directory with files to be validated")
	pkg                  = flag.String("p", "", "package with files to be validated")
	needValidatableCheck = flag.Bool("needValidatableCheck", true, "check struct on Validatable before calling Validate()")
	excludeRegexp        = flag.String("excludeRegexp", `(client\.go|client_mock\.go)`,
		"regexp file names that generator should exclude, e.g. `(client\\.go|client_mock\\.go)`")
)

func main() {
	flag.Parse()

	// if directory & package aren`t set then first argument is used for both flags
	d, p := *dir, *pkg
	if d == "" && p == "" {
		args := flag.Args()
		if len(args) != 1 {
			flag.PrintDefaults()
			os.Exit(1)
		}

		d, p = args[0], args[0]
	}

	supportedTagsSlice := strings.Split(*supportedTags, ",")
	if len(supportedTagsSlice) >= 1 {
		if supportedTagsSlice[0] != "" {
			supportedTagsSlice = append([]string{""}, supportedTagsSlice...)
		}
	}
	// supportedTags[0] = "" always to generate default validator.
	for i, tag := range supportedTagsSlice {
		supportedTagsSlice[i] = strings.ToUpper(tag)
	}

	cfg := config{
		supportedTags:    supportedTagsSlice,
		dir:              d,
		pkg:              p,
		outputFile:       *outputFile,
		excludeRegexpStr: *excludeRegexp,
	}

	generate(cfg, *needValidatableCheck)
}

func generate(cfg config, needCheck bool) {
	insp := NewInspector()
	if err := insp.Inspect(cfg); err != nil {
		log.Fatalf("unable to inspect structs for %q: %v", cfg.dir, err)
	}

	g := NewGenerator(insp.Result())

	if err := g.Generate(cfg, needCheck); err != nil {
		log.Fatalf("unable to generate validators for %q: %v", cfg.dir, err)
	}
}
