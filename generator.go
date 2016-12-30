package main

import (
	"bytes"
	"go/format"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/gojuno/genval/types"
)

type generator struct {
	structs []StructDef
	imports map[string]bool
}

func NewGenerator(structs []StructDef) generator {
	return generator{
		structs: structs,
		imports: make(map[string]bool),
	}
}

func (g generator) Generate(path, outputFile string, needCheck bool) error {
	pkg := filepath.Base(path)
	if pkg == "." {
		pkg = "main"
	}

	cfg := types.GenConfig{
		NeedValidatableCheck: needCheck,
		AddImport: func(imp string) {
			g.imports[imp] = true
		},
	}

	buf, err := g.gen(pkg, cfg)
	if err != nil {
		log.Fatalf("can't generate resulting source: %s", err)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatalf("source: %s\ncan't format resulting source: %s", buf.String(), err)
	}

	filepath := filepath.Join(path, outputFile)

	f, err := os.Create(filepath)
	if err != nil {
		log.Fatalf("can't create file %q: %v", filepath, err)
	}
	defer f.Close()

	if _, err := f.Write(formatted); err != nil {
		log.Fatalf("can't write to resulting file %q: %v", filepath, err)
	}
	return nil
}

func (g generator) gen(pkg string, cfg types.GenConfig) (*bytes.Buffer, error) {
	codeBuf := bytes.NewBuffer([]byte{})
	importsBuf := bytes.NewBuffer([]byte{})

	g.genCode(codeBuf, cfg)
	g.genImports(importsBuf, pkg, cfg.NeedValidatableCheck)

	io.WriteString(importsBuf, codeBuf.String())
	return importsBuf, nil
}

func (g generator) genCode(w io.Writer, cfg types.GenConfig) {
	for _, s := range sorted(g.structs) {
		s.Generate(w, cfg)
	}
}
func sorted(structs []StructDef) []StructDef {
	sort.Sort(AlpabetSorter(structs))
	return structs
}

type AlpabetSorter []StructDef

func (a AlpabetSorter) Len() int           { return len(a) }
func (a AlpabetSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a AlpabetSorter) Less(i, j int) bool { return a[i].Name < a[j].Name }

func (g generator) genImports(w io.Writer, pkg string, needValidatable bool) {
	const t = `
        //This file was automatically generated by the genval generator
        //Please don't modify it manually. Edit your entity tags and then
        //run go generate


        package {{ .Pkg }}
        import (
            {{ range $imp, $v := .Imports }}
                "{{ $imp }}"
            {{ end }}        
        )
	{{if .NeedValidatable}}  
        type Validatable interface {
            Validate() error
        }

        func callValidateIfValidatable(i interface{}) error {
			if v, ok := i.(Validatable); ok {
				if err := v.Validate(); err != nil {
					return err
				}
			}
			return nil
		}
	{{end}}
    `
	type params struct {
		Imports         map[string]bool
		Pkg             string
		NeedValidatable bool
	}

	p := params{
		Imports:         g.imports,
		Pkg:             pkg,
		NeedValidatable: needValidatable,
	}

	tmpl, err := template.New("imports").Parse(t)
	if err != nil {
		log.Fatalf("failed to create template for imports: %v", err)
	}
	if err := tmpl.Execute(w, p); err != nil {
		log.Fatalf("failed to execute template for imports: %v", err)
	}
}
