package auth

import (
	"github.com/gin-gonic/gin"
	v1 "go-oj/app/http/contorllers/v1"
	"go-oj/pkg/captcha"
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
