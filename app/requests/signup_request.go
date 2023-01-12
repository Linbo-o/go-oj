package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"go-oj/app/requests/validators"
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

type SignupEmailExistRequest struct {
	Email string `json:"email" valid:"email"`
}

func SignupExistEmail(data interface{}, c *gin.Context) map[string][]string {
	//1、定制规则
	rules := govalidator.MapData{
		"email": []string{"required", "email"},
	}

	//2、定制错误信息
	messages := govalidator.MapData{
		"email": []string{
			"required:邮箱是必填项",
			"email:必须要符合email格式",
		},
	}

	return validate(rules, messages, data)
}

type SignupUsingPhoneRequest struct {
	Phone           string `json:"phone,omitempty" valid:"phone"`
	VerifyCode      string `json:"verify_code,omitempty" valid:"verify_code"`
	Name            string `valid:"name" json:"name"`
	Password        string `valid:"password" json:"password,omitempty"`
	PasswordConfirm string `valid:"password_confirm" json:"password_confirm,omitempty"`
}

func SignupUsingPhone(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"phone":            []string{"required", "digits:11", "not_exists:user_basic,phone"},
		"name":             []string{"required", "alpha_num", "between:3,20", "not_exists:user_basic,name"},
		"password":         []string{"required", "min:6"},
		"password_confirm": []string{"required"},
		"verify_code":      []string{"required"},
	}

	messages := govalidator.MapData{
		"phone": []string{
			"required:手机号为必填项，参数名称 phone",
			"digits:手机号长度必须为 11 位的数字",
		},
		"name": []string{
			"required:用户名为必填项",
			"alpha_num:用户名格式错误，只允许数字和英文",
			"between:用户名长度需在 3~20 之间",
		},
		"password": []string{
			"required:密码为必填项",
			"min:密码长度需大于 6",
		},
		"password_confirm": []string{
			"required:确认密码框为必填项",
		},
		"verify_code": []string{
			"required:验证码答案必填",
		},
	}

	//1、验证表单格式是否正确，参数是否合法
	errs := validate(rules, messages, data)

	//2、验证密码和确认密码是否相等
	data_ := data.(*SignupUsingPhoneRequest)
	if data_.Password != data_.PasswordConfirm {
		errs["password_confirm"] = append(errs["password_confirm"], "两次输入密码不相同")
	}
	//3、验证验证码是否正确
	errs = validators.ValidateVerifyCode(data_.Phone, data_.VerifyCode, errs)

	return errs
}
