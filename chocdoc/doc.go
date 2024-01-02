package chocdoc

import (
	"github.com/gin-gonic/gin"

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

func Decode(engine *gin.Engine, annotations map[string]*elements.Node) Documents {
	return doc.Decode(engine, annotations)
}
