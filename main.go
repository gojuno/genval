package main

//go:generate genval examples/api

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	outputFilePtr := flag.String("outputFile", "validators_generated.go", "output file name")
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
	insp := newInspector()
	if err := insp.inspect(dir, outputFile); err != nil {
		log.Fatalf("unable to inspect structs for %q: %v", dir, err)
	}

	g := NewGenerator(insp.Structs, insp.PublicValidators)

	if err := g.Generate(dir, outputFile, needCheck); err != nil {
		log.Fatalf("unable to generate validators for %q: %v", dir, err)
	}
}

type inspector struct {
	Structs          []structDef
	PublicValidators map[string]bool
	Enums            map[string][]string
}

func newInspector() *inspector {
	return &inspector{
		PublicValidators: make(map[string]bool),
		Enums:            make(map[string][]string),
	}
}

func (insp *inspector) inspect(dir, outputFile string) error {
	files, err := getFilesForInspect(dir, outputFile)
	if err != nil {
		return err
	}
	for _, f := range files {
		fs := token.NewFileSet()
		parsedFile, err := parser.ParseFile(fs, f, nil, 0)
		if err != nil {
			log.Fatalf("Error parsing file: %s: %s", f, err)
		}
		ast.Walk(insp, parsedFile)
	}
	return nil
}

func getFilesForInspect(dir, outputFile string) ([]string, error) {
	if strings.HasPrefix(dir, "/") {
		dir = "." + dir
	}
	var result []string
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir %s: %s", dir, err)
	}
	for _, f := range files {
		if f.IsDir() || f.Name() == outputFile || strings.HasSuffix(f.Name(), "_test.go") || !strings.HasSuffix(f.Name(), ".go") {
			continue
		}
		result = append(result, dir+"/"+f.Name())
	}
	return result, nil
}

func (insp *inspector) Visit(node ast.Node) ast.Visitor {
	switch spec := node.(type) {
	case *ast.TypeSpec:
		insp.visitStruct(spec)
		return nil
	case *ast.FuncDecl: //To check if Validate() method already exist
		methodName := spec.Name.Name
		if methodName == "Validate" {
			st := spec.Recv.List[0].Type
			if star, ok := st.(*ast.StarExpr); ok {
				if x, ok := star.X.(*ast.Ident); ok {
					insp.PublicValidators[x.Name] = true
				}
			}
		}
		return nil
	case *ast.ValueSpec: //To find all consts and then generate validation with consts
		/*if spec.Names[0].Obj.Kind == ast.Con {
			valueName := spec.Names[0].Name
			if x, ok := spec.Type.(*ast.Ident); ok {
				valueType := x.Name
				if !simpleType[valueType] {
					insp.Enums[valueType] = append(insp.Enums[valueType], valueName)
				}
			}
			return nil
		}*/
	}
	return insp
}

func (insp *inspector) addStruct(s structDef) {
	insp.Structs = append(insp.Structs, s)
}

func (insp *inspector) visitStruct(astTypeSpec *ast.TypeSpec) {
	structName := astTypeSpec.Name.Name
	switch v := astTypeSpec.Type.(type) {

	case *ast.StructType:
		astFields := v.Fields.List

		s := structDef{Name: structName}

		for _, astField := range astFields {
			field := fieldDef{StructName: structName}
			field.Type = s.parseFieldType(astField.Type)
			field.Name = field.parseFieldName(astField.Names)
			field.SetTags(field.parseTags(astField.Tag))
			s.Fields = append(s.Fields, field)
		}
		insp.addStruct(s)
		return

	/*case *ast.Ident:
		if v.Obj != nil {
			if astTypeSpec, ok := v.Obj.Decl.(*ast.TypeSpec); ok {
				structDef := inspectStruct(astTypeSpec)
				if structDef != nil {
					structs[structDef.Name] = structDef
				}
			}
		}
		&aliasStruct{Name: structName, AliasType: v.Name}
		return

	case *ast.SelectorExpr:
		if x, ok := v.X.(*ast.Ident); ok {
			debug("skipping generation validator for alias: %s on %s", structName, x.Name+"."+v.Sel.Name)
		} else {
			log.Fatalf("not ident Type! : %+v: %+v", structName, v)
		}
	*/
	case *ast.MapType: //alias on map
	case *ast.InterfaceType: //alias on interface
	case *ast.ArrayType: //alias on array

	default:
		log.Fatalf("not expected Type for typeSpec: %s, %+v: %T", structName, astTypeSpec, astTypeSpec.Type)
	}
}

func (s *structDef) parseFieldType(t ast.Expr) TypeDef {
	switch v := t.(type) {
	case *ast.Ident:
		simple := getSimpleType(v.Name)
		if simple != nil {
			return simple
		}
		return NewTypeStruct()

	case *ast.SelectorExpr:
		// v.X - contains pkg : if x, ok := v.X.(*ast.Ident); ok { x.Name+"."+v.Sel.Name)}
		return NewTypeStruct()

	case *ast.ArrayType:
		return NewTypeArray(s.parseFieldType(v.Elt))

	case *ast.StarExpr:
		return NewTypePointer(s.parseFieldType(v.X))

	case *ast.InterfaceType:
		return NewTypeInterface()

	case *ast.MapType:
		return NewTypeMap(s.parseFieldType(v.Key), s.parseFieldType(v.Value))

	case *ast.FuncType:
		return NewTypeFunc()
	}
	log.Fatalf("undefined typeField for %s: %+v, %T", s.Name, t, t)
	return nil
}

func (fd *fieldDef) parseFieldName(fieldNames []*ast.Ident) string {
	if len(fieldNames) != 0 {
		return fieldNames[0].Name
	}
	return fd.Name //wrapped struct, fieldName the same as type
}

func getSimpleType(fieldType string) TypeDef {
	switch fieldType {
	case "string":
		return NewTypeString()

	case "int":
		fallthrough
	case "int8":
		fallthrough
	case "int16":
		fallthrough
	case "int32":
		fallthrough
	case "int64":
		fallthrough
	case "uint":
		fallthrough
	case "uint8":
		fallthrough
	case "uint16":
		fallthrough
	case "uint32":
		fallthrough
	case "uint64":
		return NewTypeNumber()

	case "float32":
		fallthrough
	case "float64":
		return NewTypeNumber()

	case "bool":
		return NewTypeBool()
	}
	return nil
}
