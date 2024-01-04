package openapi

import "strings"

type Option func(b *builder)

func WithInformation(info Information) Option {
	return func(b *builder) {
		b.info = info
	}
}

func WithSchema(name string, schema *Schema) Option {
	return func(b *builder) {
		b.schemas[strings.ToLower(name)] = schema
	}
}

func WithSchemas(schemas map[string]*Schema) Option {
	return func(b *builder) {
		for k, v := range schemas {
			b.schemas[strings.ToLower(k)] = v
		}
	}
}
