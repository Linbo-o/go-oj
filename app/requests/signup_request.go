package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

type SignupPhoneExistRequest struct {
	Phone string `db:"phone" json:"phone" valid:"phone"`
}

func SignupExistPhone(data interface{}, c *gin.Context) map[string][]string {
	//1、定制规则
	rules := govalidator.MapData{
		"phone": []string{"required", "digits:11"},
	}

	//2、定制错误信息
	messages := govalidator.MapData{
		"phone": []string{
			"required:手机号码为必填项",
			"digits:号码长度必须为11位",
		},
	}

	return validate(rules, messages, data)
}
