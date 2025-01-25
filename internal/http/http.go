package http

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sajad-dev/eda-architecture/internal/exception"
	"github.com/sajad-dev/eda-architecture/internal/websocket"
)

func WebServer(ws *websocket.Websocket) {

	go func() {

		// http.HandleFunc("/add-channel", httpWebHandler(ws, AddSocketChannel))
		err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil)
		exception.Log(err)
	}()
}

func httpWebHandler(handlerFunc func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlerFunc(w, r)
	}
}
