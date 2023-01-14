package validators

import (
	"github.com/gin-gonic/gin"
	"go-oj/pkg/logger"
	"go-oj/pkg/verifycode"
)

// ValidateVerifyCode 自定义规则，验证『手机/邮箱验证码』
func ValidateVerifyCode(key, answer string, errs map[string][]string) map[string][]string {
	if ok := verifycode.NewVerifyCode().CheckAnswer(key, answer); !ok {
		errs["verify_code"] = append(errs["verify_code"], "验证码错误")
	}
	return errs
}

func ValidateIsAdmin(c *gin.Context) map[string][]string {
	is, exist := c.Get("is_admin")
	errs := make(map[string][]string)
	if !exist {
		logger.WarnString("validate", "is_admin", "获取参数失败")
		errs["授权认证"] = append(errs["授权认证"], "获取参数失败")
		return errs
	}
	isAdmin := is.(int)
	if isAdmin != 1 {
		logger.WarnString("validate", "is_admin", "滚，你不是管理员[鄙视]!")
		errs["授权认证"] = append(errs["授权认证"], "滚，你不是管理员[鄙视]!")
		return errs
	}
	return nil
}
