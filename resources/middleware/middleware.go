package middleware

import (
	"fmt"

	redis "github.com/Pholey/bitAPI/redis"
	pass "github.com/Pholey/bitAPI/resources/lib"
	"github.com/gin-gonic/gin"
)

func UserRequired(inner func(*gin.Context)) func(*gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		token, parseErr := pass.ParseBasic(authHeader)
		fmt.Print("we got teh token!", token)
		if parseErr != nil {
			c.JSON(401, gin.H{"error": "Malformed authorization header"})
		} else {

			// Look our token up in redis
			username := redis.Client.Get(token).Val()

			if username == "" {
				c.JSON(404, gin.H{"error": "Unauthorized"})
			} else {
				// Continue on as normal... NOTE: Should probably do some context
				// magic with the user to pass it along the requests
				inner(c)
			}
		}
	}
}
