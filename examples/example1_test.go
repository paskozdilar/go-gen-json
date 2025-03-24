package examples_test

import (
	"encoding/json"
	"testing"

	"github.com/paskozdilar/go-gen-json/examples"
)

func TestExample1(t *testing.T) {
	e1 := &examples.Example1{}
	e2 := &examples.Example1{}
	in := []byte(`{"foo":"hello","bar":42}`)

	if err := json.Unmarshal(in, e1); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	if err := e2.ParseJSON(in); err != nil {
		t.Fatalf("e2.ParseJSON failed: %v", err)
	}

	if e1.Foo != e2.Foo {
		t.Fatalf("Foo mismatch: %v != %v", e1.Foo, e2.Foo)
	}
	if e1.Bar != e2.Bar {
		t.Fatalf("Bar mismatch: %v != %v", e1.Bar, e2.Bar)
	}
}

func BenchmarkExample1(b *testing.B) {
	e := &examples.Example1{}
	in := []byte(`{"foo":"hello","bar":42}`)

	b.Run("json.Unmarshal", func(b *testing.B) {
		for b.Loop() {
			if err := json.Unmarshal(in, e); err != nil {
				b.Fatalf("json.Unmarshal failed: %v", err)
			}
		}
	})
	b.Run("fb.ParseJSON", func(b *testing.B) {
		for b.Loop() {
			if err := e.ParseJSON(in); err != nil {
				b.Fatalf("fb.ParseJSON failed: %v", err)
			}
		}
	})
}
