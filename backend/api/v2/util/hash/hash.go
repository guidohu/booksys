package hash

import (
	"crypto/sha256"
	"fmt"

	crypt "github.com/tredoe/osutil/user/crypt/sha512_crypt"
)

// Sha256 returns a sha256 hashed version (in hex format) of the given string
func Sha256(p string) string {
	hasher := sha256.New()
	hasher.Write([]byte(p))
	hash := string(hasher.Sum(nil)[:])
	return fmt.Sprintf("%x", hash)
}

// CryptSha512 generates a sha512 hashed password representation that
// uses 5000 rounds and the given salt (same as PHP's crypt function).
func CryptSha512(p string, s string) (string, error) {
	c := crypt.New()
	saltStr := fmt.Sprintf("$6$rounds=5000$%s", s)
	salt := []byte(saltStr)
	return c.Generate([]byte(p), salt)
}