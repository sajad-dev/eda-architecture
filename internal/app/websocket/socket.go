package websocket

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/sajad-dev/eda-architecture/internal/app/exception"
	"github.com/sajad-dev/eda-architecture/internal/database/model"
)

var ActiveSocket *Websocket

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

func (ws *Websocket) AddAddr(pattern string) {
	ws.ServerMux.Mux.Handle("/app"+pattern, ws.Middleware([]MiddlewareFuncType{}, HandelFuncType(HandlerFunc)))
}

func (ws *Websocket) RunServer() {
	err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("WEBSOCKET_PORT")), ws.ServerMux.Mux)
	exception.Log(err)
}

func getAddress() model.GetOutput {
	ou :=  model.Get([]string{"public_key"}, "channels", []model.Where_st{}, "id", true)
	return ou
}

func Handler() {

	csm := NewCustomServeMux()
	ws := Websocket{MiddlewareBase: []MiddlewareFuncType{},
		Subscriber: map[string][]*websocket.Conn{},
		ServerMux:  *csm}

	for _, addr := range getAddress() {
		ws.AddAddr(
			"/" + addr["public_key"])
	}

	ActiveSocket = &ws
	go func() {
		ws.RunServer()
	}()

}
