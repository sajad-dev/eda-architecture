package websocket

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sajad-dev/eda-architecture/internal/app/exception"
)

func remove(slice []*websocket.Conn, paramToRemove *websocket.Conn) []*websocket.Conn {
	var result []*websocket.Conn
	for _, v := range slice {
		if v != paramToRemove {
			result = append(result, v)
		}
	}
	return result
}

func HandlerFunc(w http.ResponseWriter, r *http.Request, ws *Websocket) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	exception.Log(err)
	secret := r.Header.Get("secret_key")

	ws.ServerMux.Mu.Lock()
	ws.Subscriber[r.URL.Path] = append(ws.Subscriber[r.URL.Path], conn)
	ws.ServerMux.Mu.Unlock()

	conn.SetCloseHandler(func(code int, text string) error {
		fmt.Printf("Connection closed: Code=%d, Text=%s\n", code, text)
		ws.Subscriber[r.URL.Path] = remove(ws.Subscriber[r.URL.Path], conn)
		return nil
	})

	for _, subscriber := range ws.Subscriber[r.URL.Path] {
		if subscriber != conn {
			err := subscriber.WriteMessage(1, []byte("new-subscriber"))
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					fmt.Println("Connection closed normally.")
				} else {
					fmt.Println("Error reading message:", err)
				}
				ws.Subscriber[r.URL.Path] = remove(ws.Subscriber[r.URL.Path], conn)
				conn.Close()
				break
			}
		}
	}
	if !checkPrivateKey(secret, r.URL.Path[1:]) {
		return
	}
	go func() {
		defer conn.Close()
		for {

			_, msg, err := conn.ReadMessage()
			if err != nil {
				
				ws.Subscriber[r.URL.Path] = remove(ws.Subscriber[r.URL.Path], conn)
				conn.Close()
				break
			}
			for _, subscriber := range ws.Subscriber[r.URL.Path] {
				if subscriber != conn {
					err := subscriber.WriteMessage(1, msg)
					if err != nil {

					}
				}
			}
		}
	}()
}

type AddChannel struct {
	Name string `json:"name"`
}
