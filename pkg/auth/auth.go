package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-oj/app/models/user"
	"go-oj/pkg/logger"
)

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
