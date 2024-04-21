package handler

import (
	"github.com/dinizgab/golang-webserver/internal/request"
	"github.com/dinizgab/golang-webserver/internal/response"
)

type Handler struct {
	Path        string
	Method      string
	HandlerFunc HandleFunc
}

type HandleFunc func(*request.Request) *response.Response

func (f HandleFunc) Handle(req *request.Request) *response.Response {
	return f(req)
}

func New(path string, method string, handler func(*request.Request) *response.Response) *Handler {
	return &Handler{
		Path:        path,
		Method:      method,
		HandlerFunc: handler,
	}
}
