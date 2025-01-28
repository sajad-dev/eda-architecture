package websocket

import (
	"net/http"
	"regexp"

	// "encoding/json"
	"github.com/gorilla/websocket"
	"github.com/sajad-dev/eda-architecture/internal/app/exception"
)

func (ws *Websocket) unsubscribe(public_key string, conn *websocket.Conn, message Message) {
	ws.Clients[public_key][message.Channel] = removeClient(ws.Clients[public_key][message.Channel], conn)

}

func (ws *Websocket) publish(public_key string, secret_key string, message Message) {
	if !checkPrivateKey(secret_key, public_key) {
		return
	}
	for _, sub := range ws.Clients[public_key][message.Channel] {
		err := sub.WriteJSON(message)
		exception.Log(err)
	}
}

func (ws *Websocket) subscribe(public_key string, conn *websocket.Conn, message Message) {
	ws.Clients[public_key][message.Channel] = append(ws.Clients[public_key][message.Channel], conn)
}

func (ws *Websocket) HandlerFunc(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	exception.Log(err)

	regex := regexp.MustCompile(`(?<=/app/)[^/?]+`)

	match := regex.FindString(r.URL.Path)
	secret := r.Header.Get("X_Secret_Key")

	go func() {
		defer conn.Close()
		for {
			var message Message
			err := conn.ReadJSON(&message)
			if err != nil {
				exception.Log(err)
				break
			}

			switch message.Event {
			case "subscribe":
				ws.subscribe(match, conn, message)
			case "unsubscribe":
				ws.unsubscribe(match, conn, message)
			case "publish":
				ws.publish(match, secret, message)
			}

		}
	}()
}
