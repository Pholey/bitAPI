package password

import (
	"bytes"
	rand "crypto/rand"
	sha256 "crypto/sha256"
	big "math/big"

	pbkdf2 "golang.org/x/crypto/pbkdf2"
)

// FIXME: If the rand stuff in here fails, it
// will probably crash the app, there is no
// error handling here
func HashPass(password string) ([]byte, int, []byte, error) {
	salt := make([]byte, 32)
	rand.Read(salt)
	ii, err := rand.Int(rand.Reader, big.NewInt(16000))

	if err != nil {
		return nil, 0, nil, err
	}

	iterations := int(ii.Int64()) + 4000
	hash := pbkdf2.Key([]byte(password), salt, iterations, 32, sha256.New)
	return salt, iterations, hash, nil
}

func VerifyHash(password string, salt []byte, iterations int, hash []byte) bool {
	compareHash := pbkdf2.Key([]byte(password), salt, iterations, 32, sha256.New)
	return bytes.Equal(hash, compareHash)
}
