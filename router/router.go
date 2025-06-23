package router

import (
	"edulimitrate/handler"
	"edulimitrate/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// 初始化日志
	middleware.InitLogger("logs/track.log")  // 你可以用配置里的路径

	// 使用日志中间件
	r.Use(middleware.LogMiddleware())

	// 路由注册
	r.POST("/ratelimit/open", handler.OpenRegion)
	r.POST("/ratelimit/close", handler.CloseRegion)
	r.POST("/ratelimit/check", handler.CheckRegion)

	return r
}
