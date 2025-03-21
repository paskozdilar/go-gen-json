// Example of a more complex structure, with nested objects, pointers and
// arrays.

package examples

import (
	"bytes"
	"fmt"
	"math"

	"github.com/go-json-experiment/json/jsontext"
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

func (target *Example2) ParseJSON(b []byte) error {
	// Use jsontext to parse tokens
	d := jsontext.NewDecoder(bytes.NewReader(b))

	// It's an object, so expect '{'
	pos := d.InputOffset()
	t, err := d.ReadToken()
	if err != nil {
		return fmt.Errorf("error at offset %d: %w", pos, err)
	}
	if t.Kind() != '{' {
		return fmt.Errorf("error at offset %d: expected '{', got %v",
			pos, t.Kind())
	}

	// Parse the fields, throw error on unknown fields (for now)
	filled := [2]bool{false, false}
	for range 2 {
		// Read the field name
		pos = d.InputOffset()
		t, err := d.ReadToken()
		if err != nil {
			return fmt.Errorf("error at offset %d: %w", pos, err)
		}
		if t.Kind() != '"' {
			return fmt.Errorf("error at offset %d: expected '\"', got %v",
				pos, t.Kind())
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
				return fmt.Errorf("error at offset %d: expected string, got %v",
					pos, t.Kind())
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
			if t.Kind() != '{' {
				return fmt.Errorf("error at offset %d: expected '{', got %v",
					pos, t.Kind())
			}
			// Parse the fields, throw error on unknown fields (for now)
			filled := [2]bool{false, false}
			for range 2 {
				pos = d.InputOffset()
				t, err := d.ReadToken()
				if err != nil {
					return fmt.Errorf("error at offset %d: %w", pos, err)
				}
				if t.Kind() != '"' {
					return fmt.Errorf("error at offset %d: expected '\"', got %v",
						pos, t.Kind())
				}
				// Read the field value
				switch t.String() {
				case "baz":
					// Check if value already filled
					if filled[0] {
						return fmt.Errorf("error at offset %d: duplicate field 'baz'", pos)
					}
					filled[0] = true
					// Read value
					pos := d.InputOffset()
					t, err := d.ReadToken()
					if err != nil {
						return fmt.Errorf("error at offset %d: %w", pos, err)
					}
					if t.Kind() != '[' {
						return fmt.Errorf("error at offset %d: expected '[', got %v",
							pos, t.Kind())
					}
					// Parse the array
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
						// Append the value
						i := t.Int()
						if i < math.MinInt || i > math.MaxInt {
							return fmt.Errorf("error at offset %d: number out of range: %d", pos, i)
						}
						target.Bar.Baz = append(target.Bar.Baz, int(i))
					}
				case "qux":
					// Check if value already filled
					if filled[1] {
						return fmt.Errorf("error at offset %d: duplicate field 'qux'", pos)
					}
					filled[1] = true
					var qux struct {
						Quux *float64 `json:"quux"`
					}
					target.Bar.Qux = &qux
					// Read value
					pos := d.InputOffset()
					t, err := d.ReadToken()
					if err != nil {
						return fmt.Errorf("error at offset %d: %w", pos, err)
					}
					if t.Kind() != '{' {
						return fmt.Errorf("error at offset %d: expected '{', got %v",
							pos, t.Kind())
					}
					// Parse the fields, throw error on unknown fields (for now)
					filled := [1]bool{false}
					for range 1 {
						pos = d.InputOffset()
						t, err := d.ReadToken()
						if err != nil {
							return fmt.Errorf("error at offset %d: %w", pos, err)
						}
						if t.Kind() != '"' {
							return fmt.Errorf("error at offset %d: expected '\"', got %v",
								pos, t.Kind())
						}
						// Read the field value
						switch t.String() {
						case "quux":
							// Check if value already filled
							if filled[0] {
								return fmt.Errorf("error at offset %d: duplicate field 'quux'", pos)
							}
							filled[0] = true
							// Read value
							pos := d.InputOffset()
							t, err := d.ReadToken()
							if err != nil {
								return fmt.Errorf("error at offset %d: %w", pos, err)
							}
							if t.Kind() != '0' {
								return fmt.Errorf("error at offset %d: expected number, got %v",
									pos, t.Kind())
							}
							// Set the field
							f := t.Float()
							target.Bar.Qux.Quux = &f
						default:
							return fmt.Errorf("error at offset %d: unknown field %v", pos, t.String())
						}
					}
				}
			}
		default:
			return fmt.Errorf("error at offset %d: unknown field %v", pos, t.String())
		}
	}
	return nil
}
