package problem

import (
	"github.com/gin-gonic/gin"
	v1 "go-oj/app/http/contorllers/v1"
	"go-oj/app/models/problem"
	pc "go-oj/app/models/problem-category"
	"go-oj/app/models/testcase"
	"go-oj/app/requests"
	"go-oj/pkg/helpers"
	"go-oj/pkg/logger"
	"go-oj/pkg/response"
	"net/http"
)

type ProblemController struct {
	v1.BaseAPIController
}

func (pro *ProblemController) ProblemCreate(c *gin.Context) {
	//1、绑定表单参数，验证表单合法性
	request := requests.ProblemCreateRequest{}
	if ok := requests.Validate(c, &request, requests.ProblemCreate); !ok {
		return
	}
	logger.Dump(request.TestCases)
	//2、创建建数据
	//2.1、问题基础数据
	p := problem.ProblemBasic{
		Identity:   helpers.GetUUID(),
		Title:      request.Title,
		Content:    request.Content,
		MaxRuntime: request.MaxRuntime,
		MaxMem:     request.MaxMem,
	}
	//2.2、问题分类数据处理
	problemCategories := make([]*pc.ProblemCategory, 0)
	for _, id := range request.Categories {
		problemCategories = append(problemCategories, &pc.ProblemCategory{
			ProblemIdentity: p.Identity,
			CategoryId:      uint(id),
		})
	}
	p.Categories = problemCategories
	//2.3、问题测试用例数据处理
	testCases := make([]*testcase.TestCase, 0)
	for _, ca := range request.TestCases {
		testCases = append(testCases, &testcase.TestCase{
			Identity:        helpers.GetUUID(),
			ProblemIdentity: p.Identity,
			Input:           ca.Input,
			Output:          ca.Output,
		})
	}
	p.TestCases = testCases
	//2.4将数据写入数据库
	if ok := p.Create(); !ok {
		response.Abort500(c, "写入数据失败")
	} else {
		response.Success(c)
	}
}

func (pro *ProblemController) GetProblemList(c *gin.Context) {
	//1、绑定参数并检验
	request := requests.GetProblemListRequest{}
	if ok := requests.Validate(c, &request, requests.GetProblemList); !ok {
		return
	}

	//2、获取列表
	list := problem.GetProblemList(request.Size, request.Page)
	c.JSON(http.StatusOK, gin.H{
		"problem_list": list,
	})
}

func (pro *ProblemController) GetProblemDetail(c *gin.Context) {
	//1、获取表单，验证表单
	request := requests.GetProblemDetailRequest{}
	if ok := requests.Validate(c, &request, requests.GetProblemDetail); !ok {
		return
	}

	//2、获取问题信息，返回客户端
	p := problem.GetProblemDetail(request.Identity)
	if p.Identity == "" {
		response.Abort500(c, "获取信息失败")
	} else {
		c.JSON(http.StatusOK, gin.H{
			"details": p,
		})
	}
}
