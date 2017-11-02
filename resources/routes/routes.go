package routes

import (
	mw "github.com/Pholey/bitAPI/resources/middleware"
	session "github.com/Pholey/bitAPI/resources/session"
	socket "github.com/Pholey/bitAPI/resources/socket"
	user "github.com/Pholey/bitAPI/resources/user"
	"github.com/labstack/echo"
)

// Decorator middleware
type beforeHandlers []func(echo.HandlerFunc) echo.HandlerFunc

// Route - Struct containing all the info to initialize a route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc echo.HandlerFunc
	Before      beforeHandlers
}

var sock, sock_err = socket.NewSocket()

// Routes - Array of routes
var Routes = []Route{
	Route{
		"CreateUser",
		"POST",
		"/user",
		echo.HandlerFunc(user.Create),
		beforeHandlers{},
	},
	Route{
		"CreateSession",
		"POST",
		"/session",
		echo.HandlerFunc(session.Create),
		beforeHandlers{},
	},
	Route{
		"UpgradeSocket",
		"GET",
		"/channel/:id",
		sock.Listen,
		beforeHandlers{mw.UserRequiredToken},
	},
	Route {
		"GetChannels",
		"GET",
		"/channels",
		sock.GetChannels,
		beforeHandlers{},
	},
	Route{
		"MessageChannel",
		"POST",
		"/channel/:id",
		sock.HTTPForward,
		beforeHandlers{},
	},
}
