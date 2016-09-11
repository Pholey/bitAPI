package routes

import (
	mw "github.com/Pholey/bitAPI/resources/middleware"
	session "github.com/Pholey/bitAPI/resources/session"
	socket "github.com/Pholey/bitAPI/resources/socket"
	user "github.com/Pholey/bitAPI/resources/user"
	"github.com/gin-gonic/gin"
)

// Decorator middleware
type beforeHandlers []func(func(*gin.Context)) func(*gin.Context)

// Route - Struct containing all the info to initialize a route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc func(*gin.Context)
	Before      beforeHandlers
}

// Routes - Array of routes
var Routes = []Route{
	Route{
		"CreateUser",
		"POST",
		"/user",
		user.Create,
		beforeHandlers{},
	},
	Route{
		"CreateSession",
		"POST",
		"/session",
		session.Create,
		beforeHandlers{},
	},
	Route{
		"Upgrade Socket",
		"GET",
		"/socket",
		socket.Entry,
		beforeHandlers{mw.UserRequired},
	},
}
