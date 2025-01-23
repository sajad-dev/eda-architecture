package websocket

import (
	"net/http"

	"github.com/sajad-dev/eda-architecture/internal/exception"
)

func UpgraderMiddleware(http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := Upgrader.Upgrade(w, r, nil)
		exception.Log(err)
		defer conn.Close()
	})
}
