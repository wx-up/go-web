package go_web

type Handler interface {
	ServeHTTP(ctx *Context)
	Route(method string, path string, handleFunc func(ctx *Context))
}
