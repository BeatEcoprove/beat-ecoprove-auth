package services

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

const SALT_COST = 11

func GenerateSalt(length int) (string, error) {
	salt := make([]byte, length)

	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

func HashPassword(password string, salt string) (string, error) {
	saltedPassword := password + salt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func CheckPasswordHash(password string, salt, hash string) bool {
	saltedPassword := password + salt
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(saltedPassword))

	return err == nil
}

func GeneratePassword(minPasswordLength, maxPasswordLength int64) (string, error) {
	const passwordChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}|;:,.<>?/~`"

	length, err := rand.Int(rand.Reader, big.NewInt(maxPasswordLength-minPasswordLength+1))
	if err != nil {
		return "", err
	}

	length.Add(length, big.NewInt(minPasswordLength))

	password := make([]byte, length.Int64())

	for i := range password {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(passwordChars))))
		if err != nil {
			return "", err
		}
		password[i] = passwordChars[index.Int64()]
	}

	return string(password), nil
}
