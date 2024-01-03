Useage:

1. install chocdoc
```
go install github.com/go-chocolate/chocdoc/cmd/chocdoc@latest
```
2. write router and handler
```go
type HelloRequest struct{
	Name string `json:"name" doc:"must;the name balabala"`
}

type HelloResponse struct{
	Message string `json:"message" doc:"the message balabala"`
}

// Hello
// @summary say hello
// @description say hello
// @req [HelloRequest]
// @rsp [HelloResponse]
// @customKey customValue
func Hello(ctx *gin.Context) {
}

func route(r gin.IRouter) {
	r.POST("/hello", Hello)
}
```
3. run ```chocdoc``` to generate annotation code
```
./chocdoc
```
4. get documents by generated annotation, and export to openapi
```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-chocolate/chocdoc/chocdoc"
	"github.com/go-chocolate/chocdoc/export/openapi"

	example "github.com/go-chocolate/chocdoc/example/docs/chocdoc"
	"github.com/go-chocolate/chocdoc/example/internal/router"
)

//go:generate chocdoc -root ../ -output ./chocdoc
func main() {
	e := gin.New()
	router.Router(e)
	documents := chocdoc.Decode(e, example.Nodes())
	jsonText := openapi.Export(documents, openapi.Information{Title: "chocdoc example", Version: "0.0.1"})
}

```

see also [example](./example)