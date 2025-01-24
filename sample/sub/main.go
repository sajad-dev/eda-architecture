// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/sajad-dev/eda-architecture/pkg/exception"
)

func data(msg []byte) {
	fmt.Println(string(msg))
}

func main() {

	u := url.URL{Scheme: "ws", Host: ":8081", Path: "/new-channel"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	exception.Log(err)

	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				exception.Log(err)
				break
			}

			data(message)
			select {
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
