package openapi

import "github.com/go-chocolate/chocdoc/internal/doc"

func NewSchema(model *doc.Model, tag string) *Schema {
	schema := &Schema{Type: TypeObject, Properties: map[string]*Schema{}}
	buildSchema(schema, model, tag)
	return schema
}

func buildSchema(schema *Schema, m *doc.Model, tag string) (supported bool) {
	for _, v := range m.Fields {
		fieldName := v.Tags.Get(tag)
		if fieldName == "" {
			continue
		}
		supported = true
		var field *Schema
		var fieldType = formatGoType(v)
		if isArray, depth := arrayDepth(v); isArray {
			schema.Properties[fieldName] = &Schema{Type: fieldType}
			for i := 0; i < depth; i++ {
				schema.Properties[fieldName] = &Schema{Type: TypeArray, Items: schema.Properties[fieldName]}
			}
			field = schema.Properties[fieldName].Items
		} else {
			schema.Properties[fieldName] = &Schema{Type: fieldType, Properties: map[string]*Schema{}}
			field = schema.Properties[fieldName]
		}

		switch field.Type {
		case TypeObject:
			field.Properties = make(map[string]*Schema)
			buildSchema(field, v.Sub, tag)
		default:
		}
	}
	return
}
