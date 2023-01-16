package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"go-oj/pkg/logger"
	"io"
)

type SubmitRequest struct {
	ProblemIdentity string `form:"problem_identity" json:"problem_identity" valid:"problem_identity"`
	Code            []byte
}

func Submit(data interface{}, c *gin.Context) map[string][]string {
	code, err := io.ReadAll(c.Request.Body)
	if err != nil || len(code) == 0 {
		logger.LogIf(err)
		return map[string][]string{
			"code": []string{"获提交代码失败，请检查提交代码是否为空"},
		}
	}
	data.(*SubmitRequest).Code = code
	rules := govalidator.MapData{
		"problem_identity": []string{"required"},
	}

	messages := govalidator.MapData{
		"problem_identity": []string{"required:问题identity是必填项"},
	}

	return validate(rules, messages, data)
}
