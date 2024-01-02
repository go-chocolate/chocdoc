package openapi

import (
	"github.com/go-chocolate/chocdoc/chocdoc"
)

type Option func(e *exporter)

func WithInformation(info Information) Option {
	return func(e *exporter) {
		e.information = info
	}
}

type Exporter interface {
	Export(documents chocdoc.Documents) (string, error)
}

type exporter struct {
	information Information
}

func (e *exporter) Export(documents chocdoc.Documents) (string, error) {
	result := Export(documents, e.information)
	return result, nil
}

func NewExporter(options ...Option) Exporter {
	e := &exporter{information: Information{
		Title:   "Chocdoc",
		Version: "0.0.1",
	}}
	for _, opt := range options {
		opt(e)
	}
	return e
}
