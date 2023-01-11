package main

import (
	"github.com/gin-gonic/gin"
	"go-oj/bootstrap"
)

func main() {
	r := gin.New()
	bootstrap.SetupRoute(r)
	r.Run(":3000")
}
