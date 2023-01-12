package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"go-oj/app/models/testcase"
)

type ProblemCreateRequest struct {
	Title      string               `json:"title"    db:"title" valid:"title"`
	Content    string               `json:"content"  valid:"content"`
	Categories []int                `json:"categories" valid:"categories"`
	MaxRuntime int                  `json:"max_runtime"  valid:"max_runtime"`
	MaxMem     int                  `json:"max_mem"  valid:"max_mem"`
	TestCases  []*testcase.TestCase `json:"test_cases" valid:"test_cases"`
}

func ProblemCreate(data interface{}, c *gin.Context) map[string][]string {
	//1、定制规则
	rules := govalidator.MapData{
		"title":       []string{"required", "not_exists:problem_basic,title"},
		"content":     []string{"required"},
		"categories":  []string{"required"},
		"max_runtime": []string{"required"},
		"max_mem":     []string{"required"},
		"test_cases":  []string{"required"},
	}

	//2、定制错误信息
	message := govalidator.MapData{
		"title": []string{
			"required:title是必填项",
			"not_exists:该问题题目已经存在",
		},
		"content":     []string{"required:content是必填项"},
		"categories":  []string{"required:categories是必填项"},
		"max_runtime": []string{"required:content是必填项"},
		"max_mem":     []string{"required:max_mem是必填项"},
		"test_cases":  []string{"required:test_case是必填项"},
	}

	//3、对上面规则进行验证，返回错误信息
	return validate(rules, message, data)
}

type GetProblemListRequest struct {
	Size int `json:"size" valid:"size"`
	Page int `json:"page" valid:"page"`
}

func GetProblemList(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"size": []string{"required"},
		"page": []string{"required"},
	}

	messages := govalidator.MapData{
		"size": []string{"required:size是必填项"},
		"page": []string{"required:page是必填项"},
	}

	return validate(rules, messages, data)
}
