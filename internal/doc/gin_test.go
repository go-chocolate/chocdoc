package doc

import (
	"encoding/json"
	"testing"

	"github.com/gin-gonic/gin"
)

func GinExampleHelloGet(ctx *gin.Context)  {}
func GinExampleHelloPost(ctx *gin.Context) {}

func TestDecodeGin(t *testing.T) {
	eng := gin.Default()

	eng.GET("/hello", GinExampleHelloGet)
	eng.POST("/hello", GinExampleHelloPost)

	group := eng.Group("group")
	group.GET("/hello", GinExampleHelloGet)
	group.POST("/hello", GinExampleHelloPost)

	routers := DecodeGin(eng)
	b, _ := json.MarshalIndent(routers, "", "\t")
	t.Log(string(b))
}
