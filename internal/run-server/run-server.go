package runserver

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sajad-dev/eda-architecture/internal/app/exception"
	"github.com/sajad-dev/eda-architecture/internal/app/helpers"
)

func Run() {
	fmt.Println("Server is running on port 8000")
	err := http.ListenAndServe(":8000", nil)
	exception.Log(err)
	if !helpers.IfThenElse(os.Getenv("DEBUG") == "true", true, false).(bool) {
		// defer log.Panicln("END PROGRAM")
	}
}
