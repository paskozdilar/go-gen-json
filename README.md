# go-gen-json

Generator for custom-tailored JSON parsers and serializers for Go structs.

**DISCLAIMER:** This package is still in development and not production-ready.

Currently only supports a small subset of Go types and does not handle
all edge cases.

Pending issues:
- [ ] Implement MarshalJSON
- [ ] Re-use existing `MarshalJSON` and `UnmarshalJSON` methods
- [ ] Parse recursive types
- [ ] Handle JSON struct tags
- [ ] Handle JSON options (omitempty, etc.)
- [ ] Handle unexported fields
- [ ] Handle external types (just call "json.Marshal" and "json.Unmarshal"?)

## Introduction

One of the slowest parts of Go `json.Marshal` and `json.Unmarshal` functions
is the reflection part. Fortunately, `json/v2` has also introduced `jsontext`
package which is a low-level JSON parser and serializer that does not use
reflection.

This package aims to generate custom-tailored code for marshaling and
unmarshaling Go structs to and from JSON using `jsontext` package.
This is now done using a custom `//go:generate` directive, similar to
`stringer` and `mockgen`.

## Example

Given the following Go code:

```go
//go:generate go-gen-json -type=MyStruct
type MyStruct struct {
    Field1 string
    Field2 int
}
```

Runnin `go generate` will generate a `mystruct_gen_json.go` file with the
following functions:

```go
func (s *MyStruct) MarshalJSON() ([]byte, error)
func (s *MyStruct) MarshalJSONTo(*jsontext.Writer)
func (s *MyStruct) UnmarshalJSON([]byte) error
func (s *MyStruct) UnmarshalJSONFrom(*jsontext.Reader) error
```

These will be compatible with the `json.Marshal` and `json.Unmarshal`
functions, so they can be used as drop-in replacements.
