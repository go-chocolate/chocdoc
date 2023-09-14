package binding

type HelloGetRequest struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type HelloGetResponse struct {
}

type HelloPostRequest struct {
}

type HelloPostResponse struct {
}
