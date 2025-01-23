package websocket

import (
	"net/http"
	"os"

	"github.com/sajad-dev/eda-architecture/internal/exception"
)


func NewCustomServeMux() *CustomServeMux {
	return &CustomServeMux{
		Mux: http.NewServeMux(),
	}
}


func (ws *Websocket) Middleware(middlewareList []MiddlewareFuncType, handelFunc HandelFuncType) http.Handler {
	ws.ServerMux.Mu.Lock()
	defer ws.ServerMux.Mu.Unlock()
	var handler http.Handler
	handler = http.HandlerFunc(WebSocketHandler(ws, handelFunc))
	for _, middleware := range ws.MiddlewareBase {
		handler = middleware(handler)
	}
	for _, middleware := range middlewareList {
		handler = middleware(handler)
	}
	return handler
}

func WebSocketHandler(ws *Websocket, handlerFunc func(http.ResponseWriter, *http.Request, *Websocket)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlerFunc(w, r, ws)
	}
}

func (ws *Websocket) AddAddr(handlerFunc http.Handler, pattern string) {
	ws.ServerMux.Mu.Lock()
	defer ws.ServerMux.Mu.Unlock()
	ws.ServerMux.Mux.Handle(pattern, handlerFunc)
}

func (ws *Websocket) RunServer() {
	go func() {

		err := http.ListenAndServe(os.Getenv("PORT"), nil)
		exception.Log(err)
	}()
}

func Handler(addrs []Addr) Websocket {
	csm := NewCustomServeMux()
	ws := Websocket{MiddlewareBase: []MiddlewareFuncType{},
		ServerMux: csm}

	for _, addr := range addrs {
		ws.AddAddr(ws.Middleware(addr.MiddlewareList, HandelFuncType(addr.Handler)),
			addr.Pattern)
	}

	ws.RunServer()

	return ws
}
