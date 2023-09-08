package doc

import (
	"encoding/json"
	"testing"
)

type RequestExample struct {
	Name  string `doc:"must;name example"`
	Value int
	Data  []*Data
	Child *RequestExample
}

type Data struct {
	ID     int
	Enable bool
}

func TestDecodeModel(t *testing.T) {
	var req = new(RequestExample)
	m := DecodeModel(req)
	b, _ := json.MarshalIndent(m, "", "\t")
	t.Log(string(b))
}
