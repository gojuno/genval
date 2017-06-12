package types

import (
	"go/ast"
	"log"
	"strings"
)

type Tag interface {
	Key() string
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
	InnerTags []Tag
}

func (t ScopeTag) Key() string {
	return t.Name
}

func ParseTags(astTag *ast.BasicLit, logCtx string) []Tag { //example: `json:"place_type,omitempty" validate:"min=1,max=64"` OR `json:"user_id"`
	if astTag == nil {
		return nil
	}
	tagString := astTag.Value
	if tagString == "" {
		return nil
	}
	tagString = removeQuotes(tagString) //clean from `json:"place_type,omitempty" validate:"min=1,max=64"` to  json:"place_type,omitempty" validate:"min=1,max=64"
	splittedTags := strings.Split(tagString, " ")

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
			return parseFuncs(removeQuotes(v[1]), logCtx) //clean quotes from "min=1,max=64" to min=1,max=64
		}
		for _, m := range misspellValidate {
			if m == tagName {
				log.Fatalf("tag validate is misspelled for %s: %s", logCtx, tagName)
			}
		}
	}
	return nil
}

func scopeWasParsedRight(tagFunc string) bool {
	return strings.Count(tagFunc, "[") == strings.Count(tagFunc, "]")
}

func parseFuncs(functions string, logCtx string) []Tag { //min_items=5,key=[min=4,max=59],value=[min_len=1,max_len=64]
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
	var res []Tag
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
