package main

import (
	"os"

	"github.com/sajad-dev/eda-architecture/internal/command"
	connectiondb "github.com/sajad-dev/eda-architecture/internal/connection_db"
	"github.com/sajad-dev/eda-architecture/internal/api"
	"github.com/sajad-dev/eda-architecture/internal/migration"
	runserver "github.com/sajad-dev/eda-architecture/internal/run-server"
	"github.com/sajad-dev/eda-architecture/internal/websocket"
)

// func main() {
// 	err := godotenv.Load(".env")
// 	exception.Log(err)
// 	adrr := []websocket.Addr{}
// 	websocket.Handler(adrr)
// }

func main() {

	runserver.Init()

	connectiondb.Connection()
	if len(os.Args) > 2 {
		command.Handel(os.Args)
		return
	}
	websocket.Handler([]websocket.Addr{})
	migration.Handel()
	api.RouteRun()

	runserver.Run()
}
