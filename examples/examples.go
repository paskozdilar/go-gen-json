// Package examples contains example types for testing go-gen-json.
package examples

import (
	"fmt"
	"time"
)

//go:generate go run .. -type=NamedString
type NamedString string

var (
	NamedStringValue = NamedString("foo")
	NamedStringJSON  = []byte(`"foo"`)
)

//go:generate go run .. -type=EmptyStruct
type EmptyStruct struct{}

var (
	EmptyStructValue = EmptyStruct{}
	EmptyStructJSON  = []byte(`{}`)
)

//go:generate go run .. -type=BasicStruct
type BasicStruct struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Email  string `json:"email"`
	Active bool   `json:"active"`
}

var (
	BasicStructValue = BasicStruct{
		Name:   "foo",
		Age:    42,
		Email:  "foo@bar.baz",
		Active: false,
	}
	BasicStructJSON = canonicalize([]byte(`
		{
			"name": "foo",
			"age": 42,
			"email": "foo@bar.baz",
			"active": false
		}
	`))
)

//go:generate go run .. -type=NestedStruct
type NestedStruct struct {
	ID      int           `json:"id"`
	Profile []BasicStruct `json:"profile"`
	Tags    []string      `json:"tags"`
}

var (
	NestedStructValue = NestedStruct{
		ID: 1234,
		Profile: []BasicStruct{
			BasicStructValue,
		},
		Tags: []string{"foo", "bar", "baz"},
	}
	NestedStructJSON = canonicalize(fmt.Appendf(nil, `
		{
			"id": 1234,
			"profile": [
				%s
			],
			"tags": [
				"foo",
				"bar",
				"baz"
			]
		}`,
		BasicStructJSON,
	))
)

//go:generate go run .. -type=ComplexStruct
type ComplexStruct struct {
	ID        int            `json:"id"`
	Data      map[string]any `json:"data"`
	Numbers   []float64      `json:"numbers"`
	Metadata  *BasicStruct   `json:"metadata,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
}

var (
	ComplexStructValue = ComplexStruct{
		ID: 1234,
		Data: map[string]any{
			"temperature": map[string]any{
				"Celsius":    37.0,
				"Fahrenheit": 98.6,
				"Kelvin":     310.15,
			},
			"humidity": 31.4,
			"notes":    "this data is completely made up",
		},
		Numbers: []float64{
			3.1415926535,
			2.7182818284,
			1.4142135623,
			1.6180339887,
			6.02214076e23,
			2.220446049250313e-16,
			0.1,
			0.0,
		},
		Metadata:  &BasicStructValue,
		CreatedAt: time.Date(2025, 9, 21, 17, 0, 0, 0, time.Local),
	}
	ComplexStructJSON = canonicalize(fmt.Appendf(nil, `
		{
		    "id": 1234,
		    "data": {
		        "humidity": 31.4,
		        "notes": "this data is completely made up",
		        "temperature": {
		            "Celsius": 37,
		            "Fahrenheit": 98.6,
		            "Kelvin": 310.15
		        }
		    },
		    "numbers": [
		        3.1415926535,
		        2.7182818284,
		        1.4142135623,
		        1.6180339887,
		        602214076000000000000000,
		        2.220446049250313e-16,
		        0.1,
		        0
		    ],
		    "metadata": %s,
		    "created_at": "2025-09-21T17:00:00+02:00"
		}`,
		BasicStructJSON,
	))
)

//go:generate go run .. -type=EmbeddedStruct
type EmbeddedStruct struct {
	BasicStruct `json:",inline"`
	NestedStruct
	ExtraField string `json:"extra_field"`
}

var (
	EmbeddedStructValue = EmbeddedStruct{
		BasicStruct: BasicStructValue,
		ExtraField:  "extra",
	}
	EmbeddedStructJSON = canonicalize(join([]byte(`
		{
			"extra_field": "extra"
		}`),
		BasicStructJSON,
	))
)

type TaggedStruct struct {
	PublicField     string `json:"public_field"`
	PrivateField    string `json:"-"`
	OmitEmpty       string `json:"omit_empty,omitempty"`
	CustomName      string `json:"custom_name"`
	unexportedField string
}

type InterfaceStruct struct {
	Value any `json:"value"`
}

type SliceStruct struct {
	StringSlice []string         `json:"string_slice"`
	IntSlice    []int            `json:"int_slice"`
	StructSlice []BasicStruct    `json:"struct_slice"`
	MapSlice    []map[string]any `json:"map_slice"`
}
