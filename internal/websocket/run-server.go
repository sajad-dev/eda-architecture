package websocket

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sajad-dev/eda-architecture/internal/exception"
)

type Websocket struct {
	middlewareBase []MiddlewareFuncType
	serverMux      *CustomServeMux
}

type Addr struct {
	Pattern        string
	Handler        HandelFuncType
	MiddlewareList []MiddlewareFuncType
}

type CustomServeMux struct {
	mux *http.ServeMux
	mu  sync.RWMutex
}

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewCustomServeMux() *CustomServeMux {
	return &CustomServeMux{
		mux: http.NewServeMux(),
	}
}

type HandelFuncType func(w http.ResponseWriter, r *http.Request)
type MiddlewareFuncType func(http.Handler) http.Handler

func (ws *Websocket) Middleware(middlewareList []MiddlewareFuncType, handelFunc HandelFuncType) http.Handler {
	ws.serverMux.mu.Lock()
	defer ws.serverMux.mu.Unlock()
	var handler http.Handler
	handler = http.HandlerFunc(handelFunc)
	for _, middleware := range ws.middlewareBase {
		handler = middleware(handler)
	}
	for _, middleware := range middlewareList {
		handler = middleware(handler)
	}
	return handler
}

func (ws *Websocket) AddAddr(handlerFunc http.Handler, pattern string) {
	ws.serverMux.mu.Lock()
	defer ws.serverMux.mu.Unlock()
	ws.serverMux.mux.Handle(pattern, handlerFunc)
}

func (ws *Websocket) RunServer() {
	go func() {

		err := http.ListenAndServe(":8080", nil)
		exception.Log(err)
	}()
}

func Handler(addrs []Addr) Websocket {
	csm := NewCustomServeMux()
	ws := Websocket{middlewareBase: []MiddlewareFuncType{UpgraderMiddleware},
		serverMux: csm}

	for _, addr := range addrs {
		ws.AddAddr(ws.Middleware(addr.MiddlewareList, addr.Handler),
			addr.Pattern)
	}

	ws.RunServer()

	return ws
}
