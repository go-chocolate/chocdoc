package doc

import (
	"reflect"
	"runtime"

	"github.com/gin-gonic/gin"
)

func DecodeGin(e *gin.Engine) []*Router {
	var routers []*Router
	trees := reflect.ValueOf(e).Elem().FieldByName("trees")
	for i := 0; i < trees.Len(); i++ {
		tree := trees.Index(i)
		node := tree.FieldByName("root").Elem()
		method := tree.FieldByName("method").String()
		routers = append(routers, expandNode(method, node)...)
	}
	return routers
}

func expandNode(method string, node reflect.Value) []*Router {
	var routers []*Router
	fullPath := node.FieldByName("fullPath").String()
	children := node.FieldByName("children")
	for i := 0; i < children.Len(); i++ {
		child := children.Index(i).Elem()
		routers = append(routers, expandNode(method, child)...)
	}

	handlersChain := node.FieldByName("handlers")
	handlersLength := handlersChain.Len()
	if handlersLength == 0 {
		return routers
	}

	handler := handlersChain.Index(handlersLength - 1)
	handlerName := runtime.FuncForPC(handler.Pointer()).Name()
	router := new(Router)
	router.Path = fullPath
	router.Type = Function
	router.Name = handlerName
	router.Method = method
	return append(routers, router)
}
