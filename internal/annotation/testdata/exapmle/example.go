package exapmle

import "net/http"

// Request
// @summary request example
type Request struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

// Response
type Response struct {
	Message string `json:"message"`
}

// Handle
// hahahaha
// @header content-type:application/json
// @request [Request]
// @response [Response]
func Handle(w http.ResponseWriter, r *http.Request) {

}
