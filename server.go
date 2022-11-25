package go_web

import "net/http"

type Server interface {
	Route(method string, path string, handle HandlerFunc)
	Run(addr string) error
}

type defaultServer struct {
	name    string
	handler Handler
	root    HandlerFunc
}

func (d *defaultServer) Route(method string, path string, handle HandlerFunc) {
	d.handler.Route(method, path, handle)
}

func (d *defaultServer) Run(addr string) error {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if _, err := d.handler.Match(request.Method, request.URL.Path); err != nil {
			writer.WriteHeader(http.StatusNotFound)
			_, _ = writer.Write([]byte("route not register"))
			return
		}
		ctx := NewContext(request, writer)
		d.root(ctx)
	})
	return http.ListenAndServe(addr, nil)
}

func NewHttpServer(name string, filters ...FilterBuilder) Server {
	handler := NewHandlerBasedOnMap()

	// 责任链
	root := handler.ServeHTTP
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
