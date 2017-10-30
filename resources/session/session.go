package session

import (
	"math/rand"

	redis "github.com/Pholey/bitAPI/redis"
	valid "github.com/Pholey/bitAPI/lib/validator"
	user "github.com/Pholey/bitAPI/resources/user"
	"github.com/labstack/echo"
	"net/http"
)

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
	var req user.User
	if err := c.Bind(&req); err != nil {
		return err
	}

	err := valid.Validate.Struct(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	isAuthed, user, err := user.Auth(req.UserName, req.Password)

	var key string
	// TODO: This is horrendous
	if err != nil {
		return err
	} else if isAuthed == false {
		return echo.NewHTTPError(http.StatusUnauthorized)
	} else {
		// Create a new token (50 characters)
		key = randSeq(50)
		err := redis.Client.Set(key, user.Id, 0).Err()
		if err != nil {
			return err
		}

		c.JSON(200, Token{key})
	}

	return nil;
}
