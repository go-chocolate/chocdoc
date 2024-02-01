package openapi

import (
	"strings"

	"github.com/go-chocolate/chocdoc/internal/doc"
)

const (
	VERSION = "3.0.1"
)

type builder struct {
	info    Information
	schemas map[string]*Schema
}

func newBuilder(options ...Option) *builder {
	b := &builder{info: Information{Title: "Chocdoc", Version: VERSION}, schemas: map[string]*Schema{}}
	for _, option := range options {
		option(b)
	}
	return b
}

func (b *builder) build(documents doc.Documents) (*Swagger, error) {
	swagger := &Swagger{Openapi: VERSION, Info: b.info}

	swagger.Paths = map[string]Path{}
	var err error
	for path, docs := range documents.Split() {
		if swagger.Paths[path], err = b.buildPath(docs); err != nil {
			return swagger, err
		}
	}
	return swagger, nil
}

func (b *builder) buildPath(documents []*doc.Document) (Path, error) {
	var p = Path{}
	for _, document := range documents {
		api, err := b.buildAPI(document)
		if err != nil {
			return p, err
		}
		p[Method(document.Method)] = api
	}
	return p, nil
}

func (b *builder) buildAPI(document *doc.Document) (*API, error) {
	api := &API{
		Summary:     document.Summary,
		Description: document.Description,
		Tags:        document.KV.Gets("tags"),
		Parameters:  b.buildRequestParameters(document),
	}

	b.buildRequestFormSchema(api, document)
	b.buildRequestJsonSchema(api, document)
	b.buildResponseSchema(api, document)

	return api, nil
}

func (b *builder) buildRequestParameters(document *doc.Document) []*Parameter {
	var parameters []*Parameter
	for _, parameter := range splitPathParameter(document.Path) {
		parameters = append(parameters, &Parameter{
			Name:     parameter,
			In:       InPath,
			Required: true,
			Schema:   &Schema{Type: TypeString},
		})
	}

	if document.Req != nil {
		for _, field := range document.Req.Fields {
			queryFieldName := field.Tags.Get("query")
			if queryFieldName == "" {
				continue
			}
			parameters = append(parameters, &Parameter{
				Name:     queryFieldName,
				In:       InQuery,
				Required: field.Required,
				Schema:   &Schema{Type: formatGoType(field)},
			})
		}
	}
	return parameters
}

func (b *builder) buildRequestFormSchema(api *API, document *doc.Document) {
	if document.Req == nil {
		return
	}
	var supportForm bool
	var supportFormFile bool
	for _, field := range document.Req.Fields {
		formFieldName := field.Tags.Get("form")
		if formFieldName != "" {
			supportForm = true
		}
		if field.Type == "multipart.File" || field.Type == "multipart.FileHeader" {
			supportFormFile = true
		}
	}
	if !supportForm {
		return
	}

	if api.RequestBody.Content == nil {
		api.RequestBody.Content = map[string]*Content{}
	}

	var body = NewSchema(document.Req, "form")
	if extend := strings.TrimSpace(document.KV.Get("RequestExtend")); extend != "" {
		var extendName, extendField string
		if n := strings.Index(extend, "."); n > 0 {
			extendName = strings.TrimSpace(extend[:n])
			extendField = strings.TrimSpace(extend[n+1:])
		} else {
			extendName, extendField = extend, "data"
		}
		if sc := b.schemas[extendName]; sc != nil {
			base := sc.Copy()
			base.Extend(body, extendField)
			body = base
		}
	}
	if supportFormFile {
		api.RequestBody.Content["multipart/form-data"] = &Content{Schema: body}
	} else {
		api.RequestBody.Content["x-www-urlencoded-form"] = &Content{Schema: body}
	}
}

func (b *builder) buildRequestJsonSchema(api *API, document *doc.Document) {
	if document.Req == nil {
		return
	}
	var supportJson bool
	for _, field := range document.Req.Fields {
		formFieldName := field.Tags.Get("json")
		if formFieldName != "" {
			supportJson = true
			break
		}
	}
	if !supportJson {
		return
	}

	if api.RequestBody.Content == nil {
		api.RequestBody.Content = map[string]*Content{}
	}

	var body = NewSchema(document.Req, "json")
	if extend := strings.TrimSpace(document.KV.Get("RequestExtend")); extend != "" {
		var extendName, extendField string
		if n := strings.Index(extend, " "); n > 0 {
			extendName = strings.TrimSpace(extend[:n])
			extendField = strings.TrimSpace(extend[n+1:])
		} else {
			extendName, extendField = extend, "data"
		}
		if sc := b.schemas[extendName]; sc != nil {
			base := sc.Copy()
			base.Extend(body, extendField)
			body = base
		}
	}
	api.RequestBody.Content["application/json"] = &Content{Schema: body}
}

func (b *builder) buildResponseSchema(api *API, document *doc.Document) {
	if document.Rsp == nil {
		//document.KV.Get("")
		//TODO
		return
	}
	if api.Responses == nil {
		api.Responses = map[string]*Body{}
	}
	var body = NewSchema(document.Rsp, "json")
	if extend := strings.TrimSpace(document.KV.Get("ResponseExtend")); extend != "" {
		var extendName, extendField string
		if n := strings.Index(extend, " "); n > 0 {
			extendName = strings.TrimSpace(extend[:n])
			extendField = strings.TrimSpace(extend[n+1:])
		} else {
			extendName, extendField = extend, "data"
		}
		if sc := b.schemas[extendName]; sc != nil {
			base := sc.Copy()
			base.Extend(body, extendField)
			body = base
		}
	}
	api.Responses["200"] = &Body{
		Description: "",
		Content:     map[string]*Content{"application/json": {Schema: body}},
	}
}
