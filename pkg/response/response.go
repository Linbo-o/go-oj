package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Abort500(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"message": message,
	})
}

func Success(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "操作成功!",
	})
}

// Unauthorized 响应 401，未传参 msg 时使用默认消息
// 登录失败、jwt 解析失败时调用
func Unauthorized(c *gin.Context, msg ...string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"message": msg,
	})
}
