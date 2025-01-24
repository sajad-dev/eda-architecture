package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sajad-dev/eda-architecture/pkg/exception"
)

func HandlerFunc(w http.ResponseWriter, r *http.Request, ws *Websocket) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	exception.Log(err)
	// defer conn.Close()

	ws.ServerMux.Mu.Lock()
	ws.Subscriber[r.URL.Path] = append(ws.Subscriber[r.URL.Path], conn)
	ws.ServerMux.Mu.Unlock()
	for _, subscriber := range ws.Subscriber[r.URL.Path] {
		if subscriber != conn {
			err := subscriber.WriteMessage(1, []byte("new-subscriber"))
			if err != nil {
				conn.Close()
				exception.Log(err)
				break
			}
		}
	}
	go func() {
		defer conn.Close()
		for {

			_, msg, err := conn.ReadMessage()
			if err != nil {
				exception.Log(err)
				fmt.Println(ws.Subscriber)

				i := len(ws.Subscriber[r.URL.Path])
				ws.Subscriber[r.URL.Path] = ws.Subscriber[r.URL.Path][:i-1]
				conn.Close()
				fmt.Println(ws.Subscriber)
				break
			}
			fmt.Println(string(msg))
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
	err := json.NewDecoder(r.Body).Decode(&addChannel)
	exception.Log(err)
	ws.AddAddr(WebSocketHandler(ws, HandlerFunc), fmt.Sprintf("/%s", addChannel.Name))

}
