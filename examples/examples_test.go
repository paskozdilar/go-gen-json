package examples

import (
	"encoding/json/v2"
	"log"
	"reflect"
	"testing"
)

// TODO: modify into `testMarshal` and `testUnmarshal` akin to bench
func run[T, W any](t *testing.T, name string, v1 T, v2 W) {
	b1, err := json.Marshal(&v1)
	if err != nil {
		t.Fatalf("marshal %s error: %v", name, err)
	}
	b2, err := json.Marshal(&v2)
	if err != nil {
		t.Fatalf("marshal %s error: %v", name, err)
	}
	if err := json.Unmarshal(b1, &v1); err != nil {
		t.Fatalf("unmarshal %s error: %v", name, err)
	}
	if err := json.Unmarshal(b2, &v2); err != nil {
		t.Fatalf("unmarshal %s error: %v", name, err)
	}
	var c1, c2 W
	json.Unmarshal(b1, &c1)
	json.Unmarshal(b2, &c2)
	if !reflect.DeepEqual(c1, c2) {
		log.Fatalf(
			"marshal %s error: differes from json/v2, got: %s, want: %s",
			name, b1, b2,
		)
	}
}

func TestEmptyStruct(t *testing.T) {
	type W EmptyStruct
	run(t, "EmptyStruct", EmptyStructValue, W(EmptyStructValue))
}

func TestBasicStruct(t *testing.T) {
	type W BasicStruct
	run(t, "BasicStruct", BasicStructValue, W(BasicStructValue))
}

func TestNestedStruct(t *testing.T) {
	type W NestedStruct
	run(t, "NestedStruct", NestedStructValue, W(NestedStructValue))
}

func TestComplexStruct(t *testing.T) {
	type W ComplexStruct
	run(t, "ComplexStruct", ComplexStructValue, W(ComplexStructValue))
}
