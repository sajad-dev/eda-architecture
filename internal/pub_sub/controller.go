package pub_sub

import (
	"fmt"
	"net/http"

	"github.com/sajad-dev/eda-architecture/internal/exception"
	"github.com/sajad-dev/eda-architecture/internal/websocket"
)

func HandlerFunc(w http.ResponseWriter, r *http.Request, ws *websocket.Websocket) {
	conn, err := websocket.Upgrader.Upgrade(w, r, nil)
	exception.Log(err)

	ws.ServerMux.Mu.Lock()
	ws.Subscriber[r.URL.Path] = append(ws.Subscriber[r.URL.Path], conn)
	ws.ServerMux.Mu.Unlock()

	fmt.Println(ws.Subscriber)

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
