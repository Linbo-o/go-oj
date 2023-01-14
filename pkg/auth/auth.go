package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-oj/app/models/user"
	"go-oj/pkg/logger"
)

// Attempt 尝试使用用户名或者手机或者邮箱登录
func Attempt(loginId string) (user.UserBasic, error) {
	userModel := user.GetByMulti(loginId)
	if userModel.Identity == "" {
		return user.UserBasic{}, errors.New("用户账号不正确或是用户密码不正确")
	}

	return userModel, nil
}

// LoginByPhone 尝试登录指定用户，检查手机账号是否注册，返回用户结构
func LoginByPhone(phone string) (user.UserBasic, error) {
	userModel := user.GetByPhone(phone)
	if userModel.Identity == "" {
		return user.UserBasic{}, errors.New("手机号未注册")
	}

	return userModel, nil
}

func CurrentUser(c *gin.Context) user.UserBasic {
	userModel, ok := c.MustGet("current_user").(user.UserBasic)
	if !ok {
		logger.LogIf(errors.New("无法获取用户"))
		return user.UserBasic{}
	}

	return userModel
}
