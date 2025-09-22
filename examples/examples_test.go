package examples_test

import (
	"bytes"
	"encoding/json/v2"
	"log"
	"reflect"
	"testing"

	"github.com/paskozdilar/go-gen-json/examples"
)

func testUnmarshal[T any](in []byte, out T) func(*testing.T) {
	return func(t *testing.T) {
		var v T
		if _, ok := any(&v).(json.Unmarshaler); !ok {
			t.Skipf("type %T does not implement json.Unmarshaler", &v)
		}
		if err := json.Unmarshal(in, &v); err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}
		if !reflect.DeepEqual(v, out) {
			log.Fatalf(
				"unmarshal error: differes from json/v2, got: %v, want: %v, diff: %v",
				svaluef(v), svaluef(out), sdiff(v, out),
			)
		}
	}
}

func testMarshal[T any](v T, out []byte) func(*testing.T) {
	return func(t *testing.T) {
		if _, ok := any(&v).(json.Marshaler); !ok {
			t.Skipf("type %T does not implement json.Marshaler", &v)
		}
		b, err := json.Marshal(&v)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}
		if !bytes.Equal(b, out) {
			log.Fatalf(
				"marshal error: differes from json/v2, got: %s, want: %s",
				b, out,
			)
		}
	}
}

func TestNamedString(t *testing.T) {
	t.Run("Unmarshal", testUnmarshal(examples.NamedStringJSON, examples.NamedStringValue))
	t.Run("Marshal", testMarshal(examples.NamedStringValue, examples.NamedStringJSON))
}

func TestEmptyStruct(t *testing.T) {
	t.Run("Unmarshal", testUnmarshal(examples.EmptyStructJSON, examples.EmptyStructValue))
	t.Run("Marshal", testMarshal(examples.EmptyStructValue, examples.EmptyStructJSON))
}

func TestBasicStruct(t *testing.T) {
	t.Run("Unmarshal", testUnmarshal(examples.BasicStructJSON, examples.BasicStructValue))
	t.Run("Marshal", testMarshal(examples.BasicStructValue, examples.BasicStructJSON))
}

func TestNestedStruct(t *testing.T) {
	t.Run("Unmarshal", testUnmarshal(examples.NestedStructJSON, examples.NestedStructValue))
	t.Run("Marshal", testMarshal(examples.NestedStructValue, examples.NestedStructJSON))
}

func TestComplexStruct(t *testing.T) {
	t.Run("Unmarshal", testUnmarshal(examples.ComplexStructJSON, examples.ComplexStructValue))
	t.Run("Marshal", testMarshal(examples.ComplexStructValue, examples.ComplexStructJSON))
}

func TestEmbeddedStruct(t *testing.T) {
	t.Run("Unmarshal", testUnmarshal(examples.EmbeddedStructJSON, examples.EmbeddedStructValue))
	t.Run("Marshal", testMarshal(examples.EmbeddedStructValue, examples.EmbeddedStructJSON))
}
