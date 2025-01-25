package api

import (
	"net/http"

	"github.com/sajad-dev/eda-architecture/internal/controllers"
	publictypes "github.com/sajad-dev/eda-architecture/internal/public_types"
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
		Pattern:     "/add-channel",
		Method:      "POST",
		Controller:  controllers.AddChannel,
		Middlewares: MiddlewaresListType{},
	},
}
