package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"go-oj/app/models/testcase"
	"go-oj/app/requests/validators"
)

//ProblemCreateRequest 创建问题
type ProblemCreateRequest struct {
	Title      string               `json:"title"    db:"title" valid:"title"`
	Content    string               `json:"content"  valid:"content"`
	Categories []int                `json:"categories" valid:"categories"`
	MaxRuntime int                  `json:"max_runtime"  valid:"max_runtime"`
	MaxMem     int                  `json:"max_mem"  valid:"max_mem"`
	TestCases  []*testcase.TestCase `json:"test_cases" valid:"test_cases"`
}

func ProblemCreate(data interface{}, c *gin.Context) map[string][]string {
	// 先验证是否为管理员
	errs := validators.ValidateIsAdmin(c)
	if errs != nil {
		return errs
	}
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
		"title":       []string{"required:title是必填项"},
		"content":     []string{"required:content是必填项"},
		"categories":  []string{"required:categories是必填项"},
		"max_runtime": []string{"required:content是必填项"},
		"max_mem":     []string{"required:max_mem是必填项"},
		"test_cases":  []string{"required:test_case是必填项"},
	}

	//3、对上面规则进行验证，返回错误信息
	return validate(rules, message, data)
}

//ProblemModifyRequest 修改题目
type ProblemModifyRequest struct {
	Identity   string `json:"identity" valid:"identity"`
	Title      string `json:"title"    db:"title" valid:"title"`
	Content    string `json:"content"  valid:"content"`
	MaxRuntime int    `json:"max_runtime"  valid:"max_runtime"`
	MaxMem     int    `json:"max_mem"  valid:"max_mem"`
}

func ProblemModify(data interface{}, c *gin.Context) map[string][]string {
	// 先验证是否为管理员
	errs := validators.ValidateIsAdmin(c)
	if errs != nil {
		return errs
	}
	//1、定制规则
	rules := govalidator.MapData{
		"identity":    []string{"required"},
		"title":       []string{"required", "not_exists:problem_basic,title"},
		"content":     []string{"required"},
		"max_runtime": []string{"required"},
		"max_mem":     []string{"required"},
	}

	//2、定制错误信息
	message := govalidator.MapData{
		"identity":    []string{"required:identity是必填项"},
		"title":       []string{"required:title是必填项"},
		"content":     []string{"required:content是必填项"},
		"max_runtime": []string{"required:content是必填项"},
		"max_mem":     []string{"required:max_mem是必填项"},
	}

	//3、对上面规则进行验证，返回错误信息
	return validate(rules, message, data)
}

//GetProblemListRequest 获取题目列表
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

//GetProblemDetailRequest 获取题目详情
type GetProblemDetailRequest struct {
	Identity string `json:"identity" valid:"identity"`
}

func GetProblemDetail(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"identity": []string{"required"},
	}

	message := govalidator.MapData{
		"identity": []string{"required:identity是必填项"},
	}

	return validate(rules, message, data)
}
