package middleware

import (
	"fmt"

	redis "github.com/Pholey/bitAPI/redis"
	pass "github.com/Pholey/bitAPI/resources/lib"
	"github.com/labstack/echo"
)

type Error struct {
	Error string `json: "error"`
}

func UserRequired(inner echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header().Get("Authorization")
		token, parseErr := pass.ParseBasic(authHeader)
		fmt.Print("we got teh token!", token)
		if parseErr != nil {
			c.JSON(401, Error{"Malformed authorization header"})
		} else {

			// Look our token up in redis
			username := redis.Client.Get(token).Val()

			if username == "" {
				c.JSON(404, Error{"Unauthorized"})
			} else {
				// Continue on as normal... NOTE: Should probably do some context
				// magic with the user to pass it along the requests
				return inner(c)
			}
		}
		return nil;
	}
}
