package middlewares

import (
	"net/http"

	"github.com/sajad-dev/eda-architecture/internal/types"
)

func (e MiddlewaresType) HandelMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		e(w, r, next)
	})
}

func finallyHandler(controller types.ControllerType) http.Handler {
	return http.HandlerFunc(controller)
}

func Handler(middlewares []func(http.Handler) http.Handler, finally types.ControllerType) http.Handler {
	finally_co := finallyHandler(finally)
	for i := len(middlewares) - 1; i >= 0; i-- {
		finally_co = middlewares[i](finally_co)
	}
	return finally_co
}

func ConfigWriterAndReader(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		next.ServeHTTP(w, r)
	})
}

