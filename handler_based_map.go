package go_web

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
)

type HandlerBasedOnMap struct {
	// key = method + uri
	//handlers map[string]func(*Context)
	handlers sync.Map
}

func (h *HandlerBasedOnMap) ServeHTTP(ctx *Context) {
	// 计算 key
	request, writer := ctx.R, ctx.W

	handler, err := h.Match(request.Method, request.URL.Path)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		_, _ = writer.Write([]byte("route not found"))
		return
	}

	// 处理请求
	handler(NewContext(request, writer))

}

func (h *HandlerBasedOnMap) Route(method string, path string, handleFunc HandlerFunc) {
	key := h.key(method, path)

	// 重复的问题需要处理
	h.handlers.Store(key, handleFunc)
}

func (h *HandlerBasedOnMap) key(method string, path string) string {
	return fmt.Sprintf("%s#%s", method, path)
}

func (h *HandlerBasedOnMap) Match(method string, path string) (HandlerFunc, error) {
	key := h.key(method, path)
	handler, ok := h.handlers.Load(key)
	if !ok {
		return nil, errors.New("route not register")
	}
	return handler.(HandlerFunc), nil
}

var _ Handler = (*HandlerBasedOnMap)(nil)

func NewHandlerBasedOnMap() Handler {
	return &HandlerBasedOnMap{}
}
