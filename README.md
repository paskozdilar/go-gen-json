# go-gen-json

Generator for struct-specific JSON parsers / serializers in Go.

**THIS PACKAGE IS IN DEVELOPMENT AND NOT EVEN CLOSE TO READY TO USE**

## Introduction

One of the slowest parts of Go `json.Marshal` and `json.Unmarshal` functions
is the reflection part. This package aims to generate code that will be
specific to the struct you want to serialize / deserialize.

The plan is to use `//go:generate` directive to generate the code for the
structs you want to serialize / deserialize, similar to `stringer`.

E.g.:

```go
//go:generate go-gen-json -type=MyStruct
type MyStruct struct {
    Field1 string
    Field2 int
}
```

With the above directive, running `go generate` in the same directory as the
above file will generate a file `mystruct_json.go` with the following content:

```go
package mypackage

func (s *MyStruct) MarshalJSON() ([]byte, error) {
    // ...
}

func (s *MyStruct) UnmarshalJSON(data []byte) error {
    // ...
}
```

These will be entirely compatible with the `json.Marshal` and `json.Unmarshal`
functions, but will be much faster.

## Progress

- [x] Write a few examples manually to better understand the problem
- [x] Write benchmarks to prove the speedup of manually written code
- [ ] Write a parser for the struct definition
- [ ] Write a generator for the struct-specific code
    - [ ] MarshalJSON
    - [ ] UnmarshalJSON
- [ ] Write required utility functions in subpackage
    - [ ] Create subpackage "jsonutil"
    - [x] Write Discard function
