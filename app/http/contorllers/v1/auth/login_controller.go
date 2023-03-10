package auth

import (
	"github.com/gin-gonic/gin"
	v1 "go-oj/app/http/contorllers/v1"
	"go-oj/app/requests"
	"go-oj/pkg/auth"
	jwtpkg "go-oj/pkg/jwt"
	"go-oj/pkg/response"
	"net/http"
)

type LoginController struct {
	v1.BaseAPIController
}

func (lc *LoginController) LoginByPhone(c *gin.Context) {
	//1、获取表单，并验证表单参数是否合法，检查验证码是否正确
	request := requests.LoginByPhoneRequest{}
	if ok := requests.Validate(c, &request, requests.LoginByPhone); !ok {
		return
	}

	//2、表单验证通过，尝试登录，检查账号是否已被注册
	u, err := auth.LoginByPhone(request.Phone)
	if err != nil {
		//登陆失败
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"token": jwtpkg.NewJWT().IssueToken(u.Identity, u.Name),
		})
	}
}

func (lc *LoginController) LoginByPassword(c *gin.Context) {
	request := requests.LoginByPasswordRequest{}
	if ok := requests.Validate(c, &request, requests.LoginByPassword); !ok {
		return
	}

	u, _ := auth.Attempt(request.LoginId)
	if u.Identity == "" {
		response.Unauthorized(c, "账号不存在或密码错误")
	} else {
		c.JSON(http.StatusOK, gin.H{
			"token": jwtpkg.NewJWT().IssueToken(u.Identity, u.Name),
		})
	}
}
