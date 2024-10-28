package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", fmt.Errorf("failed generated hash password")
	}
	return string(hash), nil
}

func ComparePassword(hash string, plain []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), plain)
	return err == nil
}
