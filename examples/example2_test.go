package examples

import (
	"encoding/json"
	"testing"
)

func TestExample2(t *testing.T) {
	e1 := &Example2{}
	e2 := &Example2{}
	in := []byte(`{"foo":"hello","bar":{"baz":[1,2,3],"qux":{"quux":3.14}}}`)

	if err := json.Unmarshal(in, e1); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	if err := e2.ParseJSON(in); err != nil {
		t.Fatalf("e2.ParseJSON failed: %v", err)
	}

	if e1.Foo != e2.Foo {
		t.Fatalf("Foo mismatch: %v != %v", e1.Foo, e2.Foo)
	}
	if len(e1.Bar.Baz) != len(e2.Bar.Baz) {
		t.Fatalf("Bar.Baz length mismatch: %v != %v", len(e1.Bar.Baz), len(e2.Bar.Baz))
	}
	for i := range e1.Bar.Baz {
		if e1.Bar.Baz[i] != e2.Bar.Baz[i] {
			t.Fatalf("Bar.Baz mismatch: %v != %v", e1.Bar.Baz[i], e2.Bar.Baz[i])
		}
	}
	if e1.Bar.Qux == nil && e2.Bar.Qux != nil {
		t.Fatalf("Bar.Qux mismatch: %v != %v", e1.Bar.Qux, e2.Bar.Qux)
	}
	if e1.Bar.Qux != nil && e2.Bar.Qux == nil {
		t.Fatalf("Bar.Qux mismatch: %v != %v", e1.Bar.Qux, e2.Bar.Qux)
	}
	if e1.Bar.Qux != nil && e2.Bar.Qux != nil {
		if *e1.Bar.Qux.Quux != *e2.Bar.Qux.Quux {
			t.Fatalf("Bar.Qux.Quux mismatch: %v != %v", *e1.Bar.Qux.Quux, *e2.Bar.Qux.Quux)
		}
	}
}

func BenchmarkExample2(b *testing.B) {
	e := &Example2{}
	in := []byte(`{"foo":"hello","bar":{"baz":[1,2,3],"qux":{"quux":3.14}}}`)

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
