package websocket

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type HandelFuncType func(w http.ResponseWriter, r *http.Request, ws *Websocket)
type MiddlewareFuncType func(http.Handler) http.Handler

type Message struct {
	Event   string `json"event"`
	Channel string `json"channel"`
	Data string `json"data"`
}

type Websocket struct {
	ServerMux CustomServeMux
	Clients   map[string]map[string][]*websocket.Conn
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
