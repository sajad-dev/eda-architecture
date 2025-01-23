package main

import (
	"github.com/joho/godotenv"
	"github.com/sajad-dev/eda-architecture/internal/exception"
	"github.com/sajad-dev/eda-architecture/internal/pub_sub"
	"github.com/sajad-dev/eda-architecture/internal/websocket"
)

func main() {
	err := godotenv.Load(".env")
	exception.Log(err)
	adrr := []websocket.Addr{
		{Pattern: "/channel", Handler: pub_sub.HandlerFunc, MiddlewareList: []websocket.MiddlewareFuncType{}},
	}
	_, waitGroup := websocket.Handler(adrr)
	waitGroup.Wait()
}
