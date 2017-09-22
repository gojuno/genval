package types

import (
	"go/ast"
	"log"
	"strings"
)

// ValidatableTag represents a tag that uses for validating
type ValidatableTag interface {
	// Key returns a tag title
	Key() string
}

// ValidatableTags represents a set of ValidatableTag
type ValidatableTags []ValidatableTag

// Empty check whether ValidatableTags is empty
func (ts ValidatableTags) Empty() bool {
	return len(ts) == 0
}

// ContainsTag checks whether ValidatableTags contains a ValidatableTag
func (ts ValidatableTags) ContainsTag(t ValidatableTag) bool {
	for _, v := range ts {
		if v == t {
			return true
		}
	}
	return false
}

// contains {"json": "provider_id", "xml": "provider_id"} for
// `json:"provider_id,omitempty" "xml":"provider_id" validate:"min_len=1"
type FieldTagsNames map[string]string

const FieldNameFromStructDefinition = ""

func (n FieldTagsNames) Get(name string) string {
	res, ok := n[strings.ToLower(name)]
	if !ok {
		return n.GetFromStructDefinition()
	}
	return res
}

func (n FieldTagsNames) GetFromStructDefinition() string {
	return n[FieldNameFromStructDefinition]
}

type SimpleTag struct {
	Name  string
	Param string
}

func (t SimpleTag) Key() string {
	return t.Name
}

type ScopeTag struct {
	Name      string
	InnerTags []ValidatableTag
}

func (t ScopeTag) Key() string {
	return t.Name
}

func ParseTags(astTag *ast.BasicLit, logCtx string) ([]ValidatableTag, FieldTagsNames) { //example: `json:"place_type,omitempty" validate:"min=1,max=64"` OR `json:"user_id"`
	if astTag == nil {
		return nil, nil
	}
	tagString := astTag.Value
	if tagString == "" {
		return nil, nil
	}
	tagString = removeQuotes(tagString) //clean from `json:"place_type,omitempty" validate:"min=1,max=64"` to  json:"place_type,omitempty" validate:"min=1,max=64"
	splittedTags := strings.Split(tagString, " ")
	validateTags := []ValidatableTag{}
	nameTags := FieldTagsNames{}
	for _, tagWithName := range splittedTags {
		if tagWithName == "" {
			continue
		}
		v := strings.SplitN(tagWithName, ":", 2)
		if len(v) != 2 {
			log.Fatalf("invalid tag for %s: %s", logCtx, tagWithName)
		}
		tagName := strings.Trim(v[0], " ")
		if tagName == ValidateTag {
			validateTags = parseFuncs(removeQuotes(v[1]), logCtx) //clean quotes from "min=1,max=64" to min=1,max=64
			continue
		}
		for _, m := range misspellValidate {
			if m == tagName {
				log.Fatalf("tag validate is misspelled for %s: %s", logCtx, tagName)
			}
		}
		tagValues := strings.Split(strings.TrimSpace(v[1]), ",")
		if len(tagValues) == 1 { // e.g. json:"field_name"
			nameTags[tagName] = removeQuotes(tagValues[0])
		} else { // e.g. handle json:"field_name,omitempty" where
			// tagValues[0] = '"field_name', so we need to cut first char "
			nameTags[tagName] = tagValues[0][1:len(tagValues[0])]
		}
	}
	return validateTags, nameTags
}

func scopeWasParsedRight(tagFunc string) bool {
	return strings.Count(tagFunc, "[") == strings.Count(tagFunc, "]")
}

func parseFuncs(functions string, logCtx string) []ValidatableTag { //min_items=5,key=[min=4,max=59],value=[min_len=1,max_len=64]
	var funcs []string
	cur := ""
	for _, f := range strings.Split(functions, ",") {
		cur += f
		if scopeWasParsedRight(cur) {
			funcs = append(funcs, cur)
			cur = ""
		} else {
			cur += ","
		}
	}
	if cur != "" {
		log.Fatalf("parse funcs failed for %s, tag: %s", logCtx, functions)
	}
	var res []ValidatableTag
	for _, funcWithParam := range funcs {
		tv := strings.SplitN(funcWithParam, "=", 2)
		name := strings.Trim(tv[0], " ")
		param := ""
		if len(tv) > 1 {
			param = strings.Trim(tv[1], " ")
		}
		if strings.HasPrefix(param, "[") {
			res = append(res, ScopeTag{
				Name:      name,
				InnerTags: parseFuncs(removeQuotes(param), logCtx), //remove '[' and ']'
			})
		} else {
			res = append(res, SimpleTag{
				Name:  name,
				Param: param,
			})
		}
	}
	return res
}

func removeQuotes(s string) string {
	if len(s) < 2 {
		log.Fatalf("bad input for removing quotes: %s", s)
	}
	return s[1 : len(s)-1]
}
