package session

import (
	"math/rand"

	redis "github.com/Pholey/bitAPI/redis"
	pass "github.com/Pholey/bitAPI/resources/lib"

	"github.com/labstack/echo"
	"net/http"
)

type Session struct {
	Password string `json:"password"`
	UserName string `json:"userName"`
}

type Token struct {
	Token string `json: "token"`
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func Create(c echo.Context) error {
	// Grab our session data
	var req Session
	if err := c.Bind(&req); err != nil {
		return err
	}

	isAuthed, err := pass.Auth(req.UserName, req.Password)

	var key string
	// TODO: This is horrendous
	if err != nil {
		return err
	} else if isAuthed == false {
		return echo.NewHTTPError(http.StatusNotFound)
	} else {
		// Create a new token (50 characters)
		key = randSeq(50)
		err := redis.Client.Set(key, req.UserName, 0).Err()
		if err != nil {
			return err
		}

		c.JSON(200, Token{key})
	}

	return nil;
}
