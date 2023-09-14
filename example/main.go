package main

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/go-chocolate/chocdoc/example/annotation"
	"github.com/go-chocolate/chocdoc/example/internal/router"
	"github.com/go-chocolate/chocdoc/internal/doc"
)

func main() {
	e := gin.New()

	router.Route(e)

	documents := doc.Decode(e, annotation.Nodes())
	b, _ := json.MarshalIndent(documents.Group(), "", "\t")
	fmt.Println(string(b))
}
