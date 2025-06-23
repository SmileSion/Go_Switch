package router

import (
	"edulimitrate/handler"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/ratelimit/open", handler.OpenRegion)
	r.POST("/ratelimit/close", handler.CloseRegion)
	r.POST("/ratelimit/check", handler.CheckRegion)

	return r
}
