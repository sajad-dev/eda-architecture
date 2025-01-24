package main

import (
	"github.com/joho/godotenv"
	"github.com/sajad-dev/eda-architecture/pkg/exception"
	"github.com/sajad-dev/eda-architecture/pkg/websocket"
)

func main() {
	err := godotenv.Load(".env")
	exception.Log(err)
	adrr := []websocket.Addr{
		{Pattern: "/channel", Handler: websocket.HandlerFunc, MiddlewareList: []websocket.MiddlewareFuncType{}},
	}
	websocket.Handler(adrr)
}
