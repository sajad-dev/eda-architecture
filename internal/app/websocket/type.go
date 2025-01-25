package websocket

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type HandelFuncType func(w http.ResponseWriter, r *http.Request, ws *Websocket)
type MiddlewareFuncType func(http.Handler) http.Handler

type Websocket struct {
	MiddlewareBase []MiddlewareFuncType
	ServerMux      CustomServeMux
	Subscriber     map[string][]*websocket.Conn
}

type Addr struct {
	Pattern        string
	Handler        HandelFuncType
	MiddlewareList []MiddlewareFuncType
}

type CustomServeMux struct {
	Mux *http.ServeMux
	Mu  sync.RWMutex
}

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
