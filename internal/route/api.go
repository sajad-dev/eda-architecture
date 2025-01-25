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
		Pattern:     "/add-channel",
		Method:      "POST",
		Controller:  controllers.AddChannel,
		Middlewares: MiddlewaresListType{},
	},
	{
		Pattern:     "/users",
		Method:      "POST",
		Controller:  controllers.AddChannel,
		Middlewares: MiddlewaresListType{},
	},
		{
		Pattern:     "/users",
		Method:      "GET",
		Controller:  controllers.Users,
		Middlewares: MiddlewaresListType{},
	},
		{
		Pattern:     "/users",
		Method:      "DELETE",
		Controller:  controllers.DeleteUser,
		Middlewares: MiddlewaresListType{},
	},
}
