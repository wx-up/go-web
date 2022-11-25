package go_web

type HandlerFunc func(ctx *Context)

type Handler interface {
	ServeHTTP(ctx *Context)
	Route(method string, path string, handleFunc HandlerFunc)
	Match(method string, path string) (HandlerFunc, error)
}
