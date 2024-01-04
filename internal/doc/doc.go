package doc

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/go-chocolate/chocdoc/elements"
)

type HandlerType string

const (
	Function HandlerType = "function"
	Struct   HandlerType = "struct"
)

type Router struct {
	Method string
	Path   string
	Type   HandlerType
	Name   string
	Group  string //TODO
}

type document struct {
	path        string      // the url path
	method      string      // the http method
	handler     string      // the handler name
	name        string      // the document name
	summary     string      //
	description string      //
	header      KV          // headers
	group       string      // group
	request     interface{} // request model pointer
	response    interface{} // response model pointer
	kv          KV          //
}

func decode(docs []*document) Documents {
	var documents []*Document
	for _, v := range docs {
		if v.name == "" {
			//funcName := runtime.FuncForPC(reflect.ValueOf(v.Handler).Pointer()).Name()
			funcName := v.handler
			if tmp := strings.Split(funcName, "."); len(tmp) > 1 {
				v.name = tmp[len(tmp)-1]
			} else {
				v.name = funcName
			}
		}
		if v.header == nil {
			v.header = kvMap{}
		}
		if v.kv == nil {
			v.kv = kvMap{}
		}
		if v.group == "" {
			if n := strings.LastIndex(v.name, "/"); n >= 0 {
				v.group = v.name[:n]
				v.name = v.name[n+1:]
			}
		}
		doc := &Document{
			Path:        v.path,
			Name:        v.name,
			Summary:     v.summary,
			Description: v.description,
			Method:      v.method,
			Header:      v.header,
			KV:          v.kv,
			Req:         newDecoder(tree{}).decode(v.request),
			Rsp:         newDecoder(tree{}).decode(v.response),
			Group:       v.group,
		}
		documents = append(documents, doc)
	}
	return documents
}

func Decode(mux any, annotations map[string]*elements.Node) Documents {
	var routers []*Router
	switch m := mux.(type) {
	case *gin.Engine:
		routers = DecodeGin(m)
	case *http.ServeMux:
		routers = DecodeHTTPMux(m)
	default:
		panic(fmt.Errorf("unknown http router: %v", reflect.TypeOf(mux)))
	}
	docs := FromAnnotation(routers, annotations)
	return decode(docs)
}
