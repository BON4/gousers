package utils

import (
	bc "golang.org/x/crypto/bcrypt"
)

func VerifyPassword(password []byte, hash []byte) (bool, error) {
	if err := bc.CompareHashAndPassword(hash, password); err != nil {
		return false, err
	}
	return true, nil
}

func HashPassword(password []byte) (*[]byte, error) {
	hashedpass, err := bc.GenerateFromPassword(password, 10)
	if err != nil {
		return nil, err
	}
	return &hashedpass, nil
}
