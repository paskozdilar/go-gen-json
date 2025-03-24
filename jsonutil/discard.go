package jsonutil

import (
	"fmt"

	"github.com/go-json-experiment/json/jsontext"
)

// Discard consumes the next JSON value from the decoder without storing it.
func Discard(d *jsontext.Decoder) error {
	return discardAny(d)
}

func discardAny(d *jsontext.Decoder) error {
	pos := d.InputOffset()
	t, err := d.ReadToken()
	if err != nil {
		return fmt.Errorf("error at offset %d: %w", pos, err)
	}
	switch t.Kind() {
	case 'n', 'f', 't', '"', '0':
		return nil
	case '{':
		return discardObject(d)
	case '}':
		return fmt.Errorf("error at offset %d: unexpected '}'", pos)
	case '[':
		return discardArray(d)
	case ']':
		return fmt.Errorf("error at offset %d: unexpected ']'", pos)
	default:
		return fmt.Errorf("error at offset %d: undefined token kind: '%v'", pos, t.Kind())
	}
}

func discardArray(d *jsontext.Decoder) error {
	for {
		pos := d.InputOffset()
		t, err := d.ReadToken()
		if err != nil {
			return fmt.Errorf("error at offset %d: %w", pos, err)
		}
		if t.Kind() == ']' {
			return nil
		}
		if err := discardAny(d); err != nil {
			return err
		}
	}
}

func discardObject(d *jsontext.Decoder) error {
	for {
		pos := d.InputOffset()
		t, err := d.ReadToken()
		if err != nil {
			return fmt.Errorf("error at offset %d: %w", pos, err)
		}
		switch t.Kind() {
		case '}':
			return nil
		case '"':
			if err := discardAny(d); err != nil {
				return err
			}
		default:
			return fmt.Errorf("error at offset %d: expected '\"' or '}', got '%v'", pos, t.Kind())
		}
	}
}
