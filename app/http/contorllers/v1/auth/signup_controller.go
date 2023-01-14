package auth

import (
	"github.com/gin-gonic/gin"
	"go-oj/app/http/contorllers/v1"
	"go-oj/app/models/user"
	"go-oj/app/requests"
	"go-oj/pkg/hash"
	"go-oj/pkg/helpers"
	jwtpkg "go-oj/pkg/jwt"
	"go-oj/pkg/logger"
	"go-oj/pkg/response"
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

func (sc *SignupController) SignupUsingPhone(c *gin.Context) {
	//1、绑定表单数据，验证表单
	request := requests.SignupUsingPhoneRequest{}
	if ok := requests.Validate(c, &request, requests.SignupUsingPhone); !ok {
		return
	}
	//2、验证成功,创建数据
	logger.DebugJSON("signup_using_phone", "创建数据", request)
	u := user.UserBasic{
		Identity:  helpers.GetUUID(),
		Name:      request.Name,
		Password:  hash.Encrypt(request.Password),
		Phone:     request.Phone,
		PassNum:   0,
		SubmitNum: 0,
		IsAdmin:   0,
	}
	if ok := u.Create(); !ok {
		response.Abort500(c, "创建数据失败")
	} else {
		c.JSON(http.StatusOK, gin.H{
			"token": jwtpkg.NewJWT().IssueToken(u.Identity, u.Name),
		})
	}
}
