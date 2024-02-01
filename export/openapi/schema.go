package openapi

import (
	"encoding/json"
	"errors"

	"github.com/go-chocolate/chocdoc/internal/doc"
)

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

func NewSchemaFromJSON(data []byte) (*Schema, error) {
	if len(data) == 0 {
		return nil, errors.New("invalid json input")
	}
	if data[0] == '{' && data[len(data)-1] == '}' {
		model := make(map[string]*Any)
		if err := json.Unmarshal(data, &model); err != nil {
			return nil, err
		}
		var schema = new(Schema)
		err := newSchemaFromJSON(schema, model)
		return schema, err
	} else if data[0] == '[' && data[len(data)-1] == ']' {
		var model []map[string]*Any
		if err := json.Unmarshal(data, &model); err != nil {
			return nil, err
		}
		var schema = new(Schema)
		err := newSchemaFromJSON(schema, model[0])
		if err != nil {
			return nil, err
		}
		return &Schema{Type: TypeArray, Items: schema}, nil
	} else {
		return nil, errors.New("invalid json input")
	}
}

func newSchemaFromJSON(schema *Schema, model map[string]*Any) error {
	if schema.Properties == nil {
		schema.Properties = map[string]*Schema{}
	}
	for k, v := range model {
		schema.Properties[k] = &Schema{Type: v.Type()}
		switch v.Type() {
		case TypeArray:
			schema.Properties[k].Items = new(Schema)
			array, err := v.Array()
			if err != nil {
				return err
			}
			if err = newSchemaFromJSON(schema.Properties[k].Items, array[0]); err != nil {
				return err
			}
		case TypeObject:
			object, err := v.Object()
			if err != nil {
				return err
			}
			if err = newSchemaFromJSON(schema.Properties[k], object); err != nil {
				return err
			}
		}
	}
	return nil
}

type Any struct {
	val string
}

func (a *Any) Type() Type {
	switch {
	case a.val[0] == '[':
		return TypeArray
	case a.val[0] == '{':
		return TypeObject
	case a.val[0] == '"':
		return TypeString
	case a.val == "true" || a.val == "false":
		return TypeBoolean
	case a.val == "null":
		return TypeNull
	case a.val[0] >= '0' && a.val[0] <= '9':
		return TypeNumber
	default:
		return TypeUnknown
	}
}

func (a *Any) Object() (map[string]*Any, error) {
	var object = map[string]*Any{}
	err := json.Unmarshal([]byte(a.val), &object)
	return object, err
}

func (a *Any) Array() ([]map[string]*Any, error) {
	var object []map[string]*Any
	err := json.Unmarshal([]byte(a.val), &object)
	return object, err
}

func (a *Any) UnmarshalJSON(b []byte) error {
	a.val = string(b)
	return nil
}
