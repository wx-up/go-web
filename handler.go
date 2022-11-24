package go_web

import "net/http"

type Handler interface {
	http.Handler
	Route(method string, path string, handleFunc func(ctx *Context))
}
