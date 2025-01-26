package websocket

import (
	"net/http"

	"github.com/sajad-dev/eda-architecture/internal/app/exception"
)

func HandlerFunc(w http.ResponseWriter, r *http.Request, ws *Websocket) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	exception.Log(err)
	secret := r.Header.Get("secret_key")



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
	if !checkPrivateKey(secret, r.URL.Path[5:]) {
		return
	}
	go func() {
		defer conn.Close()
		for {
			
			_, msg, err := conn.ReadMessage()
			if err != nil {
				exception.Log(err)

				i := len(ws.Subscriber[r.URL.Path])
				ws.Subscriber[r.URL.Path] = ws.Subscriber[r.URL.Path][:i-1]
				conn.Close()
				break
			}
			for _, subscriber := range ws.Subscriber[r.URL.Path] {
				if subscriber != conn {
					err := subscriber.WriteMessage(1, msg)
					if err!= nil {
						
					}
				}
			}
		}
	}()
}

type AddChannel struct {
	Name string `json:"name"`
}

// func AddSocketChannel(w http.ResponseWriter, r *http.Request, ws *Websocket) {
// 	w.Header().Set("Content-Type", "application/json")

// 	var addChannel AddChannel
// 	err := json.NewDecoder(r.Body).Decode(&addChannel)
// 	exception.Log(err)
// 	ws.AddAddr(WebSocketHandler(ws, HandlerFunc), fmt.Sprintf("/%s", addChannel.Name))

// }
