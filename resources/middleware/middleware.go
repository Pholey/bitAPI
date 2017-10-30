package middleware

import (
	"errors"
	"strings"

	redis "github.com/Pholey/bitAPI/redis"
	"github.com/labstack/echo"

	"net/http"
)

func parseBasic(authHeader string) (string, error) {
	err := errors.New("Authorization header required")

	if len(authHeader) == 0 {
		return "", nil
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		return "", err
	}

	kind := strings.ToLower(parts[0])
	credentials := parts[1]

	if kind != "bearer" {
		return "", err
	}

	return credentials, nil
}

type ContextWithUser struct {
	echo.Context
	UserId    int64
}

func userRequired(token string, inner echo.HandlerFunc, c echo.Context) error {
	// Look our token up in redis
	userId, err := redis.Client.Get(token).Int64()

	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	} else {
		// Continue on as normal...

		cu := &ContextWithUser{c, userId}

		cu.UserId = userId
		return inner(cu)
	}
}

func UserRequiredBearer(inner echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header["Authorization"][0]
		token, parseErr := parseBasic(authHeader)

		if parseErr != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		} else {
			return userRequired(token, inner, c)
		}
	}
}

func UserRequiredToken(inner echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if len(c.Request().Header["Sec-Websocket-Protocol"]) == 0 {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		token := c.Request().Header["Sec-Websocket-Protocol"][0]

		return userRequired(token, inner, c)
	}
}
