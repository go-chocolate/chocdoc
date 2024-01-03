package doc

import (
	"net/http"
	"reflect"
	"runtime"
)

func DecodeHTTPMux(mux *http.ServeMux) []*Router {
	if mux == nil {
		mux = http.DefaultServeMux
	}
	muxMap := reflect.ValueOf(mux).Elem().FieldByName("m")

	var routers []*Router
	for _, key := range muxMap.MapKeys() {
		router := &Router{}

		val := muxMap.MapIndex(key)
		router.Path = val.FieldByName("pattern").String()

		handler := val.FieldByName("h").Elem()
		if handler.Kind() == reflect.Pointer {
			handler = handler.Elem()
		}

		switch handler.Type().Kind() {
		case reflect.Func:
			router.Type = Function
			router.Name = runtime.FuncForPC(handler.Pointer()).Name()
		case reflect.Struct:
			router.Type = Struct
			router.Name = handler.Type().PkgPath() + "." + handler.Type().Name()
		}
		routers = append(routers, router)
	}
	return routers
}
