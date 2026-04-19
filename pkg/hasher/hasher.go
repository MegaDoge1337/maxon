package hasher

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (hash string, err error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	hash = string(hashBytes)
	return hash, nil
}

func CompareHashAndPassword(hashedPassword string, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}
