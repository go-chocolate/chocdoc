package doc

import (
	"encoding/json"
	"net/http"
	"testing"
)

func hello(w http.ResponseWriter, r *http.Request) {

}

type sayHello struct{}

func (h *sayHello) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func TestDecodeHTTPMux(t *testing.T) {
	http.HandleFunc("/hello", hello)
	http.Handle("/v2/hello", &sayHello{})

	for _, v := range DecodeHTTPMux(nil) {
		b, _ := json.Marshal(v)
		t.Log(string(b))
	}

}
