package main

import (
	"os"

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
	os.Remove("godoc example.json")
	os.WriteFile("godoc example.json", []byte(jsonText), 0644)
}
