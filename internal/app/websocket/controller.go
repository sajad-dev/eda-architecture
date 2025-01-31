package websocket

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/sajad-dev/eda-architecture/internal/app/exception"
	"github.com/sajad-dev/eda-architecture/internal/app/response"
)

// HandlerFunc upgrades an HTTP request to a WebSocket connection
// and listens for incoming messages.
func (ws *Websocket) HandlerFunc(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil) // Upgrade HTTP to WebSocket
	exception.Log(err) // Log any errors

	ws.CountConn++ // Increment the connection counter

	// Extract the public key from the URL using regex
	regex := regexp.MustCompile(`/app/([^/?]+)`)
	match := regex.FindStringSubmatch(r.URL.Path)

	// Send an initial connection message to the client
	conn.WriteJSON(Message{Event: "pusher:connection_established", Data: map[string]interface{}{
		"socket_id":        ws.CountConn,
		"activity_timeout": "120",
	}})

	var messageChan = make(chan MessageChan, 1) // Channel for handling messages

	go func() {
		for {
			var message Message
			message.Data = make(map[string]string) // Initialize message data map

			_, ou, err := conn.ReadMessage() // Read incoming message from WebSocket
			if err != nil {
				exception.Log(err) // Log the error if reading fails
				break
			}

			err = json.Unmarshal(ou, &message) // Parse JSON message
			exception.Log(err)

			// Send the received message to the handler
			messageChan <- MessageChan{connection: conn, public_key: match[1], message: &message}

			go func() {
				ws.eventHandel(messageChan) // Handle the received event asynchronously
			}()
		}
	}()
}

// handleTriggerAPI processes API requests to trigger WebSocket messages
func (ws *Websocket) handleTriggerAPI(w http.ResponseWriter, r *http.Request) {
	var message TriggerBody

	err := json.NewDecoder(r.Body).Decode(&message) // Decode incoming JSON request
	exception.Log(err)

	// If "Channel" is set, append it to the "Channels" list
	if len(message.Channel) != 0 {
		message.Channels = append(message.Channels, message.Channel)
	}

	queryParams := r.URL.Query() // Retrieve query parameters

	// Validate the private key before proceeding
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

	// Publish the message using the authentication key
	ws.publish(queryParams.Get("auth_key"), message)

	// Send a success response back to the client
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.SuccessResponse{
		Message: "Success ",
		Status:  "200 OK",
		Data:    message,
	})
}
