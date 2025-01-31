package exception

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/fatih/color"
	"github.com/sajad-dev/eda-architecture/internal/app/helpers"
	"github.com/sajad-dev/eda-architecture/internal/app/response"
)

// Response500 sends an HTTP 500 Internal Server Error response with detailed logs
func Response500(w http.ResponseWriter, exception string) {
	_, file, line, ok := runtime.Caller(1) // Retrieve caller file and line number

	// If debugging is enabled, send a detailed error response
	if helpers.IfThenElse(os.Getenv("DEBUG") == "true", true, false).(bool) {
		res := fmt.Sprintf("Error occurred in %s:%d - %s", file, line, exception)
		json.NewEncoder(w).Encode(response.ErrorResponse{Messages: res, Code: 500, Status: false})
		return
	}

	// Log error details if available
	if ok {
		log.Printf("Error occurred in %s:%d - %s", file, line, exception)
	} else {
		log.Println("Error occurred:", exception)
	}

	// Send a generic error message to the client
	json.NewEncoder(w).Encode(response.ErrorResponse{Messages: "Internal Server Error", Code: 500, Status: false})
}

// Response405 sends an HTTP 405 Method Not Allowed response
func Response405(w http.ResponseWriter) {
	_, file, line, ok := runtime.Caller(1) // Retrieve caller file and line number

	// If debugging is enabled, send a detailed error response
	if helpers.IfThenElse(os.Getenv("DEBUG") == "true", true, false).(bool) {
		res := fmt.Sprintf("Error occurred in %s:%d - %s", file, line, "Method Not Allowed")
		json.NewEncoder(w).Encode(response.ErrorResponse{Messages: res, Code: 500, Status: false})
		return
	}

	// Log error details if available
	if ok {
		log.Printf("Error occurred in %s:%d - %s", file, line, "Method Not Allowed")
	} else {
		log.Println("Error occurred:", "Method Not Allowed")
	}

	// Send HTTP 405 response
	w.WriteHeader(http.StatusMethodNotAllowed)
	json.NewEncoder(w).Encode(response.ErrorResponse{Messages: "Method Not Allowed", Code: 405, Status: false})
}

// Response404 sends an HTTP 404 Not Found response
func Response404(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(response.ErrorResponse{Messages: "Not Found", Code: 404, Status: false})
}

// Log logs an error message, with optional debugging details
func Log(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1) // Retrieve caller file and line number

		// Log the error differently based on debugging mode
		if !helpers.IfThenElse(os.Getenv("DEBUG") == "true", true, false).(bool) {
			log.Println(err) // Log the error normally
		} else {
			erro := fmt.Sprintf("%s - line %d - file %s", err.Error(), line, file)
			color.Red(erro) // Print error message in red
		}
	}
}
