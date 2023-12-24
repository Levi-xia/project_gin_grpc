package main

import (
	"strconv"

	"com.levi/project-api/config"
	"com.levi/project-api/router"
	"com.levi/project-api/rpc"
	"com.levi/project-common/base"
	"com.levi/project-common/bootstrap"
	commonConfig "com.levi/project-common/config"
	"com.levi/project-common/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// 初始化全局配置
	commonConfig.InitConfig()
	// 初始化配置
	config.InitConfig()
	// 初始化日志
	base.InitLog()
	// 注册中间件
	middleware.InitMiddleware(r)
	// 初始化grpc客户端
	rpc.InitRpcClient()
	// 注册路由
	router.InitRouter(r)
	// 启动服务
	bootstrap.Run(r, config.Conf.Server.Name, ":"+strconv.Itoa(config.Conf.Server.Port), nil)
}
