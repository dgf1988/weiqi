package h

import (
	"net/http"
)

//Handler 处理器接口, 和标准库的处理器比，多了一个参数args。
// 这个参数是路由器解析请求路径后得到，并传递给处理器。
type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, args []string)
}

//HandlerFunc 函数处理器
type HandlerFunc func(w http.ResponseWriter, r *http.Request, args []string)

func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request, args []string) {
	h(w, r, args)
}
