package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"strings"

	"github.com/l1va/genval/types"
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
			insp.enums[s.Name] = nil //mark as used
		}
	}
	for k, v := range insp.enums {
		s := StructDef{
			Name:       k,
			EnumValues: v,
		}
		if p, ok := insp.publicValidators[s.Name]; ok {
			s.PublicValidatorExist = p
		}
		res = append(res, s)
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
		if methodName == "Validate" {
			st := spec.Recv.List[0].Type
			if star, ok := st.(*ast.StarExpr); ok {
				if x, ok := star.X.(*ast.Ident); ok {
					insp.publicValidators[x.Name] = true
				}
			}
		}
		return nil
	case *ast.ValueSpec: //To find all consts and then generate validation with consts
		if spec.Names[0].Obj.Kind == ast.Con {
			valueName := spec.Names[0].Name
			if x, ok := spec.Type.(*ast.Ident); ok {
				valueType := x.Name
				isSimple := nil != getSimpleType(valueType)
				if !isSimple {
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
		s := NewStruct(structName)
		for _, astField := range astFields {
			fieldType := parseFieldType(astField.Type, fmt.Sprintf("struct %s", structName))
			fieldName := parseFieldName(astField.Names, fieldType)
			tags := ParseTags(astField.Tag, fmt.Sprintf("struct %s, field %s", structName, fieldName))

			field, err := NewField(fieldName, fieldType, tags)
			if err != nil {
				log.Fatalf("field creatinon failed for struct %s, %s", structName, err)
			}
			s.AddField(*field)
		}
		insp.addStruct(s)
		return
	case *ast.Ident: //alias on simple type
	case *ast.SelectorExpr: //alias on struct
	case *ast.StarExpr: //alias on pointer
	case *ast.MapType: //alias on map
	case *ast.InterfaceType: //alias on interface
	case *ast.ArrayType: //alias on array
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
		return types.NewStructType(v.Name)
	case *ast.SelectorExpr:
		// v.X - contains pkg : if x, ok := v.X.(*ast.Ident); ok { x.Name+"."+v.Sel.Name)}
		return types.NewStructType(v.Sel.Name)
	case *ast.ArrayType:
		return types.NewArrayType(parseFieldType(v.Elt, logCtx))
	case *ast.StarExpr:
		return types.NewPointerType(parseFieldType(v.X, logCtx))
	case *ast.InterfaceType:
		return types.NewInterfaceType()
	case *ast.MapType:
		return types.NewMapType(parseFieldType(v.Key, logCtx), parseFieldType(v.Value, logCtx))
	case *ast.FuncType:
		return types.NewFuncType()
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

func getSimpleType(fieldType string) types.TypeDef {
	switch fieldType {
	case "string":
		return types.NewStringType()
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return types.NewNumberType(fieldType)
	case "float32", "float64":
		return types.NewNumberType(fieldType)
	case "bool":
		return types.NewBoolType()
	}
	return nil
}
