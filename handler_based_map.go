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

func (h *HandlerBasedOnMap) key(method string, path string) string {
	return fmt.Sprintf("%s#%s", method, path)
}

func NewHandlerBasedOnMap() *HandlerBasedOnMap {
	return &HandlerBasedOnMap{
		handlers: make(map[string]func(*Context), 10),
	}
}
