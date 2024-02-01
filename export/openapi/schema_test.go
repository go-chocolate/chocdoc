package openapi

import (
	"testing"

	"github.com/go-chocolate/chocdoc/internal/doc"
)

type Response struct {
	Code    int    `json:"code"`
	Data    Data   `json:"data"`
	Message string `json:"message"`
}

type Data struct {
	Data string `json:"data"`
}

type Info struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func TestSchemaExtend(t *testing.T) {
	response := NewSchema(doc.DecodeModel(&Response{}), "json")
	info := NewSchema(doc.DecodeModel(&Info{}), "json")

	t.Log(response)
	t.Log(info)
	t.Log(response.Extend(info, "data.data"))
	t.Log(response)
}

func TestNewSchemaFromJSON(t *testing.T) {
	schema, err := NewSchemaFromJSON([]byte(`{"code":200,"data":[{"name":"1"}]}`))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(schema)
	t.Log(schema.String())
}
