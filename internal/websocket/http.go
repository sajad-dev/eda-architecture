package websocket

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
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

func (ws *Websocket) RunServer(waitGroup sync.WaitGroup) {
	go func() {
		defer waitGroup.Done()
		err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("WEBSOCKET_PORT")), ws.ServerMux.Mux)
		exception.Log(err)
	}()
}

func WebServer(waitGroup sync.WaitGroup, ws *Websocket) {
	
	go func() {
		
		defer waitGroup.Done()
		http.HandleFunc("/add-channel", WebSocketHandler(ws, AddSocketChannel))
		err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil)
		exception.Log(err)
	}()
}

func Handler(addrs []Addr) (Websocket, sync.WaitGroup) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(2)

	csm := NewCustomServeMux()
	ws := Websocket{MiddlewareBase: []MiddlewareFuncType{},
		Subscriber: map[string][]*websocket.Conn{},
		ServerMux:  *csm}

	for _, addr := range addrs {
		ws.AddAddr(ws.Middleware(addr.MiddlewareList, HandelFuncType(addr.Handler)),
			addr.Pattern)
	}

	ws.RunServer(waitGroup)
	WebServer(waitGroup, &ws)

	return ws, waitGroup
}
