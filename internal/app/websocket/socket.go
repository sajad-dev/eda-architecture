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

func (ws *Websocket) AddAddr(pattern string) {
	ws.ServerMux.Mux.HandleFunc(pattern, ws.HandlerFunc)
}

func (ws *Websocket) runServer() {
	err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("WEBSOCKET_PORT")), ws.ServerMux.Mux)
	exception.Log(err)
}

func getAddress() model.GetOutput {
	ou := model.Get([]string{"public_key"}, "channels", []model.Where_st{}, "id", true)
	return ou
}

func Handler() {

	csm := NewCustomServeMux()
	ws := Websocket{
		Clients: map[string]map[string][]*websocket.Conn{},
		ServerMux:  *csm}

	for _, addr := range getAddress() {
		ws.AddAddr(
			"/app/" + addr["public_key"])
	}

	ActiveSocket = &ws
	go func() {
		ws.runServer()
	}()

}
