package examples_test

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

// svaluef returns a string representation of v suitable for debugging and test
// failure messages.
func svaluef(v any) string {
	var b strings.Builder
	var visit func(reflect.Value)

	visit = func(val reflect.Value) {
		if !val.IsValid() {
			b.WriteString("<invalid>")
			return
		}

		// Dereference pointers
		for val.Kind() == reflect.Pointer {
			if val.IsNil() {
				b.WriteString("<nil>")
				return
			}
			val = val.Elem()
		}

		switch val.Kind() {
		case reflect.Struct:
			b.WriteString(val.Type().Name())
			b.WriteString("{")

			type Field struct {
				idx  int
				name string
			}
			fields := []Field{}
			for i := 0; i < val.NumField(); i++ {
				field := val.Type().Field(i)
				if field.IsExported() {
					fields = append(fields, Field{i, field.Name})
				}
			}
			slices.SortFunc(fields, func(a, b Field) int {
				return strings.Compare(a.name, b.name)
			})

			for _, f := range fields {
				i := f.idx
				field := val.Type().Field(i)
				if i > 0 {
					b.WriteString(" ")
				}
				b.WriteString(field.Name)
				b.WriteString(":")
				visit(val.Field(i))
			}
			b.WriteString("}")
		case reflect.Slice, reflect.Array:
			b.WriteString("[")
			for i := 0; i < val.Len(); i++ {
				if i > 0 {
					b.WriteString(" ")
				}
				visit(val.Index(i))
			}
			b.WriteString("]")
		case reflect.Map:
			b.WriteString("{")
			keys := val.MapKeys()
			slices.SortFunc(keys, func(a, b reflect.Value) int {
				return strings.Compare(fmt.Sprintf("%v", a.Interface()), fmt.Sprintf("%v", b.Interface()))
			})
			for i, key := range keys {
				if i > 0 {
					b.WriteString(" ")
				}
				visit(key)
				b.WriteString(":")
				visit(val.MapIndex(key))
			}
			b.WriteString("}")
		case reflect.Interface:
			if val.IsNil() {
				b.WriteString("<nil>")
				return
			}
			visit(val.Elem())
		default:
			b.WriteString(fmt.Sprintf("%v", val.Interface()))
		}
	}

	visit(reflect.ValueOf(v))
	return b.String()
}

// sdiff returns a list of human-readable differences between a and b.
func sdiff(a, b any) []string {
	var diffs []string
	var walk func(path string, va, vb reflect.Value)

	walk = func(path string, va, vb reflect.Value) {
		if !va.IsValid() || !vb.IsValid() {
			if va.IsValid() != vb.IsValid() {
				diffs = append(diffs, path+" (validity mismatch)")
			}
			return
		}

		if va.Type() != vb.Type() {
			diffs = append(diffs,
				fmt.Sprintf("%s (type mismatch: %T vs %T)", path, va.Interface(), vb.Interface()))
			return
		}

		// Dereference pointers
		for va.Kind() == reflect.Pointer && !va.IsNil() {
			va = va.Elem()
		}
		for vb.Kind() == reflect.Pointer && !vb.IsNil() {
			vb = vb.Elem()
		}

		if !va.IsValid() || !vb.IsValid() {
			if va.IsValid() != vb.IsValid() {
				diffs = append(diffs, path+" (validity mismatch)")
			}
			return
		}

		switch va.Kind() {
		case reflect.Struct:
			for i := 0; i < va.NumField(); i++ {
				field := va.Type().Field(i)
				if !field.IsExported() {
					continue
				}
				walk(path+"."+field.Name, va.Field(i), vb.Field(i))
			}
		case reflect.Slice, reflect.Array:
			if va.Len() != vb.Len() {
				diffs = append(diffs, fmt.Sprintf("%s (len mismatch %d vs %d)", path, va.Len(), vb.Len()))
				return
			}
			for i := 0; i < va.Len(); i++ {
				walk(fmt.Sprintf("%s[%d]", path, i), va.Index(i), vb.Index(i))
			}
		case reflect.Map:
			if va.Len() != vb.Len() {
				diffs = append(diffs, fmt.Sprintf("%s (map len mismatch)", path))
				return
			}
			for _, key := range va.MapKeys() {
				walk(fmt.Sprintf("%s[%v]", path, key), va.MapIndex(key), vb.MapIndex(key))
			}
		case reflect.Interface:
			if va.IsNil() || vb.IsNil() {
				if va.IsNil() != vb.IsNil() {
					diffs = append(diffs, path+" (interface nil mismatch)")
				}
				return
			}
			walk(path, va.Elem(), vb.Elem())
		default:
			areEqual := reflect.DeepEqual(va.Interface(), vb.Interface())
			if !areEqual {
				diffs = append(diffs, fmt.Sprintf("%s (value mismatch: %v vs %v)", path, va.Interface(), vb.Interface()))
				return
			}
		}
	}

	walk("", reflect.ValueOf(a), reflect.ValueOf(b))
	return diffs
}
