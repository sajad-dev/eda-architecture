package exception

import (
	"fmt"
	// "log"
	"runtime"

	"github.com/fatih/color"
)

func Log(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		erro := fmt.Sprintf("%s - line %d - file %s", err.Error(), line, file)
		color.Red(erro)
		// log.Panicln(erro)

	}
}
