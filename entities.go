package main

// type structDef interface {
// 	Name() string
// }

type structDef struct {
	Name   string
	Fields []fieldDef
}
type fieldDef struct {
	Name string
	Type typeDef
	Tags []tagDef
}

type typeDef interface {
	Name() string
}

type typeStruct struct {
	NameStr string
}

func (t typeStruct) Name() string {
	return t.NameStr
}

type typeArray struct {
	InnerType typeDef
}

func (t typeArray) Name() string {
	return "[]" + t.InnerType.Name()
}

type typePointer struct {
	InnerType typeDef
}

func (t typePointer) Name() string {
	return "*" + t.InnerType.Name()
}

type typeMap struct {
	Key   typeDef
	Value typeDef
}

func (t typeMap) Name() string {
	return "map[" + t.Key.Name() + "]" + t.Value.Name()
}

type typeInterface struct {
}

func (t typeInterface) Name() string {
	return "interface{}"
}

type tagDef struct {
	Name  string
	Param string
}

/*
type aliasStruct struct {
	Name      string
	AliasType string
}

func (a aliasStruct) Name() string {
	return a.Name
}
*/
