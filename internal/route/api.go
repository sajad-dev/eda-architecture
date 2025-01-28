package api

import (
	"net/http"

	"github.com/sajad-dev/eda-architecture/internal/app/controllers"
	publictypes "github.com/sajad-dev/eda-architecture/internal/types"
)

type MiddlewaresListType []func(http.Handler) http.Handler
type ApiType struct {
	Pattern     string
	Method      methodType
	Controller  publictypes.ControllerType
	Middlewares MiddlewaresListType
}

var RouteList = []ApiType{
	{
		Pattern:     "/api/create-key",
		Method:      "POST",
		Controller:  controllers.CreateKey,
		Middlewares: MiddlewaresListType{},
	},
}
