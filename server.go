package go_web

import "net/http"

type Server interface {
	Route(method string, path string, handle func(ctx *Context))
	Run(addr string) error
}

type defaultServer struct {
	name    string
	handler Handler
}

func (d *defaultServer) Route(method string, path string, handle func(ctx *Context)) {
	d.handler.Route(method, path, handle)
}

func (d *defaultServer) Run(addr string) error {
	return http.ListenAndServe(addr, d.handler)
}

func NewHttpServer(name string) Server {
	return &defaultServer{
		name:    name,
		handler: NewHandlerBasedOnMap(),
	}
}
