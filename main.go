package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

/*type Inner struct {
    StartDate time.Time
}

type Outer struct {
    InnerStructField *Inner
    CreatedAt time.Time      `validate:"ltecsfield=InnerStructField.StartDate"`
}*/

func main() {
	debugFlagPtr := flag.Bool("debug", false, "for debug logging")
	outputFilePtr := flag.String("outputFile", "validators_generated.go", "output file name")

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

	mainLogic(dir, *debugFlagPtr, *outputFilePtr)

}

func mainLogic(dir string, debug bool, outputFile string) {
	insp := newInspector(debug)
	structs, err := insp.inspect(dir, outputFile)
	if err != nil {
		log.Fatalf("unable to inspect structs for %q: %v", dir, err)
	}

	gen := generator{
		Debug:   debug,
		Structs: structs,
	}

	if err := gen.generate(dir); err != nil {
		log.Fatalf("unable to generate validators for %q: %v", dir, err)
	}
}

type inspector struct {
	Debug            bool
	Structs          []structDef
	PublicValidators map[string]bool
	Enums            map[string][]string
}

func newInspector(debug bool) *inspector {
	return &inspector{
		Debug:            debug,
		PublicValidators: make(map[string]bool),
		Enums:            make(map[string][]string),
	}
}

func (insp *inspector) inspect(dir, outputFile string) ([]structDef, error) {
	files, err := getFilesForInspect(dir, outputFile)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		fs := token.NewFileSet()
		parsedFile, err := parser.ParseFile(fs, f, nil, 0)
		if err != nil {
			log.Fatalf("Error parsing file: %s: %s", f, err)
		}
		ast.Walk(insp, parsedFile)
	}
	return insp.Structs, nil
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
		if f.IsDir() || f.Name() == outputFile || strings.HasSuffix(f.Name(), "_test.go") {
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
			return nil
		}
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

		for _, field := range astFields {

			fieldType := s.parseFieldType(field.Type)
			fieldName := s.parseFieldName(field.Names, fieldType)
			s.parseTagsAndAddField(fieldName, fieldType, field.Tag)

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

	case *ast.MapType:
		debug("skipping generation validator for alias on map: %s on %#v", structName, v)

	case *ast.InterfaceType:
		debug("skipping generation validator for alias on interface: %s on %#v", structName, v)

	case *ast.ArrayType:
		debug("skipping generation validator for alias on array: %s on %#v", structName, v)*/

	default:
		log.Fatalf("not expected Type for typeSpec: %s, %+v: %+v", structName, astTypeSpec, astTypeSpec.Type)
	}
}

func (s *structDef) parseFieldType(t ast.Expr) typeDef {
	switch v := t.(type) {
	case *ast.Ident:
		simple := getSimpleType(v.Name)
		if simple != nil {
			return simple
		}
		return typeStruct{NameStr: v.Name}

	case *ast.SelectorExpr:
		// v.X - contains pkg : if x, ok := v.X.(*ast.Ident); ok { x.Name+"."+v.Sel.Name)}
		return typeStruct{NameStr: v.Sel.Name}

	case *ast.ArrayType:
		return typeArray{InnerType: s.parseFieldType(v.Elt)}

	case *ast.StarExpr:
		return typePointer{InnerType: s.parseFieldType(v.X)}

	case *ast.InterfaceType:
		return typeInterface{}

	case *ast.MapType:
		return typeMap{Key: s.parseFieldType(v.Key), Value: s.parseFieldType(v.Value)}
	}
	panic(fmt.Errorf("Undefined typeField for %s: %+v", s.Name, t))
}

func (s *structDef) parseFieldName(fieldNames []*ast.Ident, fieldType typeDef) string {
	if len(fieldNames) != 0 {
		return fieldNames[0].Name
	}
	return fieldType.name() //wrapped struct, fieldName the same as type
}

func (s *structDef) parseTagsAndAddField(fieldName string, fieldType typeDef, astTag *ast.BasicLit) {

	tags := s.parseTags(astTag)

	fieldType.setTags(tags)

	s.Fields = append(s.Fields, fieldDef{Name: fieldName, Type: fieldType})
}

func (s *structDef) parseTags(astTag *ast.BasicLit) []tagDef { //example: `json:"place_type,omitempty" validate:"min=1,max=64"` OR `json:"user_id"`
	if astTag == nil {
		return nil
	}
	tagString := astTag.Value
	if tagString == "" {
		return nil
	}
	tagString = tagString[1 : len(tagString)-1] //clean from `json:"place_type,omitempty" validate:"min=1,max=64"` to  json:"place_type,omitempty" validate:"min=1,max=64"
	splittedTags := strings.Split(tagString, " ")

	for _, tagWithName := range splittedTags {
		if tagWithName == "" {
			continue
		}
		v := strings.SplitN(tagWithName, ":", 2)
		if len(v) != 2 {
			log.Fatalf("invalid tag for %s: %+v: %+v", s.Name, tagString, tagWithName)
		}
		tagName := strings.Trim(v[0], " ")
		if tagName == ValidateTag {
			var tags []tagDef
			functions := v[1][1 : len(v[1])-1] //clean quotes from "min=1,max=64" to min=1,max=64
			allFunctions := strings.Split(functions, ",")
			for _, functionWithParam := range allFunctions {
				tag := s.parseTagFunc(functionWithParam)
				tags = append(tags, tag)
			}
			return tags
		}
		//TODO: check for misspeling
	}
	return nil
}

func (s *structDef) parseTagFunc(functionWithParam string) tagDef {
	//TODO: any other tags for maps etc.
	tag := tagDef{}
	tv := strings.SplitN(functionWithParam, "=", 2)
	tag.Name = strings.Trim(tv[0], " ")

	if len(tv) > 1 {
		tag.Param = strings.Trim(tv[1], " ")
		if _, err := strconv.ParseFloat(tag.Param, 64); err != nil { // just validation, can be deleted
			log.Fatalf("not number tag param in %s: for %+v", s.Name, tag.Param)
		}
	}
	return tag
}

func getSimpleType(fieldType string) typeDef {
	switch fieldType {
	case "string":
		return typeMap{}

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
		return &typeNumber{}

	case "float32":
		fallthrough
	case "float64":
		return &typeNumber{}

	case "bool":
		return &typeBool{}
	}
	return nil
}
