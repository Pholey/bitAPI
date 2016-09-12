package logger

import (
	"log"
	"time"

	"github.com/labstack/echo"
)

// Logger - Logs info about incoming requests
func Logger(inner echo.HandlerFunc, name string) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()


		log.Printf(
			"%s\t%s\t%s\t%s",
			c.Request().Method(),
			c.Request().URI(),
			name,
			time.Since(start),
		)
		return inner(c)
	}
}
