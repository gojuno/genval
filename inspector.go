package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"strings"

	"github.com/gojuno/genval/types"
)

type inspector struct {
	structs          []StructDef
	publicValidators map[string]bool
	enums            map[string][]string
}

func NewInspector() *inspector {
	return &inspector{
		publicValidators: make(map[string]bool),
		enums:            make(map[string][]string),
	}
}

func (insp *inspector) Inspect(dir, outputFile string) error {
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
func (insp *inspector) Result() []StructDef {
	res := insp.structs
	for i, s := range res {
		if v, ok := insp.publicValidators[s.Name]; ok {
			res[i].PublicValidatorExist = v
		}
		if v, ok := insp.enums[s.Name]; ok {
			res[i].EnumValues = v
		}
	}
	return res
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
		if methodName == "Validate" && spec.Recv != nil {
			st := spec.Recv.List[0].Type
			if x, ok := st.(*ast.Ident); ok {
				insp.publicValidators[x.Name] = true
			} else {
				log.Fatalf("method Validate should be: 'func (s Struct)Validate() error{...}' not on pointer, %+v,%T", st, st)
			}
		}
		return nil
	case *ast.ValueSpec: //To find all consts and then generate validation with consts
		if spec.Names[0].Obj.Kind == ast.Con {
			valueName := spec.Names[0].Name
			if x, ok := spec.Type.(*ast.Ident); ok {
				valueType := x.Name
				if !isSimple(valueType) {
					insp.enums[valueType] = append(insp.enums[valueType], valueName)
				}
			}
			return nil
		}
	}
	return insp
}

func (insp *inspector) addStruct(s StructDef) {
	insp.structs = append(insp.structs, s)
}

func (insp *inspector) visitStruct(astTypeSpec *ast.TypeSpec) {
	structName := astTypeSpec.Name.Name
	switch v := astTypeSpec.Type.(type) {
	case *ast.StructType:
		astFields := v.Fields.List
		s := NewFieldsStruct(structName)
		for _, astField := range astFields {
			fieldType := parseFieldType(astField.Type, fmt.Sprintf("struct %s", structName))
			fieldName := parseFieldName(astField.Names, fieldType)
			tags := types.ParseTags(astField.Tag, fmt.Sprintf("struct %s, field %s", structName, fieldName))

			field, err := NewField(fieldName, fieldType, tags)
			if err != nil {
				log.Fatalf("field creation failed for struct %s, %s", structName, err)
			}
			s.AddField(*field)
		}
		insp.addStruct(s)
		return
	case *ast.Ident, *ast.SelectorExpr, *ast.StarExpr, *ast.MapType, *ast.InterfaceType, *ast.ArrayType: //aliases
		aliasType := parseFieldType(v, fmt.Sprintf("struct %s", structName))
		insp.addStruct(NewAliasStruct(structName, aliasType))
	default:
		log.Fatalf("not expected Type for typeSpec: %s, %+v: %T", structName, astTypeSpec, astTypeSpec.Type)
	}
}

func parseFieldType(t ast.Expr, logCtx string) types.TypeDef {
	switch v := t.(type) {
	case *ast.Ident:
		simple := getSimpleType(v.Name)
		if simple != nil {
			return simple
		}
		return types.NewStruct(v.Name)
	case *ast.SelectorExpr:
		return types.NewExternalStruct(v.Sel.Name)
	case *ast.ArrayType:
		return types.NewArray(parseFieldType(v.Elt, logCtx))
	case *ast.StarExpr:
		return types.NewPointer(parseFieldType(v.X, logCtx))
	case *ast.InterfaceType:
		return types.NewInterface()
	case *ast.MapType:
		return types.NewMap(parseFieldType(v.Key, logCtx), parseFieldType(v.Value, logCtx))
	case *ast.FuncType:
		return types.NewFunc()
	}
	log.Fatalf("undefined typeField for %s: %+v, %T", logCtx, t, t)
	return nil
}

func parseFieldName(fieldNames []*ast.Ident, fieldType types.TypeDef) string {
	if len(fieldNames) != 0 {
		return fieldNames[0].Name
	}
	return fieldType.Type() //wrapped struct, fieldName the same as type
}

func isSimple(fieldType string) bool {
	return getSimpleType(fieldType) != nil
}

func getSimpleType(fieldType string) types.TypeDef {
	switch fieldType {
	case "string":
		return types.NewString()
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return types.NewNumber(fieldType)
	case "float32", "float64":
		return types.NewNumber(fieldType)
	case "bool":
		return types.NewBool()
	}
	return nil
}
