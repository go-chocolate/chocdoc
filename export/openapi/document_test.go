package openapi

import (
	"testing"

	"github.com/go-chocolate/chocdoc/chocdoc"
)

func TestExport(t *testing.T) {
	var documents = chocdoc.Documents{
		{
			Path:        "/api/example/{id}",
			Name:        "example",
			Summary:     "example get",
			Description: "",
			Method:      "get",
			Req: &chocdoc.Model{
				Name: "Request",
				Fields: []*chocdoc.Field{
					{
						Name:     "Name",
						Type:     "string",
						Array:    1,
						Required: false,
						Comment:  "namenamename",
						Tags: map[string][]string{
							"query": {"name"},
						},
					},
				},
				Array: 0,
			},
			Rsp: &chocdoc.Model{
				Name: "Response",
				Fields: []*chocdoc.Field{
					{
						Name:     "Code",
						Type:     "int",
						Required: true,
						Comment:  "code",
						Tags:     map[string][]string{"json": {"code"}},
					},
					{
						Name:     "Data",
						Type:     "Struct",
						Array:    0,
						Required: true,
						Comment:  "data",
						Option:   "",
						Sub: &chocdoc.Model{
							Name: "Data",
							Fields: []*chocdoc.Field{
								{Name: "Name", Type: "string", Tags: map[string][]string{"json": {"name"}}},
								{Name: "Value", Type: "int", Tags: map[string][]string{"json": {"value"}}},
							},
							Array: 1,
						},
						Tags: map[string][]string{"json": {"data"}},
					},
				},
				Array: 0,
			},
		},
	}
	result := Export(documents, Information{Title: "hello", Version: "0.0.1"})
	t.Log(result)
}

func TestBuildJsonSchema(t *testing.T) {
	var schema = &JSONSchema{
		Type:       TypeObject,
		Properties: make(map[string]*JSONSchema),
	}
	buildJSONSchema(schema, &chocdoc.Model{
		Name: "Demo",
		Fields: []*chocdoc.Field{
			{
				Name:  "Value",
				Type:  "[]int",
				Array: 1,
			},
			{
				Name:  "Data",
				Type:  "[]struct",
				Array: 1,
				Sub: &chocdoc.Model{
					Name:   "DemoSub",
					Fields: []*chocdoc.Field{{Name: "id", Type: "int"}},
					Array:  0,
				},
			},
		},
	})

	t.Log(schema.JSON())
}
