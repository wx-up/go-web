package go_web

import "net/http"

type Server interface {
	Route(method string, path string, handle func(ctx *Context))
	Run(addr string) error
}

type defaultServer struct {
	name    string
	handler Handler
	root    Filter
}

func (d *defaultServer) Route(method string, path string, handle func(ctx *Context)) {
	d.handler.Route(method, path, handle)
}

func (d *defaultServer) Run(addr string) error {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		ctx := NewContext(request, writer)
		d.root(ctx)
	})
	return http.ListenAndServe(addr, nil)
}

func NewHttpServer(name string, filters ...FilterBuilder) Server {
	handler := NewHandlerBasedOnMap()
	root := func(ctx *Context) {
		handler.ServeHTTP(ctx.W, ctx.R)
	}
	for i := len(filters) - 1; i >= 0; i-- {
		filter := filters[i]
		root = filter(root)
	}
	return &defaultServer{
		name:    name,
		handler: handler,
		root:    root,
	}
}
