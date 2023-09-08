package doc

import (
	"reflect"
	"runtime"

	"github.com/gin-gonic/gin"
)

type GINRouter struct {
	Method   string
	Path     string
	FuncName string
	Group    string //TODO
}

func DecodeGin(e *gin.Engine) []*GINRouter {
	var routers []*GINRouter
	trees := reflect.ValueOf(e).Elem().FieldByName("trees")
	for i := 0; i < trees.Len(); i++ {
		tree := trees.Index(i)
		node := tree.FieldByName("root").Elem()
		method := tree.FieldByName("method").String()
		routers = append(routers, expandNode(method, node)...)
	}
	return routers
}

func expandNode(method string, node reflect.Value) []*GINRouter {
	var routers []*GINRouter
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
	router := new(GINRouter)
	router.Path = fullPath
	router.FuncName = handlerName
	router.Method = method
	return append(routers, router)
}
