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
