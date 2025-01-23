package websocket

import (
	"net/http"

	"github.com/sajad-dev/eda-architecture/internal/exception"
	websocket "github.com/sajad-dev/eda-architecture/internal/websocket"
)

func HandlerFunc(w http.ResponseWriter, r *http.Request, ws *websocket.Websocket) {
	conn, err := websocket.Upgrader.Upgrade(w, r, nil)
	exception.Log(err)
	defer conn.Close()

	
	ws.ServerMux.Mu.Lock()
	defer ws.ServerMux.Mu.Unlock()
	ws.Subscriber[r.URL.Path] = conn

	for {
		_, msg, err := conn.ReadMessage()
		exception.Log(err)
		for _, subscriber := range ws.Subscriber {
			err := subscriber.WriteMessage(1, msg)
			exception.Log(err)
		}
	}

}
