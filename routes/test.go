package routes

import (
	"github.com/gin-gonic/gin"
	"go-oj/app/http/middlewares"
	"go-oj/app/models/user"
	"go-oj/pkg/auth"
	"go-oj/pkg/helpers"
	jwtpkg "go-oj/pkg/jwt"
	"go-oj/pkg/response"
	"net/http"
	"sync"
)

func TestRoute(r *gin.Engine) {
	testGroup := r.Group("/test")
	{
		//1、通过jwt获取用户信息
		testGroup.GET("/user-detail", middlewares.AuthJWT(), userDetail)
		//2、创建管理员账号
		testGroup.POST("/admin-create", createAdminUser)
	}
}

var once sync.Once

func userDetail(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"user": auth.CurrentUser(c),
	})
}

func createAdminUser(c *gin.Context) {
	once.Do(func() {
		u := user.UserBasic{
			Identity:  helpers.GetUUID(),
			Name:      "admin-lb",
			Password:  "makabaka",
			Phone:     "12345678910",
			PassNum:   100,
			SubmitNum: 100,
			IsAdmin:   1,
		}
		if ok := u.Create(); !ok {
			response.Abort500(c, "创建数据失败")
		} else {
			c.JSON(http.StatusOK, gin.H{
				"token": jwtpkg.NewJWT().IssueToken(u.Identity, u.Name),
			})
		}
	})
}
