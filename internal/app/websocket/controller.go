package websocket

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/gorilla/websocket"
	"github.com/sajad-dev/eda-architecture/internal/app/exception"
)

func (ws *Websocket) unsubscribe(public_key string, conn *websocket.Conn, message Message) {
	if ws.Clients[public_key] == nil {
		ws.Clients[public_key] = make(map[string][]*websocket.Conn)
	}
	ws.Clients[public_key][message.Channel] = removeClient(ws.Clients[public_key][message.Channel], conn)
	response := Message{
		Event:   "pusher_internal:unsubscribed",
		Channel: message.Channel,
		Data:    map[string]string{"status": "unsubscribed"},
	}
	err := conn.WriteJSON(response)
	exception.Log(err)
}

func (ws *Websocket) publish(public_key string, secret_key string, message Message) {
	if !checkPrivateKey(secret_key, public_key) {
		return
	}
	for _, sub := range ws.Clients[public_key][message.Channel] {
		response := Message{
			Event:   message.Event,
			Channel: message.Channel,
			Data:    message.Data,
		}
		err := sub.WriteJSON(response)
		exception.Log(err)
	}
}

func (ws *Websocket) subscribe(public_key string, conn *websocket.Conn, message Message) {

	if ws.Clients[public_key] == nil {
		ws.Clients[public_key] = make(map[string][]*websocket.Conn)
	}

	ws.Clients[public_key][message.Channel] = append(ws.Clients[public_key][message.Channel], conn)
	response := Message{
		Event:   "pusher_internal:subscription_succeeded",
		Channel: message.Channel,
		Data:    map[string]string{"status": "subscribed"},
	}
	err := conn.WriteJSON(response)
	exception.Log(err)
}

func (ws *Websocket) HandlerFunc(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	exception.Log(err)

	regex := regexp.MustCompile(`/app/([^/?]+)`)
	match := regex.FindStringSubmatch(r.URL.Path)

	secret := r.Header.Get("secret_key")

	go func() {
		defer conn.Close()
		for {
			var message Message
			message.Data = make(map[string]string)
			err := conn.ReadJSON(&message)
			if err != nil {
				fmt.Println(err)
				break
			}

			switch message.Event {
			case "subscribe":
				ws.subscribe(match[1], conn, message)
			case "unsubscribe":
				ws.unsubscribe(match[1], conn, message)
			default:
				ws.publish(match[1], secret, message)
			}

		}
	}()
}
