package user

import "go-oj/pkg/database"

func IsPhoneExist(phone string) bool {
	return database.IsExist("user_basic", "phone", phone)
}

func IsEmailExist(email string) bool {
	return database.IsExist("user_basic", "mail", email)
}
