package requests

import (
	"github.com/gin-gonic/gin"
)

type ValidatorFunc func(interface{}, *gin.Context) map[string][]string
