package go_web

import "net/http"

type Server interface {
	Route(path string, handle func(ctx *Context))
	Run(addr string) error
}

type defaultServer struct {
}

func (d *defaultServer) Route(path string, handle func(ctx *Context)) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(r, w)
		handle(ctx)
	})
}

func (d *defaultServer) Run(addr string) error {
	return http.ListenAndServe(addr, nil)
}

var server = &defaultServer{}

func NewServer() Server {
	return server
}
