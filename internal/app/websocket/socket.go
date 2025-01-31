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

// NewCustomServeMux creates a new instance of CustomServeMux
func NewCustomServeMux() *CustomServeMux {
	return &CustomServeMux{
		Mux: http.NewServeMux(),
	}
}

// AddAddr registers a WebSocket handler for a given pattern
func (ws *Websocket) AddAddr(pattern string) {
	ws.ServerMux.Mux.HandleFunc(pattern, ws.HandlerFunc)
}

// runServer starts the WebSocket server on the configured port
func (ws *Websocket) runServer() {
	err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("WEBSOCKET_PORT")), ws.ServerMux.Mux)
	exception.Log(err)
}

// getAddress retrieves public keys from the database
func getAddress() model.GetOutput {
	ou := model.Get([]string{"public_key"}, "channels", []model.Where_st{}, "id", true)
	return ou
}

// Handler initializes the WebSocket server and registers event handlers
func Handler() {
	csm := NewCustomServeMux()
	ws := Websocket{
		Clients:   map[string]map[string][]*websocket.Conn{},
		ServerMux: *csm,
	}

	// Register WebSocket handlers for each address retrieved from the database
	for _, addr := range getAddress() {
		ws.AddAddr(
			"/app/" + addr["public_key"])
	}
	ws.ServerMux.Mux.HandleFunc("/apps/local/events", ws.handleTriggerAPI)

	ActiveSocket = &ws
	go func() {
		ws.runServer()
	}()
}
