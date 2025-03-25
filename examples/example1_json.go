package examples

import (
	"bytes"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/paskozdilar/go-gen-json/jsonutil"
)

type (
	field_Example1     = Example1
	field_Example1_Foo = string
	field_Example1_Bar = int
)

func (target *Example1) ParseJSON(b []byte) error {
	d := jsontext.NewDecoder(bytes.NewReader(b))
	return parseJSON_Example1(target, d)
}
