package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/sajad-dev/eda-architecture/internal/app/websocket"
	"github.com/sajad-dev/eda-architecture/internal/command"
	connectiondb "github.com/sajad-dev/eda-architecture/internal/database/connection_db"
	"github.com/sajad-dev/eda-architecture/internal/database/migration"
	api "github.com/sajad-dev/eda-architecture/internal/route"
	runserver "github.com/sajad-dev/eda-architecture/internal/run_server"
)

func main() {

	runserver.Init()

	connectiondb.Connection()
	if len(os.Args) > 2 {
		command.Handel(os.Args)
		return
	}
	color.Blue(fmt.Sprintf("Address Local : http://127.0.0.1:%s", os.Getenv("PORT")))

	websocket.Handler()
	migration.Handel()
	api.RouteRun()

	runserver.Run()
}
