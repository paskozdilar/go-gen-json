// Example of a more complex structure, with nested objects, pointers and
// arrays.

package examples

import (
	"bytes"
	"fmt"
	"math"

	"github.com/go-json-experiment/json/jsontext"
	"github.com/paskozdilar/go-gen-json/jsonutil"
)

type Example2 struct {
	Foo string `json:"foo"`
	Bar struct {
		Baz []int `json:"baz"`
		Qux *struct {
			Quux *float64 `json:"quux"`
		}
	}
}

type (
	field_Example2_Foo = string
	field_Example2_Bar = struct {
		Baz field_Example2_Bar_Baz `json:"baz"`
		Qux *field_Example2_Bar_Qux
	}
	field_Example2_Bar_Baz = []int
	field_Example2_Bar_Qux = struct {
		Quux *field_Example2_Bar_Qux_Quux `json:"quux"`
	}
	field_Example2_Bar_Qux_Quux = float64
)

func (target *Example2) ParseJSON(b []byte) error {
	d := jsontext.NewDecoder(bytes.NewReader(b))
	return parseJSON_Example2(target, d)
}

func parseJSON_Example2(target *Example2, d *jsontext.Decoder) error {
	pos := d.InputOffset()
	t, err := d.ReadToken()
	if err != nil {
		return fmt.Errorf("error at offset %d: %w", pos, err)
	}
	if t.Kind() != '{' {
		return fmt.Errorf("error at offset %d: expected '{', got %v", pos, t.Kind())
	}
	for {
		pos = d.InputOffset()
		t, err := d.ReadToken()
		if err != nil {
			return fmt.Errorf("error at offset %d: %w", pos, err)
		}
		switch t.Kind() {
		default:
			return fmt.Errorf("error at offset %d: expected '\"' or '}', got %v", pos, t.Kind())
		case '}':
			return nil
		case '"':
		}
		switch t.String() {
		case "foo":
			err = parseJSON_Example2_Foo(&target.Foo, d)
		case "bar":
			err = parseJSON_Example2_Bar(&target.Bar, d)
		default:
			err = jsonutil.Discard(d)
		}
		if err != nil {
			return err
		}
	}
}

func parseJSON_Example2_Foo(target *field_Example2_Foo, d *jsontext.Decoder) error {
	pos := d.InputOffset()
	t, err := d.ReadToken()
	if err != nil {
		return fmt.Errorf("error at offset %d: %w", pos, err)
	}
	if t.Kind() != '"' {
		return fmt.Errorf("error at offset %d: expected string, got %v", pos, t.Kind())
	}
	*target = t.String()
	return nil
}

func parseJSON_Example2_Bar(target *field_Example2_Bar, d *jsontext.Decoder) error {
	pos := d.InputOffset()
	t, err := d.ReadToken()
	if err != nil {
		return fmt.Errorf("error at offset %d: %w", pos, err)
	}
	if t.Kind() != '{' {
		return fmt.Errorf("error at offset %d: expected '{', got %v", pos, t.Kind())
	}
	for {
		pos = d.InputOffset()
		t, err := d.ReadToken()
		if err != nil {
			return fmt.Errorf("error at offset %d: %w", pos, err)
		}
		switch t.Kind() {
		default:
			return fmt.Errorf("error at offset %d: expected '\"' or '}', got %v", pos, t.Kind())
		case '}':
			return nil
		case '"':
		}
		switch t.String() {
		case "baz":
			err = parseJSON_Example2_Bar_Baz(&target.Baz, d)
		case "qux":
			err = parseJSON_Example2_Bar_Qux(&target.Qux, d)
		default:
			err = jsonutil.Discard(d)
		}
		if err != nil {
			return err
		}
	}
}

func parseJSON_Example2_Bar_Baz(target *field_Example2_Bar_Baz, d *jsontext.Decoder) error {
	pos := d.InputOffset()
	t, err := d.ReadToken()
	if err != nil {
		return fmt.Errorf("error at offset %d: %w", pos, err)
	}
	if t.Kind() != '[' {
		return fmt.Errorf("error at offset %d: expected '[', got %v",
			pos, t.Kind())
	}
	for {
		pos := d.InputOffset()
		t, err := d.ReadToken()
		if err != nil {
			return fmt.Errorf("error at offset %d: %w", pos, err)
		}
		if t.Kind() == ']' {
			break
		}
		if t.Kind() != '0' {
			return fmt.Errorf("error at offset %d: expected number, got %v", pos, t.Kind())
		}
		i := t.Int()
		if i < math.MinInt || i > math.MaxInt {
			return fmt.Errorf("error at offset %d: number out of range: %d", pos, i)
		}
		*target = append(*target, int(i))
	}
	return nil
}

func parseJSON_Example2_Bar_Qux(target **field_Example2_Bar_Qux, d *jsontext.Decoder) error {
	pos := d.InputOffset()
	t, err := d.ReadToken()
	if err != nil {
		return fmt.Errorf("error at offset %d: %w", pos, err)
	}
	if t.Kind() != '{' {
		return fmt.Errorf("error at offset %d: expected '{', got %v", pos, t.Kind())
	}
	*target = &field_Example2_Bar_Qux{}
	for {
		pos = d.InputOffset()
		t, err := d.ReadToken()
		if err != nil {
			return fmt.Errorf("error at offset %d: %w", pos, err)
		}
		switch t.Kind() {
		default:
			return fmt.Errorf("error at offset %d: expected '\"' or '}', got %v", pos, t.Kind())
		case '}':
			return nil
		case '"':
		}
		switch t.String() {
		case "quux":
			err = parseJSON_Example2_Bar_Qux_Quux(&(*target).Quux, d)
		default:
			err = jsonutil.Discard(d)
		}
		if err != nil {
			return err
		}
	}
}

func parseJSON_Example2_Bar_Qux_Quux(target **field_Example2_Bar_Qux_Quux, d *jsontext.Decoder) error {
	pos := d.InputOffset()
	t, err := d.ReadToken()
	if err != nil {
		return fmt.Errorf("error at offset %d: %w", pos, err)
	}
	if t.Kind() != '0' {
		return fmt.Errorf("error at offset %d: expected number, got %v", pos, t.Kind())
	}
	f := t.Float()
	*target = &f
	return nil
}
