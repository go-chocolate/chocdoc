package chocdoc

import (
	"github.com/go-chocolate/chocdoc/elements"
	"github.com/go-chocolate/chocdoc/internal/doc"
)

type (
	Document      = doc.Document
	Documents     = doc.Documents
	DocumentGroup = doc.DocumentGroup
	Model         = doc.Model
	Field         = doc.Field
	KV            = doc.KV
)

func Decode(engine any, annotations map[string]*elements.Node) Documents {
	return doc.Decode(engine, annotations)
}
