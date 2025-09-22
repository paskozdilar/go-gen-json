package examples

import (
	"encoding/json/v2"
	"testing"
)

func benchUnmarshal[T any](in []byte) func(*testing.B) {
	return func(b *testing.B) {
		var v T
		b.ResetTimer()
		for b.Loop() {
			json.Unmarshal(in, &v)
		}
	}
}

// func benchMarshal[T any](v T) func (*testing.B){
//	return func(b *testing.B) {
//		b.ResetTimer()
//		for b.Loop() {
//			json.Marshal(&v)
//		}
//	}
// }

func BenchmarkEmptyStruct(b *testing.B) {
	type _EmptyStruct EmptyStruct
	b.Run("UnmarshalGen", benchUnmarshal[EmptyStruct](EmptyStructJSON))
	b.Run("UnmarshalNoGen", benchUnmarshal[_EmptyStruct](EmptyStructJSON))
	// b.Run("MarshalGen", benchMarshal[EmptyStruct](EmptyStructJSON))
	// b.Run("MarshalNoGen", benchMarshal[_EmptyStruct](EmptyStructJSON))
}

func BenchmarkBasicStruct(b *testing.B) {
	type _BasicStruct BasicStruct
	b.Run("UnmarshalGen", benchUnmarshal[BasicStruct](BasicStructJSON))
	b.Run("UnmarshalNoGen", benchUnmarshal[_BasicStruct](BasicStructJSON))
	// b.Run("MarshalGen", benchMarshal[BasicStruct](BasicStructJSON))
	// b.Run("MarshalNoGen", benchMarshal[_BasicStruct](BasicStructJSON))
}

func BenchmarkNestedStruct(b *testing.B) {
	type _NestedStruct NestedStruct
	b.Run("UnmarshalGen", benchUnmarshal[NestedStruct](NestedStructJSON))
	b.Run("UnmarshalNoGen", benchUnmarshal[_NestedStruct](NestedStructJSON))
	// b.Run("MarshalGen", benchMarshal[NestedStruct](NestedStructJSON))
	// b.Run("MarshalNoGen", benchMarshal[_NestedStruct](NestedStructJSON))
}

func BenchmarkComplexStruct(b *testing.B) {
	type _ComplexStruct ComplexStruct
	b.Run("UnmarshalGen", benchUnmarshal[ComplexStruct](ComplexStructJSON))
	b.Run("UnmarshalNoGen", benchUnmarshal[_ComplexStruct](ComplexStructJSON))
	// b.Run("MarshalGen", benchMarshal[ComplexStruct](ComplexStructJSON))
	// b.Run("MarshalNoGen", benchMarshal[_ComplexStruct](ComplexStructJSON))
}
