package hash

import (
	"go-oj/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

func Encrypt(password string) string {
	b, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		logger.LogIf(err)
		return ""
	}
	return string(b)
}

func CheckPassword(hashPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	logger.LogIf(err)
	return err == nil
}
