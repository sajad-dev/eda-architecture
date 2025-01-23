package websocket

import "net/http"

func UpgraderMiddleware(http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Upgrader.Upgrade(w, r, nil)
	})
}
