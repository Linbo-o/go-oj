package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"go-oj/pkg/captcha"
)

type VerifyCodePhoneRequest struct {
	CaptchaID     string `json:"captcha_id,omitempty" valid:"captcha_id"`
	CaptchaAnswer string `json:"captcha_answer,omitempty" valid:"captcha_answer"`

	Phone string `json:"phone,omitempty" valid:"phone"`
}

func VerifyCodePhone(data interface{}, c *gin.Context) map[string][]string {
	//1、定制规则
	rules := govalidator.MapData{
		"captcha_id":     []string{"required"},
		"captcha_answer": []string{"required"},
		"phone":          []string{"required", "digits:11"},
	}

	//2、定制错误信息
	messages := govalidator.MapData{
		"captcha_id":     []string{"required:图案验证码id为必填项"},
		"captcha_answer": []string{"required:图案验证码我必填项"},
		"phone": []string{
			"required:手机号码为必填项",
			"digits:手机号码必须为11位数字",
		},
	}

	errs := validate(rules, messages, data)

	//3、验证图片验证码是否正确
	data_ := data.(*VerifyCodePhoneRequest)
	if ok := captcha.NewCaptcha().Verify(data_.CaptchaID, data_.CaptchaAnswer, false); !ok {
		errs["captcha"] = append(errs["captcha"], "图案验证码不正确")
	}
	return errs
}
