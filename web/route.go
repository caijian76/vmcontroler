package web

import (
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Route(route *gin.Engine) {

	// 提供 noVNC 静态文件服务
	// newfs 是从嵌入的文件系统中提取的 noVNC-1.6.0 子目录
	newfs, _ := fs.Sub(novncFS, "noVNC-1.6.0")
	// 将 /novnc 路径映射到 noVNC-1.6.0 目录的静态文件
	route.StaticFS("/novnc", http.FS(newfs))

	// 处理 VNC WebSocket 连接的路由
	route.GET("/vm/:vmname/vnc", WsVnc1)

	// 登录接口
	route.POST("/login", Login)

	// 创建一个需要 JWT 验证的路由组
	authGroup := route.Group("/")
	authGroup.Use(JWTAuthMiddleware())
	{
		// 创建一个 /api 路由组，用于管理所有 API 接口
		api := authGroup.Group("/api")
		{
			// Vm List API：获取虚拟机列表
			api.GET("/vm", Listvm)

			// Vm Start API：启动指定名称的虚拟机
			api.GET("/vm/:vmname/start", StartVm)

			// Vm Stop API：停止指定名称的虚拟机
			api.GET("/vm/:vmname/stop", StopVM)

			// Vm Status API：获取指定名称虚拟机的状态
			api.GET("/vm/:vmname", StatusVm)
		}

	}
}
