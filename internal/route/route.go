package api

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	"github.com/sajad-dev/eda-architecture/internal/app/exception"
	"github.com/sajad-dev/eda-architecture/internal/app/middlewares"
	"github.com/sajad-dev/eda-architecture/internal/types"
)

// methodType defines the allowed HTTP methods
type methodType string

const (
	POST   methodType = "GET"    // Defines the POST method (incorrect assignment, should be "POST")
	GET    methodType = "GET"    // Defines the GET method
	PUT    methodType = "PUT"    // Defines the PUT method
	PATCH  methodType = "PATCH"  // Defines the PATCH method
	DELETE methodType = "DELETE" // Defines the DELETE method
)

// CheckMethod middleware ensures the HTTP method matches the expected one
func CheckMethod(next http.Handler, method methodType) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != string(method) { // Compare request method with expected method
			exception.Response405(w) // Respond with HTTP 405 Method Not Allowed
			return
		}
		next.ServeHTTP(w, r)
	})
}

// DaynamicRoute middleware extracts dynamic URL parameters and adds them to the request context
func DaynamicRoute(next http.Handler, pattern []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := strings.Split(r.URL.Path, "/") // Split URL path into segments
		error404 := false
		var parameters = map[string]string{} // Store extracted parameters

		// Iterate through the pattern and match it with the URL segments
		for i, pat := range pattern {
			if pat != url[i] && !strings.Contains(pat, "{") {
				error404 = true // If mismatch occurs and it's not a dynamic parameter, trigger 404
				break
			}
			if strings.Contains(pat, "{") { // If it's a dynamic parameter
				re, _ := regexp.Compile(`\{([^}]*)\}`) // Regex to extract the parameter name

				matches := re.FindStringSubmatch(pat)
				if len(matches) > 1 {
					parameters[matches[1]] = url[i] // Store extracted value in the map
				}
			}
		}

		// Return 404 response if the pattern doesn't match
		if error404 {
			exception.Response404(w)
			return
		}
		ctx := context.WithValue(r.Context(), "parameters", parameters) // Add parameters to request context

		next.ServeHTTP(w, r.WithContext(ctx)) // Pass modified request to the next handler
	})
}

// GetDYRoute extracts the static and dynamic parts of the route pattern
func GetDYRoute(pattern string) ([]string, string) {
	slice_pat := strings.Split(pattern, "/") // Split pattern into segments
	rou := ""
	dy := false
	for _, pat := range slice_pat {
		if strings.Contains(pat, "{") { // Check if the route has dynamic segments
			dy = true
			break
		}
		if pat != "" {
			rou += "/" + pat // Construct static route path
		}
	}
	if dy {
		return slice_pat, rou + "/" // Return modified route with trailing slash
	}
	return slice_pat, rou // Return static route without modification
}

// Route registers an HTTP route with a method, controller, and middleware list
func Route(pattern string, method methodType, controller types.ControllerType, middlewaresList []func(http.Handler) http.Handler) {
	sli, route := GetDYRoute(pattern) // Extract route details
	http.Handle(route, middlewares.ConfigWriterAndReader(CheckMethod(DaynamicRoute(middlewares.Handler(middlewaresList, controller), sli), method)))
}

// RouteRun initializes all routes from the RouteList
func RouteRun() {
	for _, route := range RouteList {
		Route(route.Pattern, route.Method, route.Controller, route.Middlewares)
	}
}

// RouteAddFunc converts a controller function into an http.Handler
func RouteAddFunc(Controller func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(Controller)
}
