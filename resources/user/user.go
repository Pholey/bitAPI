package user

import (
	http "net/http"

	"database/sql"
	db "github.com/Pholey/bitAPI/db"
	pass "github.com/Pholey/bitAPI/lib/password"
	valid "github.com/Pholey/bitAPI/lib/validator"

	"github.com/labstack/echo"
)

// User - Struct for dealing with users
type User struct {
	Id         int64  `db:"id"`
	Salt       []byte `db:"password_salt"`
	Iterations int    `db:"password_iterations"`
	Hash       []byte `db:"password_hash"`
	Password   string `json:"password" validate:"required"`
	UserName   string `db:"username" json:"username" validate:"required"`
	Email      string `db:"email"    json:"email"`
}

// Users - Our array of users
type Users []User

// Create - Create a user
func Create(c echo.Context) error {
	// Grab our user data
	var user User
	if err:= c.Bind(&user); err != nil {
		return err
	}

	err := valid.Validate.Struct(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	// TODO(pholey): Password validation, refactoring
	// Hash our password
	salt, iterations, hash, err := pass.HashPass(user.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Create our statement
	query := db.Session.
		InsertInto("user").
		Columns("username", "email", "password_hash",
			"password_iterations", "password_salt").
		Values(user.UserName, user.Email, hash, iterations, salt)

	rows, err := query.Exec()

	_ = (rows)

	if err != nil {
		return err
	}

	c.Response().Header().Set("Content-Type", "application/json;charset=UTF-8")
	// Since we were sent JSON, we should be nice and return an empty JSON object
	c.JSON(http.StatusCreated, struct{}{})

	return nil;
}

func GetUser(username string) (*User, error) {
	var user User

	query := db.Session.
		Select("*").
		From("\"user\"").
		Where("username = ?", username)

	if _, err := query.Load(&user); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func Auth(username string, password string) (bool, *User, error) {
	user, err := GetUser(username)

	if err != nil {
		panic(err)
	}

	// This code smell is getting old really fast
	if user == nil {
		return false, nil, nil
	} else if err != nil {
		return false, nil, err
	}

	return pass.VerifyHash(password, user.Salt, user.Iterations, user.Hash), user, nil
}
