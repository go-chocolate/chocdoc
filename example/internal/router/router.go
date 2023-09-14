package router

import (
	"github.com/gin-gonic/gin"

	"github.com/go-chocolate/chocdoc/example/internal/handler"
)

func Route(e *gin.Engine) {

	e.GET("/hello", handler.HelloGET)
	e.POST("/hello", handler.HelloPost)

	e.POST("/anonymous", func(ctx *gin.Context) {
	}, handler.HelloDoc)
}
