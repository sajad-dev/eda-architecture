package runserver

import (
	"net/http"
	"os"

	"github.com/sajad-dev/eda-architecture/internal/app/exception"
	"github.com/sajad-dev/eda-architecture/internal/app/helpers"
)

func Run() {
	err := http.ListenAndServe(":8000", nil)
	exception.Log(err)
	if !helpers.IfThenElse(os.Getenv("DEBUG") == "true", true, false).(bool) {
		// defer log.Panicln("END PROGRAM")
	}
}
