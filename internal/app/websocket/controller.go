package websocket

import (
	"encoding/json"

	"net/http"
	"regexp"


	"github.com/sajad-dev/eda-architecture/internal/app/exception"
	"github.com/sajad-dev/eda-architecture/internal/app/response"
)

func (ws *Websocket) HandlerFunc(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	exception.Log(err)

	ws.CountConn++

	regex := regexp.MustCompile(`/app/([^/?]+)`)
	match := regex.FindStringSubmatch(r.URL.Path)

	conn.WriteJSON(Message{Event: "pusher:connection_established", Data: map[string]interface{}{
		"socket_id":        ws.CountConn,
		"activity_timeout": "120",
	}})

	var messageChan = make(chan MessageChan, 1)

	go func() {
		for {
			var message Message

			message.Data = make(map[string]string)

			_, ou, err := conn.ReadMessage()
			if err != nil {
				exception.Log(err)
				break
			}

			err = json.Unmarshal(ou, &message)
			exception.Log(err)

			messageChan <- MessageChan{connection: conn, public_key: match[1], message: &message}

			go func() {
				ws.eventHandel(messageChan)
			}()

		}
	}()
}

func (ws *Websocket) handleTriggerAPI(w http.ResponseWriter, r *http.Request) {

	var message TriggerBody

	err := json.NewDecoder(r.Body).Decode(&message)
	exception.Log(err)

	queryParams := r.URL.Query()

	if !checkPrivateKey(r.URL.Query(), r.URL.Path, r.Method) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.ErrorResponse{
			Messages: map[string]string{
				"messages": "Parameter value not available (1)",
			},
			Code:   400,
			Status: false,
		})
		return
	}

	ws.publish(queryParams.Get("auth_key"), message)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.SuccessResponse{
		Message: "Success ",
		Status:  "200 OK",
		Data:    message,
	})
}
