package controllers

import (
	"fmt"
	"net/http"

	"github.com/sajad-dev/eda-architecture/internal/app/utils"
	"github.com/sajad-dev/eda-architecture/internal/app/websocket"
	"github.com/sajad-dev/eda-architecture/internal/database/model"
)

// CreateKey generates a new public and secret key, stores them in the database, 
// and dynamically adds a new WebSocket route for the created key.
func CreateKey(w http.ResponseWriter, r *http.Request) {
	// Generate a random public key
	public := utils.GenerateRandomString(16)

	// Generate a random secret key
	secret := utils.GenerateRandomString(16)

	// Insert the generated keys into the "channels" table
	model.Insert(map[string]string{
		"public_key": public,
		"secret_key": secret,
	}, "channels")

	// If the WebSocket server is active, add a new route for the generated public key
	if websocket.ActiveSocket != nil {
		websocket.ActiveSocket.AddAddr("/app/" + public)
	}

	// Respond with the generated keys in JSON format
	w.Write([]byte(fmt.Sprintf(`{
	"public_key": "%s",
	"secret_key": "%s"
	}`, public, secret)))
}
