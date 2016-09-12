package session

import (
	json "encoding/json"
	"io"
	"io/ioutil"
	"math/rand"

	redis "github.com/Pholey/bitAPI/redis"
	pass "github.com/Pholey/bitAPI/resources/lib"

	"github.com/labstack/echo"
)

type Session struct {
	Password string `json:"password"`
	UserName string `json:"userName"`
}

type Error struct {
	Error string `json: "error"`
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

	// TODO(pholey): Abstract this out or find a better lib
	body, err := ioutil.ReadAll(io.LimitReader(c.Request().Body(), 1048576))

	if err != nil {
		// TODO(pholey): Proper error handing
		panic(err)
	}

	// Grab our session data
	var req Session
	if err := json.Unmarshal(body, &req); err != nil {
		panic(err)
	}

	isAuthed, err := pass.Auth(req.UserName, req.Password)

	var key string
	// TODO: This is horrendous
	if err != nil {
		panic(err)
	} else if isAuthed == false {
		c.JSON(404, Error{"User not found"})
	} else {
		// Create a new token (50 characters)
		key = randSeq(50)
		err := redis.Client.Set(key, req.UserName, 0).Err()
		if err != nil {
			panic(err)
		}

		c.JSON(200, Token{key})
	}

	return nil;
}
