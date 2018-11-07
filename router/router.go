package router

import (
	"github.com/gin-gonic/gin"
	LogController "shell-exec/http/controllers/api/v1"
)

func InitRouter() *gin.Engine {

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		// 打包日志下载
		v1.GET("/pack_log", LogController.Packing)
	}

	return r
}
