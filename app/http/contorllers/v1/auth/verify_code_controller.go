package auth

import (
	"github.com/gin-gonic/gin"
	v1 "go-oj/app/http/contorllers/v1"
	"go-oj/app/requests"
	"go-oj/pkg/captcha"
	"go-oj/pkg/response"
	"go-oj/pkg/verifycode"
	"net/http"
)

type VerifyCodeController struct {
	v1.BaseAPIController
}

func (vc *VerifyCodeController) ShowCaptcha(c *gin.Context) {
	id, bs64, err := captcha.NewCaptcha().GenCaptcha()
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"captcha_id":    id,
		"captcha_image": bs64,
	})
}

func (vc *VerifyCodeController) VerifyCodePhone(c *gin.Context) {
	//1、绑定请求参数，检查参数合法性，检查图片验证码是否正确
	request := requests.VerifyCodePhoneRequest{}
	if ok := requests.Validate(c, &request, requests.VerifyCodePhone); !ok {
		return
	}

	//2、通过检验，发送短信验证码
	if ok := verifycode.NewVerifyCode().SendSMS(request.Phone); !ok {
		response.Abort500(c, "发送短信失败")
	} else {
		response.Success(c)
	}
}
