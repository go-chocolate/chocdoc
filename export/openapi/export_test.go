package openapi

import (
	"testing"

	"github.com/go-chocolate/chocdoc/internal/doc"
)

func TestExport(t *testing.T) {
	var documents = doc.Documents{
		{
			Path:        "/api/example/:id",
			Name:        "example",
			Summary:     "example get",
			Description: "",
			Method:      "get",
			Req: &doc.Model{
				Name: "Request",
				Fields: []*doc.Field{
					{
						Name:     "Name",
						Type:     "string",
						Array:    1,
						Required: false,
						Comment:  "namenamename",
						Tags: doc.NewKV(map[string]string{
							"query": "name",
						}),
					},
				},
				Array: 0,
			},
			Rsp: &doc.Model{
				Name: "Response",
				Fields: []*doc.Field{
					{Name: "Name", Type: "string", Tags: doc.NewKV(map[string]string{"json": "name"})},
					{Name: "Value", Type: "int", Tags: doc.NewKV(map[string]string{"json": "value"})},
				},
				Array: 0,
			},
			KV: doc.NewKV(map[string]string{"ResponseExtend": "response data"}),
		},
	}

	schemas := map[string]*Schema{"response": NewSchema(doc.DecodeModel(&Response{}), "json")}
	swagger, err := Export(documents, WithSchemas(schemas))
	if err != nil {
		t.Error(err)
	}
	t.Log(swagger.JSON())
}

type Response struct {
	Code    int    `json:"code"`
	Data    string `json:"data"`
	Message string `json:"message"`
}

type Info struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func TestSchemaExtend(t *testing.T) {
	response := NewSchema(doc.DecodeModel(&Response{}), "json")
	info := NewSchema(doc.DecodeModel(&Info{}), "json")

	t.Log(response)
	response.Extend(info, "data")
	t.Log(response)
}
