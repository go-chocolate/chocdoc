// Code generated by annotation. DO NOT EDIT.
// version: v0.0.1

package annotation

import(
    "github.com/go-chocolate/chocdoc/elements"

    "github.com/go-chocolate/chocdoc/example/internal/binding"
    "github.com/go-chocolate/chocdoc/example/internal/handler"
)

var nodes = map[string]*elements.Node{
    "github.com/go-chocolate/chocdoc/example/internal/handler.HelloGET": {
		Type: elements.TypeFunc,
		Name: "HelloGET",
		Ptr: handler.HelloGET,
		Path: "github.com/go-chocolate/chocdoc/example/internal/handler.HelloGET",
		Comments: []string{"HelloGET", "@summary HelloGET", "@description HelloGET", "@request [binding.HelloGetRequest]", "@response [binding.HelloGetResponse]"},
		Annotations: []*elements.Annotation{
            {
				Raw: "@summary HelloGET",
				Content: "summary HelloGET",
				Relation: []interface{}{},
			},
            {
				Raw: "@description HelloGET",
				Content: "description HelloGET",
				Relation: []interface{}{},
			},
            {
				Raw: "@request [binding.HelloGetRequest]",
				Content: "request [github.com/go-chocolate/chocdoc/example/internal/binding.HelloGetRequest]",
				Relation: []interface{}{new(binding.HelloGetRequest), },
			},
            {
				Raw: "@response [binding.HelloGetResponse]",
				Content: "response [github.com/go-chocolate/chocdoc/example/internal/binding.HelloGetResponse]",
				Relation: []interface{}{new(binding.HelloGetResponse), },
			},
		},
	},
    "github.com/go-chocolate/chocdoc/example/internal/handler.HelloPost": {
		Type: elements.TypeFunc,
		Name: "HelloPost",
		Ptr: handler.HelloPost,
		Path: "github.com/go-chocolate/chocdoc/example/internal/handler.HelloPost",
		Comments: []string{"HelloPost", "@summary HelloPost", "@description HelloGET", "@request [binding.HelloPostRequest]", "@response [binding.HelloPostResponse]"},
		Annotations: []*elements.Annotation{
            {
				Raw: "@summary HelloPost",
				Content: "summary HelloPost",
				Relation: []interface{}{},
			},
            {
				Raw: "@description HelloGET",
				Content: "description HelloGET",
				Relation: []interface{}{},
			},
            {
				Raw: "@request [binding.HelloPostRequest]",
				Content: "request [github.com/go-chocolate/chocdoc/example/internal/binding.HelloPostRequest]",
				Relation: []interface{}{new(binding.HelloPostRequest), },
			},
            {
				Raw: "@response [binding.HelloPostResponse]",
				Content: "response [github.com/go-chocolate/chocdoc/example/internal/binding.HelloPostResponse]",
				Relation: []interface{}{new(binding.HelloPostResponse), },
			},
		},
	},
    "github.com/go-chocolate/chocdoc/example/internal/handler.HelloDoc": {
		Type: elements.TypeFunc,
		Name: "HelloDoc",
		Ptr: handler.HelloDoc,
		Path: "github.com/go-chocolate/chocdoc/example/internal/handler.HelloDoc",
		Comments: []string{"HelloDoc", "@summary HelloDoc", "@description HelloDoc"},
		Annotations: []*elements.Annotation{
            {
				Raw: "@summary HelloDoc",
				Content: "summary HelloDoc",
				Relation: []interface{}{},
			},
            {
				Raw: "@description HelloDoc",
				Content: "description HelloDoc",
				Relation: []interface{}{},
			},
		},
	},
}

func Nodes() map[string]*elements.Node {
	return nodes
}
