package main

import (
	"github.com/joho/godotenv"
	"github.com/sajad-dev/eda-architecture/internal/exception"
	"github.com/sajad-dev/eda-architecture/internal/websocket"
)

func main() {
	err := godotenv.Load(".env")
	exception.Log(err)
	adrr := []websocket.Addr{}
	websocket.Handler(adrr)
}
