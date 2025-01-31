package runserver

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sajad-dev/eda-architecture/internal/app/exception"
)

func Run() {
	err := http.ListenAndServe(fmt.Sprintf(":%s",os.Getenv("PORT")), nil)
	exception.Log(err)

}
