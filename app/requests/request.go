package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"go-oj/pkg/logger"
	"net/http"
)

type ValidatorFunc func(interface{}, *gin.Context) map[string][]string

func Validate(c *gin.Context, obj interface{}, handler ValidatorFunc) bool {
	err := c.ShouldBind(obj)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
	}
	errs := handler(obj, c)
	if len(errs) > 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errs,
		})
		logger.ErrorJSON("validate", "验证失败", errs)
		return false
	}
	return true
}

func validate(rules, message govalidator.MapData, data interface{}) map[string][]string {
	opts := govalidator.Options{
		Data:          data,
		Rules:         rules,
		Messages:      message,
		TagIdentifier: "valid",
	}
	return govalidator.New(opts).ValidateStruct()
}
