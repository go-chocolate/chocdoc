package openapi

import "github.com/go-chocolate/chocdoc/internal/doc"

type Exporter interface {
	Export(documents doc.Documents) (*Swagger, error)
}

type exporter struct {
	*builder
}

func (e *exporter) Export(documents doc.Documents) (*Swagger, error) {
	return e.build(documents)
}

func NewExporter(options ...Option) Exporter {
	return &exporter{builder: newBuilder(options...)}
}

func Export(documents doc.Documents, options ...Option) (*Swagger, error) {
	return NewExporter(options...).Export(documents)
}
