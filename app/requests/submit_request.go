package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

type SubmitRequest struct {
	ProblemIdentity string `json:"problem_identity" valid:"problem_identity"`
	Code            string `json:"code" valid:"code"`
}

func Submit(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"problem_identity": []string{"required"},
		"code":             []string{"required"},
	}

	messages := govalidator.MapData{
		"problem_identity": []string{"required:问题identity是必填项"},
		"code":             []string{"required:提交代码不能为空"},
	}

	return validate(rules, messages, data)
}
