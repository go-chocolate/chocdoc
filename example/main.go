package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/go-chocolate/chocdoc/example/internal/router"
)

func main() {
	e := gin.New()

	router.Router(e)

	http.ListenAndServe(":8080", e)
}
