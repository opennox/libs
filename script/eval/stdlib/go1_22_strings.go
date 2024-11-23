// Code generated by 'yaegi extract strings'. DO NOT EDIT.

//go:build go1.22
// +build go1.22

package stdlib

import (
	"reflect"
	"strings"
)

func init() {
	Symbols["strings/strings"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"Clone":          reflect.ValueOf(strings.Clone),
		"Compare":        reflect.ValueOf(strings.Compare),
		"Contains":       reflect.ValueOf(strings.Contains),
		"ContainsAny":    reflect.ValueOf(strings.ContainsAny),
		"ContainsFunc":   reflect.ValueOf(strings.ContainsFunc),
		"ContainsRune":   reflect.ValueOf(strings.ContainsRune),
		"Count":          reflect.ValueOf(strings.Count),
		"Cut":            reflect.ValueOf(strings.Cut),
		"CutPrefix":      reflect.ValueOf(strings.CutPrefix),
		"CutSuffix":      reflect.ValueOf(strings.CutSuffix),
		"EqualFold":      reflect.ValueOf(strings.EqualFold),
		"Fields":         reflect.ValueOf(strings.Fields),
		"FieldsFunc":     reflect.ValueOf(strings.FieldsFunc),
		"HasPrefix":      reflect.ValueOf(strings.HasPrefix),
		"HasSuffix":      reflect.ValueOf(strings.HasSuffix),
		"Index":          reflect.ValueOf(strings.Index),
		"IndexAny":       reflect.ValueOf(strings.IndexAny),
		"IndexByte":      reflect.ValueOf(strings.IndexByte),
		"IndexFunc":      reflect.ValueOf(strings.IndexFunc),
		"IndexRune":      reflect.ValueOf(strings.IndexRune),
		"Join":           reflect.ValueOf(strings.Join),
		"LastIndex":      reflect.ValueOf(strings.LastIndex),
		"LastIndexAny":   reflect.ValueOf(strings.LastIndexAny),
		"LastIndexByte":  reflect.ValueOf(strings.LastIndexByte),
		"LastIndexFunc":  reflect.ValueOf(strings.LastIndexFunc),
		"Map":            reflect.ValueOf(strings.Map),
		"NewReader":      reflect.ValueOf(strings.NewReader),
		"NewReplacer":    reflect.ValueOf(strings.NewReplacer),
		"Repeat":         reflect.ValueOf(strings.Repeat),
		"Replace":        reflect.ValueOf(strings.Replace),
		"ReplaceAll":     reflect.ValueOf(strings.ReplaceAll),
		"Split":          reflect.ValueOf(strings.Split),
		"SplitAfter":     reflect.ValueOf(strings.SplitAfter),
		"SplitAfterN":    reflect.ValueOf(strings.SplitAfterN),
		"SplitN":         reflect.ValueOf(strings.SplitN),
		"Title":          reflect.ValueOf(strings.Title),
		"ToLower":        reflect.ValueOf(strings.ToLower),
		"ToLowerSpecial": reflect.ValueOf(strings.ToLowerSpecial),
		"ToTitle":        reflect.ValueOf(strings.ToTitle),
		"ToTitleSpecial": reflect.ValueOf(strings.ToTitleSpecial),
		"ToUpper":        reflect.ValueOf(strings.ToUpper),
		"ToUpperSpecial": reflect.ValueOf(strings.ToUpperSpecial),
		"ToValidUTF8":    reflect.ValueOf(strings.ToValidUTF8),
		"Trim":           reflect.ValueOf(strings.Trim),
		"TrimFunc":       reflect.ValueOf(strings.TrimFunc),
		"TrimLeft":       reflect.ValueOf(strings.TrimLeft),
		"TrimLeftFunc":   reflect.ValueOf(strings.TrimLeftFunc),
		"TrimPrefix":     reflect.ValueOf(strings.TrimPrefix),
		"TrimRight":      reflect.ValueOf(strings.TrimRight),
		"TrimRightFunc":  reflect.ValueOf(strings.TrimRightFunc),
		"TrimSpace":      reflect.ValueOf(strings.TrimSpace),
		"TrimSuffix":     reflect.ValueOf(strings.TrimSuffix),

		// type definitions
		"Builder":  reflect.ValueOf((*strings.Builder)(nil)),
		"Reader":   reflect.ValueOf((*strings.Reader)(nil)),
		"Replacer": reflect.ValueOf((*strings.Replacer)(nil)),
	}
}