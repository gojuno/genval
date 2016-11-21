package main

import "fmt"

type structDef struct {
	Name   string
	Fields []fieldDef
}
type fieldDef struct {
	Name string
	Type typeDef
}

type tagDef struct {
	Name  string
	Param string
	Used  bool
}

type tagDefs []tagDef

//TYPES
type typeDef interface {
	name() string
	setTags(tagDefs)
}

//String
type typeString struct {
	MinLen *tagDef
	MaxLen *tagDef
}

func (t *typeString) name() string {
	return "string"
}

func (t *typeString) setTags(tags tagDefs) {
	t.MinLen = tags.useTag(StringMinLen)
	t.MaxLen = tags.useTag(StringMaxLen)

	if !tags.allUsed() {
		panic(fmt.Errorf("not all tags used: %+v", tags))
	}
}

//Number
type typeNumber struct {
	Min *tagDef
	Max *tagDef
}

func (t *typeNumber) name() string {
	return "number"
}

func (t *typeNumber) setTags(tags tagDefs) {
	t.Min = tags.useTag(NumberMin)
	t.Max = tags.useTag(NumberMax)

	if !tags.allUsed() {
		panic(fmt.Errorf("not all tags used: %+v", tags))
	}
}

//Bool
type typeBool struct {
}

func (t *typeBool) name() string {
	return "bool"
}

func (t *typeBool) setTags(tags tagDefs) {
	if !tags.allUsed() {
		panic(fmt.Errorf("not all tags used: %+v", tags))
	}
}

//Struct
type typeStruct struct {
	NameStr string
	Tags    []tagDef
}

func (t typeStruct) name() string {
	return t.NameStr
}

func (t typeStruct) setTags(tags tagDefs) {
	t.Tags = tags.notUsed()
}

//Array
type typeArray struct {
	InnerType typeDef
	MinItems  *tagDef
	MaxItems  *tagDef
}

func (t typeArray) name() string {
	return "[]" + t.InnerType.name()
}
func (t typeArray) setTags(tags tagDefs) {
	t.MinItems = tags.useTag(ArrayMinItems)
	t.MaxItems = tags.useTag(ArrayMaxItems)

	t.InnerType.setTags(tags)
}

//Pointer
type typePointer struct {
	InnerType typeDef
	Nullable  *tagDef
}

func (t typePointer) name() string {
	return "*" + t.InnerType.name()
}
func (t typePointer) setTags(tags tagDefs) {
	t.Nullable = tags.useTag(PointerNullable)

	t.InnerType.setTags(tags)
}

//Map
type typeMap struct {
	Key   typeDef
	Value typeDef
}

func (t typeMap) name() string {
	return "map[" + t.Key.name() + "]" + t.Value.name()
}

func (t typeMap) setTags(tags tagDefs) {
	//TODO: implement me
}

//Interface
type typeInterface struct {
}

func (t typeInterface) name() string {
	return "interface{}"
}
func (t typeInterface) setTags(tags tagDefs) {
	if !tags.allUsed() {
		panic(fmt.Errorf("not all tags used: %+v", tags))
	}
}

//END of TYPES

func (tags tagDefs) useTag(name string) *tagDef {
	for i, tag := range tags {
		if tag.Used {
			continue
		}
		if tag.Name == name {
			tags[i].Used = true
			return &tag
		}
	}
	return nil
}

func (tags tagDefs) allUsed() bool {
	for _, tag := range tags {
		if !tag.Used {
			return false
		}
	}
	return true
}

func (tags tagDefs) notUsed() []tagDef {
	var res []tagDef
	for _, tag := range tags {
		if !tag.Used {
			res = append(res, tag)
		}
	}
	return res
}
