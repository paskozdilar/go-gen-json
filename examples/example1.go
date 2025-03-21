// Example of a simple struct with two fields.

package examples

import (
	"bytes"
	"fmt"
	"io"
	"math"

	"github.com/go-json-experiment/json/jsontext"
)

// Input:
type Example1 struct {
	Foo string `json:"foo"`
	Bar int    `json:"bar"`
}

// Output:
func (target *Example1) ParseJSON(b []byte) error {
	// Use jsontext to parse tokens
	d := jsontext.NewDecoder(bytes.NewReader(b))

	// It's an object, so expect '{'
	pos := d.InputOffset()
	t, err := d.ReadToken()
	if err != nil {
		return fmt.Errorf("error at offset %d: %w", pos, err)
	}
	if t.Kind() != '{' {
		return fmt.Errorf("error at offset %d: expected '{', got %v", pos, t.Kind())
	}

	// Parse the fields, throw error on unknown fields (for now)
	filled := []bool{false, false}
	for i := 0; i < 2; i++ {
		// Read the field name
		pos := d.InputOffset()
		t, err := d.ReadToken()
		if err != nil {
			return err
		}
		if t.Kind() != '"' {
			return fmt.Errorf("error at offset %d: expected string, got %v", pos, t.Kind())
		}
		// Read the field value
		switch t.String() {
		case "foo":
			// Check if value already filled
			if filled[0] {
				return fmt.Errorf("error at offset %d: duplicate field 'foo'", pos)
			}
			filled[0] = true
			// Read value
			pos := d.InputOffset()
			t, err := d.ReadToken()
			if err != nil {
				return fmt.Errorf("error at offset %d: %w", pos, err)
			}
			if t.Kind() != '"' {
				return fmt.Errorf("error at offset %d: expected string, got %v", pos, t.Kind())
			}
			// Set the field
			target.Foo = t.String()
		case "bar":
			// Check if value already filled
			if filled[1] {
				return fmt.Errorf("error at offset %d: duplicate field 'bar'", pos)
			}
			filled[1] = true
			// Read value
			pos := d.InputOffset()
			t, err := d.ReadToken()
			if err != nil {
				return fmt.Errorf("error at offset %d: %w", pos, err)
			}
			if t.Kind() != '0' {
				return fmt.Errorf("error at offset %d: expected number, got %v", pos, t.Kind())
			}
			// Set the field
			i := t.Int()
			if i > math.MaxInt || i < math.MinInt {
				return fmt.Errorf("error at offset %d: number out of range: %d", pos, i)
			}
			target.Bar = int(i)
		default:
			return fmt.Errorf("error at offset %d: unknown field %s", pos, t.String())
		}
	}

	// Expect end of object
	pos = d.InputOffset()
	t, err = d.ReadToken()
	if err != nil {
		return fmt.Errorf("error at offset %d: %w", pos, err)
	}
	if t.Kind() != '}' {
		return fmt.Errorf("error at offset %d: expected '}', got %v", pos, t.Kind())
	}

	// Expect EOF
	pos = d.InputOffset()
	t, err = d.ReadToken()
	if err != io.EOF {
		return fmt.Errorf("error at offset %d: expected EOF, got %v", pos, t.Kind())
	}
	return nil
}
