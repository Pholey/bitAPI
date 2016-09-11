package lib

import (
	"bytes"
	rand "crypto/rand"
	sha256 "crypto/sha256"
	"database/sql"
	"errors"
	big "math/big"
	"strings"

	"github.com/Pholey/Exgo/db"
	pbkdf2 "golang.org/x/crypto/pbkdf2"
)

type User struct {
	Salt       []byte `db:"password_salt"`
	Iterations int    `db:"password_iterations"`
	Hash       []byte `db:"password_hash"`
	Password   string `json:"password"`
	UserName   string `db:"username" json:"userName"`
	Name       string `db:"name"     json:"name"`
	Email      string `db:"email"    json:"email"`
}

// FIXME: If the rand stuff in here fails, it
// will probably crash the app, there is no
// error handling here
func HashPass(password string) ([]byte, int, []byte) {
	salt := make([]byte, 32)
	rand.Read(salt)
	ii, _ := rand.Int(rand.Reader, big.NewInt(16000))
	iterations := int(ii.Int64()) + 64000
	hash := pbkdf2.Key([]byte(password), salt, iterations, 32, sha256.New)
	return salt, iterations, hash
}

func VerifyHash(password string, salt []byte, iterations int, hash []byte) bool {
	compareHash := pbkdf2.Key([]byte(password), salt, iterations, 32, sha256.New)
	return bytes.Equal(hash, compareHash)
}

func ParseBasic(authHeader string) (string, error) {
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

func GetUserAuthInfo(username string) (*User, error) {
	var user User

	query, args, err := db.Sq.
		Select("password_salt", "password_iterations", "password_hash").
		From("\"user\"").
		Where("username = $1", username).
		ToSql()

	if err != nil {
		return nil, err
	}

	if err = db.Client.Get(&user, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func Auth(username string, password string) (bool, error) {
	user, err := GetUserAuthInfo(username)

	if err != nil {
		panic(err)
	}

	// This code smell is getting old really fast
	if user == nil {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return VerifyHash(password, user.Salt, user.Iterations, user.Hash), nil
}
