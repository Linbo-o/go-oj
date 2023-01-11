package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterAPIRoutes(r *gin.Engine) {
	v1 := r.Group("/v1")
	{
		v1.POST("/test", func(context *gin.Context) {
			context.JSON(http.StatusOK, gin.H{
				"message": "hello world",
			})
		})
	}
}
