package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/go-chocolate/chocdoc/example/internal/binding"
)

// HelloGET
// @summary HelloGET
// @description HelloGET
// @request [binding.HelloGetRequest]
// @response [binding.HelloGetResponse]
func HelloGET(ctx *gin.Context) {
	// if the package 'binding' is not used in current go file,
	// 'ann' cannot find the import path of the linked struct/function
	_ = new(binding.HelloGetRequest)

}

// HelloPost
// @summary HelloPost
// @description HelloGET
// @request [binding.HelloPostRequest]
// @response [binding.HelloPostResponse]
func HelloPost(ctx *gin.Context) {

}

// HelloDoc
// @summary HelloDoc
// @description HelloDoc
func HelloDoc(ctx *gin.Context) {

}
