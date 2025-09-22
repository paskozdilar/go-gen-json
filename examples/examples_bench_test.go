package examples_test

import (
	"encoding/json/v2"
	"testing"

	"github.com/paskozdilar/go-gen-json/examples"
)

func benchUnmarshal[T, W any](in []byte) (func(*testing.B), func(*testing.B)) {
	fnt := func(b *testing.B) {
		b.ResetTimer()
		for b.Loop() {
			var v T
			json.Unmarshal(in, &v)
		}
	}
	fnw := func(b *testing.B) {
		b.ResetTimer()
		for b.Loop() {
			var w W
			json.Unmarshal(in, &w)
		}
	}
	return fnt, fnw
}

func benchMarshal[T, W any](v T, w W) (func(*testing.B), func(*testing.B)) {
	fnt := func(b *testing.B) {
		b.ResetTimer()
		for b.Loop() {
			json.Marshal(&v)
		}
	}
	fnw := func(b *testing.B) {
		b.ResetTimer()
		for b.Loop() {
			json.Marshal(&w)
		}
	}
	return fnt, fnw
}

func bench[T, W any](b *testing.B, in []byte, v T, w W) {
	b.Run("Unmarshal", func(b *testing.B) {
		if _, ok := any(&v).(json.Unmarshaler); !ok {
			b.Skipf("type %T does not implement json.Unmarshaler", &v)
		}
		fnt, fnw := benchUnmarshal[T, W](in)
		b.Run("Gen", fnt)
		b.Run("NoGen", fnw)
	})
	b.Run("Marshal", func(b *testing.B) {
		if _, ok := any(&v).(json.Marshaler); !ok {
			b.Skipf("type %T does not implement json.Marshaler", &v)
		}
		fnt, fnw := benchMarshal(v, w)
		b.Run("Gen", fnt)
		b.Run("NoGen", fnw)
	})
}

func BenchmarkNamedString(b *testing.B) {
	type _NamedString examples.NamedString
	bench(b, examples.NamedStringJSON, examples.NamedStringValue, _NamedString(examples.NamedStringValue))
}

func BenchmarkEmptyStruct(b *testing.B) {
	type _EmptyStruct examples.EmptyStruct
	bench(b, examples.EmptyStructJSON, examples.EmptyStructValue, _EmptyStruct(examples.EmptyStructValue))
}

func BenchmarkBasicStruct(b *testing.B) {
	type _BasicStruct examples.BasicStruct
	bench(b, examples.BasicStructJSON, examples.BasicStructValue, _BasicStruct(examples.BasicStructValue))
}

func BenchmarkNestedStruct(b *testing.B) {
	type _NestedStruct examples.NestedStruct
	bench(b, examples.NestedStructJSON, examples.NestedStructValue, _NestedStruct(examples.NestedStructValue))
}

func BenchmarkComplexStruct(b *testing.B) {
	type _ComplexStruct examples.ComplexStruct
	bench(b, examples.ComplexStructJSON, examples.ComplexStructValue, _ComplexStruct(examples.ComplexStructValue))
}

func BenchmarkEmbeddedStruct(b *testing.B) {
	type _EmbeddedStruct examples.EmbeddedStruct
	bench(b, examples.EmbeddedStructJSON, examples.EmbeddedStructValue, _EmbeddedStruct(examples.EmbeddedStructValue))
}
