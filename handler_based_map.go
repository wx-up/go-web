package go_web

import (
	"fmt"
	"net/http"
)

type HandlerBasedOnMap struct {
	// key = method + uri
	handlers map[string]func(*Context)
}

func (h *HandlerBasedOnMap) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// 计算 key
	key := h.key(request.Method, request.URL.Path)

	// 判断路由是否存在
	// method + 路由 才唯一
	handler, ok := h.handlers[key]
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		_, _ = writer.Write([]byte("route not found"))
		return
	}

	// 处理请求
	handler(NewContext(request, writer))

}

func (h *HandlerBasedOnMap) Route(method string, path string, handleFunc func(*Context)) {
	key := h.key(method, path)

	// 重复的问题需要处理
	h.handlers[key] = handleFunc
}

func (h *HandlerBasedOnMap) key(method string, path string) string {
	return fmt.Sprintf("%s#%s", method, path)
}

var _ Handler = (*HandlerBasedOnMap)(nil)

func NewHandlerBasedOnMap() Handler {
	return &HandlerBasedOnMap{
		handlers: make(map[string]func(*Context), 10),
	}
}
