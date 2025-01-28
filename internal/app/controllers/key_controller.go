package controllers

import (
	"fmt"
	"net/http"

	"github.com/sajad-dev/eda-architecture/internal/app/utils"
	"github.com/sajad-dev/eda-architecture/internal/app/websocket"
	"github.com/sajad-dev/eda-architecture/internal/database/model"
)

func CreateKey(w http.ResponseWriter, r *http.Request) {
	public := utils.GenerateRandomString(16)
	secret := utils.GenerateRandomString(16)

	model.Insert(map[string]string{
		"public_key": public,
		"secret_key": secret,
	}, "channels")

	if websocket.ActiveSocket != nil {
		websocket.ActiveSocket.AddAddr("/"+public)
	}

	w.Write([]byte(fmt.Sprintf(`{
	"public_key": "%s",
	"secret_key": "%s"
	}`, public, secret)))
}
