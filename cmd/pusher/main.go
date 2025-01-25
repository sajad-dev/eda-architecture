package main

import (
	"os"

	"github.com/fatih/color"
	"github.com/sajad-dev/eda-architecture/internal/app/websocket"
	"github.com/sajad-dev/eda-architecture/internal/command"
	connectiondb "github.com/sajad-dev/eda-architecture/internal/database/connection_db"
	"github.com/sajad-dev/eda-architecture/internal/database/migration"
	api "github.com/sajad-dev/eda-architecture/internal/route"
	runserver "github.com/sajad-dev/eda-architecture/internal/run-server"
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
	color.Blue("Adress Local : http://127.0.0.1:8000")

	go func() {
		websocket.Handler([]websocket.Addr{})
	}()
	migration.Handel()
	api.RouteRun()

	runserver.Run()
}
