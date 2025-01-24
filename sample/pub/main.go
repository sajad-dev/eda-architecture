// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/sajad-dev/eda-architecture/pkg/exception"
)

func data(upd chan int) string {
	return "Hi"
}

func main() {

	jsonStr := []byte(`{name: 'new-channel'}`)
	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/add-channel", bytes.NewBuffer(jsonStr))
	exception.Log(err)

	client := &http.Client{}
	_, err = client.Do(req)
	exception.Log(err)

	u := url.URL{Scheme: "ws", Host: ":8081", Path: "/new-channel"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	exception.Log(err)

	defer c.Close()

	done := make(chan struct{})
	upd := make(chan int)

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				exception.Log(err)
				break
			}

			if "new-subscriber" == string(message) {
				err := c.WriteMessage(websocket.TextMessage, []byte(data(upd)))
				exception.Log(err)
			}

			select {
			case <-upd:
				err := c.WriteMessage(websocket.TextMessage, []byte(data(upd)))
				exception.Log(err)
			case <-done:
				break
			}
		}
	}()

	for {
		select {
		case <-done:
			return
		}
	}
}
