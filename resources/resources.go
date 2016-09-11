package resources

import (
	L "github.com/Pholey/bitAPI/logger"
	routes "github.com/Pholey/bitAPI/resources/routes"
	"github.com/gin-gonic/gin"
)

// NewRouter - Returns a router with iniialized routes
func NewRouter() *gin.Engine {
	router := gin.Default()

	handlerMap := map[string]func(string, ...gin.HandlerFunc) gin.IRoutes{
		"GET":     router.GET,
		"POST":    router.POST,
		"PUT":     router.PUT,
		"DELETE":  router.DELETE,
		"PATCH":   router.PATCH,
		"HEAD":    router.HEAD,
		"OPTIONS": router.OPTIONS,
	}

	for _, route := range routes.Routes {
		// Set up logging for each request
		handler := L.Logger(route.HandlerFunc, route.Name)

		for _, beforeFunc := range route.Before {
			handler = beforeFunc(handler)
		}

		handlerMap[route.Method](route.Pattern, handler)
	}

	return router
}
