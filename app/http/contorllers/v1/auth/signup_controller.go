package auth

import (
	"github.com/gin-gonic/gin"
	"go-oj/app/http/contorllers/v1"
	"go-oj/app/models/user"
	"go-oj/app/requests"
	"net/http"
)

type SignupController struct {
	v1.BaseAPIController
}

func (sc *SignupController) IsPhoneExist(c *gin.Context) {
	//1、获取信息，验证表单
	request := requests.SignupPhoneExistRequest{}
	if ok := requests.Validate(c, &request, requests.SignupExistPhone); !ok {
		return
	}

	//查看是否存在，返回信息
	c.JSON(http.StatusOK, gin.H{
		"exist": user.IsPhoneExist(request.Phone),
	})
}

func (sc *SignupController) IsEmailExist(c *gin.Context) {
	//1、获取信息，验证表单
	request := requests.SignupEmailExistRequest{}
	if ok := requests.Validate(c, &request, requests.SignupExistEmail); !ok {
		return
	}

	//查看是否存在，返回信息
	c.JSON(http.StatusOK, gin.H{
		"exist": user.IsEmailExist(request.Email),
	})
}
