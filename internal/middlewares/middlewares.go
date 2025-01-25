package middlewares

import (
	"net/http"

	publictypes "github.com/sajad-dev/eda-architecture/internal/public_types"
)

func (e MiddlewaresType) HandelMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		e(w, r, next)
	})
}

func finallyHandler(controller publictypes.ControllerType) http.Handler {
	return http.HandlerFunc(controller)
}

func Handler(middlewares []func(http.Handler) http.Handler, finally publictypes.ControllerType) http.Handler {
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

func DaynamicRoute(w http.ResponseWriter, r *http.Request, next http.Handler) {

}
