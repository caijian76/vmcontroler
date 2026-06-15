package web

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// 包装embed fs给static中间件
type embedFileSystem struct {
	http.FileSystem
}

func (e embedFileSystem) Exists(prefix string, path string) bool {
	_, err := e.Open(path)
	return err == nil
}

func serveIndex(ctx *gin.Context, fsys fs.FS) {
	data, err := fs.ReadFile(fsys, "index.html")
	if err != nil {
		ctx.String(http.StatusInternalServerError, "index.html not found")
		return
	}
	ctx.Data(http.StatusOK, "text/html; charset=utf-8", data)
}

func Route(route *gin.Engine) {

	// 提供 noVNC 静态文件服务
	// novncfs 是从嵌入的文件系统中提取的 noVNC-1.6.0 子目录
	novncfs, _ := fs.Sub(novncFS, "noVNC-1.6.0")
	// 将 /novnc 路径映射到 noVNC-1.6.0 目录的静态文件
	route.StaticFS("/novnc", http.FS(novncfs))

	uifs, _ := fs.Sub(uiFS, "ui/dist")
	assetsFS, _ := fs.Sub(uifs, "assets")
	route.StaticFS("/assets", http.FS(assetsFS))
	route.GET("/favicon.ico", func(ctx *gin.Context) {
		ctx.FileFromFS("favicon.ico", http.FS(uifs))
	})

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

			api.GET("/nodes", ListNode)

			// Vm Start API：启动指定名称的虚拟机
			api.GET("/vm/:vmname/start", StartVm)

			// Vm Stop API：停止指定名称的虚拟机
			api.GET("/vm/:vmname/stop", StopVM)

			// Vm Status API：获取指定名称虚拟机的状态
			api.GET("/vm/:vmname", StatusVm)

			// Vm Create API：创建一个新的虚拟机
			api.POST("/vm/:vmname/create", Createvm)
		}

	}
	route.NoRoute(func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if strings.HasPrefix(path, "/api") {
			ctx.JSON(404, gin.H{"message": "api not found"})
			return
		}
		serveIndex(ctx, uifs)
	})
}
