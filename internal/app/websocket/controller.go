package websocket

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sajad-dev/eda-architecture/internal/app/exception"
	"github.com/sajad-dev/eda-architecture/internal/database/model"
)

type msgS struct {
	conn       *websocket.Conn
	message    *Message
	public_key string
	sin        string
}

func (ws *Websocket) unsubscribe(public_key string, conn *websocket.Conn, message Message) {
	if ws.Clients[public_key] == nil {
		ws.Clients[public_key] = make(map[string][]*websocket.Conn)
	}
	for _, channel := range message.Channels {
		ws.Clients[public_key][channel] = removeClient(ws.Clients[public_key][channel], conn)
	}
	response := Message{
		Name:     "pusher_internal:unsubscribed",
		Channels: message.Channels,
		Data:     map[string]string{"status": "unsubscribed"},
	}
	err := conn.WriteJSON(response)
	exception.Log(err)
}

type Message2 struct {
	Event   string      `json:"event"`
	Channel string      `json:"channel,omitempty"`
	Data    interface{} `json:"data"`
}

func (ws *Websocket) publish(public_key string, message Message) {
	for _, channel := range message.Channels {
		for _, sub := range ws.Clients[public_key][channel] {
			response := Message2{
				Event:   message.Name, // تغییر `name` به `event`
				Channel: channel,      // تغییر `channels` به `channel`
				Data: map[string]string{
					"message": fmt.Sprintf("%v", message.Data), // تبدیل `data` به JSON Object
				},
			}

			err := sub.WriteJSON(response)
			exception.Log(err)
		}
	}
}

func (ws *Websocket) subscribe(public_key string, conn *websocket.Conn, message Message) {
	fmt.Println("HHHH")

	if ws.Clients[public_key] == nil {
		ws.Clients[public_key] = make(map[string][]*websocket.Conn)
	}

	for _, channel := range message.Channels {
		ws.Clients[public_key][channel] = append(ws.Clients[public_key][channel], conn)
	}

	response := Message{
		Name:     "pusher_internal:subscription_succeeded",
		Channels: message.Channels,
		Data:     map[string]string{"status": "subscribed"},
	}
	err := conn.WriteJSON(response)
	exception.Log(err)
}

func (ws *Websocket) HandlerFunc(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	exception.Log(err)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	regex := regexp.MustCompile(`/app/([^/?]+)`)
	match := regex.FindStringSubmatch(r.URL.Path)

	ws.Clients["d6kAd89bMqDrLrFh"] = make(map[string][]*websocket.Conn)
	ws.Clients["d6kAd89bMqDrLrFh"]["test"] = append(ws.Clients["d6kAd89bMqDrLrFh"]["test"], conn)
	go func() {
		for {
			time.Sleep(30 * time.Second)
			for _, clients := range ws.Clients["d6kAd89bMqDrLrFh"] {
				for _, conn := range clients {
					err := conn.WriteMessage(websocket.PingMessage, nil)
					if err != nil {
						fmt.Println("Ping failed:", err)
						removeClient(clients, conn)
					}
				}
			}
		}
	}()

	go func() {
		defer conn.Close()
		for {
			var message Message
			message.Data = make(map[string]string)
			_, ou, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println(string(ou))

			err = json.Unmarshal(ou, &message)
			exception.Log(err)

			var msg = make(chan msgS, 1)

			msg_ := msgS{conn: conn, public_key: match[1], message: &message}

			msg <- msg_

			ws.eventHandel(msg)

		}
	}()
}

func (ws *Websocket) handleTriggerAPI(w http.ResponseWriter, r *http.Request) {

	var message Message
	body, err := ioutil.ReadAll(r.Body)
	exception.Log(err)
	fmt.Println("message")
	err = json.Unmarshal(body, &message)
	exception.Log(err)
	queryParams := r.URL.Query()
	auth_signature := queryParams.Get("auth_signature")
	delete(queryParams, "auth_signature")

	var keys []string
	for key := range queryParams {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var stringToSign strings.Builder
	stringToSign.WriteString(r.Method + "\n" + r.URL.Path + "\n")

	for i, key := range keys {
		if i > 0 {
			stringToSign.WriteString("&")
		}
		stringToSign.WriteString(key + "=" + queryParams.Get(key))
	}
	ou := model.Get([]string{"public_key", "secret_key"}, "channels", []model.Where_st{
		{Key: "public_key", Value: queryParams.Get("auth_key"), After: "", Operator: "="},
	}, "id", true)

	expectedSignature := generateHMACSHA256(ou[0]["secret_key"], stringToSign.String())

	if expectedSignature != auth_signature {
		return
	}

	ws.publish(queryParams.Get("auth_key"), message)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}

func (ws *Websocket) eventHandel(msg chan msgS) {
	select {
	case message := <-msg:
		switch message.message.Name {
		case "subscribe":
			ws.subscribe(message.public_key, message.conn, *message.message)
		case "unsubscribe":
			ws.unsubscribe(message.public_key, message.conn, *message.message)

		}

	}

}
