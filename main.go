package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-oj/bootstrap"
	"go-oj/config"
	pkgConfig "go-oj/pkg/config"
)

func init() {
	config.Init()
}

func main() {
	// 配置初始化，依赖命令行 --env 参数
	var env string
	flag.StringVar(&env, "env", "", "加载 .env 文件，如 --env=testing 加载的是 .env.testing 文件")
	flag.Parse()
	pkgConfig.InitConfig(env)

	r := gin.New()
	//1、设置日志
	bootstrap.SetupLogger()
	//gin.SetMode(gin.ReleaseMode)
	//2、设置sql数据库
	bootstrap.SetupDatabase()
	//3、设置redis
	bootstrap.SetupRedis()
	//4、设置路由
	bootstrap.SetupRoute(r)

	err := r.Run(fmt.Sprintf(":%v", pkgConfig.Get("app.port")))
	if err != nil {
		fmt.Println(err)
	}
}
