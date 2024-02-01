package handler

import "github.com/gin-gonic/gin"

type Book struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Author string `json:"author"`
}

type CreateBookRequest struct {
	Name   string `json:"name"`
	Author string `json:"author"`
}

type CreateBookResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// CreateBook
// @request [CreateBookRequest]
// @response [CreateBookResponse]
func CreateBook(ctx *gin.Context) {}

type UpdateBookRequest struct {
	Name   string `json:"name"`
	Author string `json:"author"`
}

type UpdateBookResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// UpdateBook
// @request [UpdateBookRequest]
// @response [UpdateBookResponse]
func UpdateBook(ctx *gin.Context) {}

type BooksResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []Book `json:"data"`
}

// Books
// @request
// @response [BooksResponse]
func Books(ctx *gin.Context) {}

// DeleteBook
// @response [UpdateBookResponse]
func DeleteBook(ctx *gin.Context) {}

type BookDetailResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    Book   `json:"data"`
}

// BookDetail
// @response [BookDetailResponse]
func BookDetail(ctx *gin.Context) {

}

type ExampleRequest struct{}
type ExampleResponse struct{}

// ExampleDoc
// @request [ExampleRequest]
// @response [ExampleResponse]
func ExampleDoc(ctx *gin.Context) {}
