package go_web

import "net/http"

type Server interface {
	Route(method string, path string, handle func(ctx *Context))
	Run(addr string) error
}

type defaultServer struct {
	name     string
	handlers *HandlerBasedOnMap
}

func (d *defaultServer) Route(method string, path string, handle func(ctx *Context)) {
	key := d.handlers.key(method, path)

	// 重复的问题需要处理
	d.handlers.handlers[key] = handle
}

func (d *defaultServer) Run(addr string) error {
	return http.ListenAndServe(addr, d.handlers)
}

func NewServer(name string) Server {
	return &defaultServer{
		name:     name,
		handlers: NewHandlerBasedOnMap(),
	}
}
