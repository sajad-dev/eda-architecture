package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type Websocket struct {
	middlewareBase []MiddlewareFuncType
}

type Addr struct {
	Pattern        string
	Handler        HandelFuncType
	MiddlewareList []MiddlewareFuncType
}

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type HandelFuncType func(w http.ResponseWriter, r *http.Request)
type MiddlewareFuncType func(http.Handler) http.Handler

func (w *Websocket) Middleware(middlewareList []MiddlewareFuncType, handelFunc HandelFuncType) http.Handler {
	var handler http.Handler
	handler = http.HandlerFunc(handelFunc)
	for _, middleware := range w.middlewareBase {
		handler = middleware(handler)
	}
	for _, middleware := range middlewareList {
		handler = middleware(handler)
	}
	return handler
}

func (w *Websocket) AddAddr(handlerFunc http.Handler, pattern string) {
	http.Handle("", handlerFunc)
}

func (w *Websocket) RunServer() {

	http.HandleFunc("", nil)
}

func Handler(addrs []Addr) {
	ws := Websocket{middlewareBase: []MiddlewareFuncType{UpgraderMiddleware}}

	for _, addr := range addrs {
		ws.AddAddr(ws.Middleware(addr.MiddlewareList, addr.Handler),
			addr.Pattern)
	}
	

	ws.RunServer()
}
