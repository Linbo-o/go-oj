package auth

import (
	"errors"
	"go-oj/app/models/user"
)

// LoginByPhone 尝试登录指定用户，检查手机账号是否注册，返回用户结构
func LoginByPhone(phone string) (user.UserBasic, error) {
	userModel := user.GetByPhone(phone)
	if userModel.Identity == "" {
		return user.UserBasic{}, errors.New("手机号未注册")
	}

	return userModel, nil
}
