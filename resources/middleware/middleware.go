package middleware

import (
	"fmt"

	redis "github.com/Pholey/bitAPI/redis"
	pass "github.com/Pholey/bitAPI/resources/lib"
	"github.com/labstack/echo"

	"net/http"
)

func UserRequired(inner echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header().Get("Authorization")
		token, parseErr := pass.ParseBasic(authHeader)
		fmt.Print("we got teh token!", token)
		if parseErr != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		} else {

			// Look our token up in redis
			username := redis.Client.Get(token).Val()

			if username == "" {
				return echo.NewHTTPError(http.StatusUnauthorized)
			} else {
				// Continue on as normal... NOTE: Should probably do some context
				// magic with the user to pass it along the requests
				return inner(c)
			}
		}
	}
}
