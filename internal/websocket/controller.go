package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sajad-dev/eda-architecture/internal/exception"
)


func HandlerFunc(w http.ResponseWriter, r *http.Request, ws *Websocket) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	exception.Log(err)

	ws.ServerMux.Mu.Lock()
	ws.Subscriber[r.URL.Path] = append(ws.Subscriber[r.URL.Path], conn)
	ws.ServerMux.Mu.Unlock()

	go func() {

		for {
			_, msg, err := conn.ReadMessage()
			exception.Log(err)
			for _, subscriber := range ws.Subscriber[r.URL.Path] {
				if subscriber != conn {
					err := subscriber.WriteMessage(1, msg)
					exception.Log(err)
				}
			}
		}
	}()
}



type AddChannel struct {
	Name string `json:"name"`
}

func AddSocketChannel(w http.ResponseWriter, r *http.Request, ws *Websocket) {
	w.Header().Set("Content-Type", "application/json")

	var addChannel AddChannel
	json.NewDecoder(r.Body).Decode(&addChannel)
	ws.AddAddr(WebSocketHandler(ws, HandlerFunc), fmt.Sprintf("/%s", addChannel.Name))

}
