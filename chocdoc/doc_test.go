package chocdoc

import (
	"encoding/json"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/go-chocolate/chocdoc/elements"
)

func GinExampleHelloGet(ctx *gin.Context)  {}
func GinExampleHelloPost(ctx *gin.Context) {}

func TestDecode(t *testing.T) {

	eng := gin.Default()

	eng.GET("/hello", GinExampleHelloGet)
	eng.POST("/hello", GinExampleHelloPost)

	var doc Documents
	var group *DocumentGroup
	doc = Decode(eng, map[string]*elements.Node{})
	group = doc.Group()

	b, _ := json.MarshalIndent(group, "", "\t")
	t.Log(string(b))
}
