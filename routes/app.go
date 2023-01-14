package routes

import (
	"github.com/gin-gonic/gin"
	"go-oj/app/http/contorllers/v1/auth"
	"go-oj/app/http/contorllers/v1/problem"
	"go-oj/app/http/middlewares"
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
			authGroup.POST("/login/using-password", loc.LoginByPassword)

			vcc := new(auth.VerifyCodeController)
			authGroup.POST("/verify-code/captcha", vcc.ShowCaptcha)
			authGroup.POST("/verify-code/phone", vcc.VerifyCodePhone)
		}

		problemGroup := v1.Group("/problem", middlewares.AuthJWT())
		{
			pbc := new(problem.ProblemController)
			problemGroup.POST("/create", pbc.ProblemCreate)
			problemGroup.POST("/modify", pbc.ProblemModify)
			problemGroup.POST("/problem-list", pbc.GetProblemList)
			problemGroup.POST("/problem-detail", pbc.GetProblemDetail)
			problemGroup.POST("/problem-judge", pbc.ProblemJudge)
			problemGroup.POST("/judge-status", pbc.ProblemJudgeResult)
		}

		//用于测试
		TestRoute(r)
	}
}
