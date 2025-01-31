package middlewares

import (
	"net/http"

	"github.com/sajad-dev/eda-architecture/internal/types"
)

// HandelMiddleware wraps an HTTP handler with a middleware function
func (e MiddlewaresType) HandelMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		e(w, r, next) // Execute the middleware function
	})
}

// finallyHandler wraps a controller function in an HTTP handler
func finallyHandler(controller types.ControllerType) http.Handler {
	return http.HandlerFunc(controller)
}

// Handler applies a list of middleware functions to a controller and returns the final HTTP handler
func Handler(middlewares []func(http.Handler) http.Handler, finally types.ControllerType) http.Handler {
	finally_co := finallyHandler(finally) // Convert controller to an HTTP handler
	
	// Apply middlewares in reverse order (last middleware wraps the previous ones)
	for i := len(middlewares) - 1; i >= 0; i-- {
		finally_co = middlewares[i](finally_co)
	}
	return finally_co
}

// ConfigWriterAndReader adds default headers for JSON responses and CORS support
func ConfigWriterAndReader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
