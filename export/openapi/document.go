package openapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-chocolate/chocdoc/chocdoc"
	"github.com/go-chocolate/chocdoc/internal/utils/jsonutil"
)

const (
	VERSION = "3.0.1"
)

var pathParameterReg = regexp.MustCompile("\\{\\w+\\}")

type Information struct {
	Title       string   `json:"title"`
	Version     string   `json:"version"`
	Description string   `json:"description,omitempty"`
	License     *License `json:"license,omitempty"`
}

type License struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func (i *Information) JSON() string {
	b, _ := json.Marshal(i)
	return string(b)
}

func Export(documents chocdoc.Documents, information Information) string {
	builder := jsonutil.NewJsonBuilder()
	builder.WriteField("openapi", VERSION)
	builder.WriteJSON("info", information.JSON())
	builder.WriteJSON("paths", buildPaths(documents))
	return builder.String()
}

func split(documents chocdoc.Documents) map[string][]*chocdoc.Document {
	var results = make(map[string][]*chocdoc.Document)
	for _, doc := range documents {
		results[doc.Path] = append(results[doc.Path], doc)
	}
	return results
}

func buildPaths(documents chocdoc.Documents) string {
	builder := jsonutil.NewJsonBuilder()
	for path, docs := range split(documents) {
		builder.WriteJSON(path, buildRouters(docs))
	}
	return builder.String()
}

func buildRouters(documents []*chocdoc.Document) string {
	builder := jsonutil.NewJsonBuilder()
	for _, doc := range documents {
		builder.WriteJSON(strings.ToLower(doc.Method), buildDocument(doc))
	}
	return builder.String()
}

func buildDocument(document *chocdoc.Document) string {
	builder := jsonutil.NewJsonBuilder()
	builder.WriteField("summary", document.Summary)
	builder.WriteField("description", document.Description)
	if document.Group != "" {
		builder.WriteJSON("tags", fmt.Sprintf("[\"%s\"]", document.Group))
	}
	var parameters []string
	for _, pathParameter := range pathParameterReg.FindAllString(document.Path, -1) {
		parameters = append(parameters, (&Parameter{
			Name:     strings.Trim(pathParameter, "{}"),
			In:       "path",
			Required: true,
			Schema:   &JSONSchema{Type: "string"},
		}).JSON())
	}

	if document.Req != nil {
		for _, field := range document.Req.Fields {
			if name := field.Tags.Get("query"); name != "" {
				parameters = append(parameters, (&Parameter{
					Name:        name,
					In:          "query",
					Description: field.Comment,
					Required:    field.Required,
					Schema:      &JSONSchema{Type: formatGoType(field)},
				}).JSON())
			}
		}
	}

	builder.WriteJSONArray("parameters", parameters...)

	if hasRequestBody(document) && document.Req != nil {
		requestSchema := &JSONSchema{}
		if document.Req.Array > 0 {
			//TODO
		} else {
			requestSchema.Type = TypeObject
			requestSchema.Properties = make(map[string]*JSONSchema)
		}

		buildJSONSchema(requestSchema, document.Req)

		var request = &RequestBody{Content: map[string]*Schema{
			"application/json": {Schema: requestSchema},
		}}
		builder.WriteJSON("requestBody", request.JSON())
	}

	if document.Rsp != nil {
		responseSchema := &JSONSchema{Type: TypeObject, Properties: make(map[string]*JSONSchema)}
		buildJSONSchema(responseSchema, document.Rsp)
		var response = map[string]*ResponseBody{
			"200": {
				Description: "成功",
				Content: map[string]*Schema{
					"application/json": {Schema: responseSchema},
				},
			},
		}
		responseJSONBytes, _ := json.Marshal(response)
		builder.WriteJSON("responses", string(responseJSONBytes))
	}
	return builder.String()
}

type Parameter struct {
	Name        string      `json:"name"`
	In          string      `json:"in"`
	Description string      `json:"description"`
	Required    bool        `json:"required"`
	Schema      *JSONSchema `json:"schema"`
}

func (p *Parameter) JSON() string {
	b, _ := json.Marshal(p)
	return string(b)
}

func formatGoType(field *chocdoc.Field) Type {
	switch strings.TrimLeft(field.Type, "[]") {
	case "string":
		return TypeString
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64":
		return TypeInteger
	case "bool":
		return TypeBoolean
	case "float32", "float64":
		return TypeDouble
	default:
		return TypeObject
	}
}

func arrayDepth(field *chocdoc.Field) (bool, int) {
	if strings.HasPrefix(field.Type, "[]") {
		return true, strings.Count(field.Type, "[]")
	}
	return false, 0
}

type Type string

const (
	TypeUnknown Type = "unknown"
	TypeObject  Type = "object"
	TypeString  Type = "string"
	TypeInteger Type = "integer"
	TypeArray   Type = "array"
	TypeBoolean Type = "boolean"
	TypeDouble  Type = "double"
)

func (t Type) String() string {
	return string(t)
}

type JSONSchema struct {
	Type       Type                   `json:"type"`
	Items      *JSONSchema            `json:"items,omitempty"`
	Properties map[string]*JSONSchema `json:"properties,omitempty"`
}

type Schema struct {
	Schema *JSONSchema `json:"schema,omitempty"`
}

func (v *JSONSchema) JSON() string {
	b, _ := json.Marshal(v)
	return string(b)
}

func buildJSONSchema(schema *JSONSchema, m *chocdoc.Model) {
	for _, v := range m.Fields {
		name := v.Name
		if jname := v.Tags.Get("json"); jname != "" {
			name = strings.Split(jname, ",")[0]
		}

		var field *JSONSchema

		if isArray, depth := arrayDepth(v); isArray {
			schema.Properties[name] = &JSONSchema{Type: formatGoType(v)}
			for i := 0; i < depth; i++ {
				schema.Properties[name] = &JSONSchema{Type: TypeArray, Items: schema.Properties[name]}
			}
			field = schema.Properties[name].Items
		} else {
			schema.Properties[name] = &JSONSchema{Type: formatGoType(v), Properties: map[string]*JSONSchema{}}
			field = schema.Properties[name]
		}

		switch field.Type {
		case TypeObject:
			field.Properties = make(map[string]*JSONSchema)
			buildJSONSchema(field, v.Sub)
		default:
		}
	}
}

type RequestBody struct {
	Content map[string]*Schema `json:"content"`
	//Example
}

func (v *RequestBody) JSON() string {
	b, _ := json.Marshal(v)
	return string(b)
}

type ResponseBody struct {
	Description string             `json:"description"`
	Content     map[string]*Schema `json:"content"`
}

func hasRequestBody(document *chocdoc.Document) bool {
	switch strings.ToUpper(document.Method) {
	case http.MethodGet, http.MethodHead, http.MethodDelete:
		return false
	}
	return true
}
