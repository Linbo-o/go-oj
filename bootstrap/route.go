package bootstrap

import (
	"github.com/gin-gonic/gin"
	"go-oj/routes"
	"net/http"
)

func SetupRoute(r *gin.Engine) {
	//1、注册全局中间件
	registerGlobalMiddlewares(r)
	//2、注册路由
	routes.RegisterAPIRoutes(r)
	//3、注册404路由
	register404Route(r)
}

// registerGlobalMiddlewares 注册全局中间件
func registerGlobalMiddlewares(r *gin.Engine) {
	r.Use(
		gin.Logger(),
		gin.Recovery(),
	)
}

// register404Route 注册404路由
func register404Route(r *gin.Engine) {
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error_code":    404,
			"error_message": "路由未定义，请确认 url 和请求方法是否正确。",
		})
	})
}
