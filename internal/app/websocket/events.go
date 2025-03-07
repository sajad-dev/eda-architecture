package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/sajad-dev/eda-architecture/internal/app/exception"
)

// eventHandel processes WebSocket events asynchronously
func (ws *Websocket) eventHandel(msg chan MessageChan) {
	select {
	case message := <-msg:
		switch message.message.Event {
		case "pusher:subscribe":
			ws.subscribe(message.public_key, message.connection, *message.message)
		}
	}
}

// publish sends messages to all subscribed clients
func (ws *Websocket) publish(public_key string, message TriggerBody) {
	for _, channel := range message.Channels {
		for _, sub := range ws.Clients[public_key][channel] {
			response := Message{
				Event:   message.Name,
				Channel: channel,
				Data:    message.Data,
			}

			err := sub.WriteJSON(response)
			exception.Log(err)
		}
	}
}

// subscribe adds a client to a specified channel for WebSocket communication
func (ws *Websocket) subscribe(public_key string, conn *websocket.Conn, message Message) {
	if ws.Clients[public_key] == nil {
		ws.Clients[public_key] = make(map[string][]*websocket.Conn)
	}

	data, _ := message.Data.(map[string]interface{})
	channel, _ := data["channel"].(string)
	ws.Clients[public_key][channel] = append(ws.Clients[public_key][channel], conn)

	response := Message{
		Event:   "pusher_internal:subscription_succeeded",
		Channel: channel,
		Data:    map[string]string{},
	}
	err := conn.WriteJSON(response)
	exception.Log(err)
}