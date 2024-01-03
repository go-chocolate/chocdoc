package router

import (
	"github.com/gin-gonic/gin"

	"github.com/go-chocolate/chocdoc/example/internal/handler"
)

func Router(e gin.IRouter) {
	group := e.Group("/store")
	group.GET("/books", handler.Books)
	group.GET("/book/{id}", handler.BookDetail)
	group.POST("/book", handler.CreateBook)
	group.PUT("/book/{id}", handler.UpdateBook)
	group.DELETE("/book/{id}", handler.DeleteBook)
}
