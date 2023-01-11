package auth

import (
	"github.com/gin-gonic/gin"
	"go-oj/app/http/v1/contorllers"
	"go-oj/app/models/user"
	"go-oj/app/requests"
	"net/http"
)

type SignupController struct {
	contorllers.BaseAPIController
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
