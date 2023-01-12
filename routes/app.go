package routes

import (
	"github.com/gin-gonic/gin"
	"go-oj/app/http/contorllers/v1/auth"
	"go-oj/app/http/contorllers/v1/problem"
)

func RegisterAPIRoutes(r *gin.Engine) {
	v1 := r.Group("/v1")
	{
		authGroup := v1.Group("/auth")
		{
			suc := new(auth.SignupController)
			// 判断手机是否已注册
			authGroup.POST("/signup/phone/exist", suc.IsPhoneExist)
			authGroup.POST("/signup/email/exist", suc.IsEmailExist)
			authGroup.POST("/signup/using-phone", suc.SignupUsingPhone)

			loc := new(auth.LoginController)
			authGroup.POST("/login/using-phone", loc.LoginByPhone)

			vcc := new(auth.VerifyCodeController)
			authGroup.POST("/verify-code/captcha", vcc.ShowCaptcha)
			authGroup.POST("/verify-code/phone", vcc.VerifyCodePhone)
		}

		problemGroup := v1.Group("/problem")
		{
			pbc := new(problem.ProblemController)
			problemGroup.POST("/create", pbc.ProblemCreate)
			problemGroup.POST("/problem-list", pbc.GetProblemList)
			problemGroup.POST("/problem-detail", pbc.GetProblemDetail)
		}
	}
}
