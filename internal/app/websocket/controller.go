package websocket

import (
	"net/http"

	// "encoding/json"
	"github.com/sajad-dev/eda-architecture/internal/app/exception"
)

func (ws *Websocket) unSubscribe() {

}

func (ws *Websocket) publish() {

}

func (ws *Websocket) subscribe() {

}

func (ws *Websocket) HandlerFunc(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	exception.Log(err)
	// secret := r.Header.Get("secret_key")

	go func() {
		defer conn.Close()
		for {
			var message Message
			err := conn.ReadJSON(&message)
			exception.Log(err)

			switch message.Event {
			case "subscribe":
				continue
			case "unsubscribe":
				continue
			case "publish":
				continue
			}

		}
	}()
}
