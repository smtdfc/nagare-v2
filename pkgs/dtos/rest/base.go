package rest

type ApiResponse[T any] struct {
	IsSuccess bool             `json:"is_success"`
	Data      T                `json:"data"`
	Error     ApiResponseError `json:"error"`
}

type ApiResponseError struct {
	Message string `json:"message"`
}

type ApiMethod string

const (
	GET  ApiMethod = "get"
	POST ApiMethod = "post"
)
