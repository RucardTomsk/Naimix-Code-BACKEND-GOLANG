package auth

import (
	"crypto/sha256"
	"fmt"
)

// Hasher uses SHA256 to hash passwords with provided salt.
type Hasher struct {
	salt string
}

func NewHasher(salt string) *Hasher {
	return &Hasher{salt: salt}
}

// Hash creates SHA256 hash of given password.
func (h *Hasher) Hash(password string) (string, error) {
	hash := sha256.New()

	if _, err := hash.Write([]byte(password)); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt))), nil
}
