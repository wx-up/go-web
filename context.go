package go_web

import "net/http"

// Context 一次完整的请求包含请求和响应，使用 context 来抽象这一过程
type Context struct {
	W http.ResponseWriter
	R *http.Request
}

func NewContext(r *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		W: w,
		R: r,
	}
}
