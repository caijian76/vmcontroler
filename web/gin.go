// web 包负责处理 Web 相关的功能，包括 HTTP 路由和 WebSocket 连接管理
package web

import (
	"embed"
	"time"

	"github.com/gin-gonic/gin"
)

// 定义常量，用于设置 WebSocket 连接的空闲超时和 Ping 消息发送间隔
const (
	// idleTimeout 为 WebSocket 连接的空闲超时时间，超过该时间无活动则关闭连接
	idleTimeout = 30 * time.Second
	// pingInterval 为定期发送 Ping 消息的间隔，应小于空闲超时时间，确保能及时检测连接状态
	pingInterval = (idleTimeout * 9) / 10
)

// 使用 go:embed 指令将 noVNC-1.6.0 目录嵌入到程序中，
// 这样在编译时该目录的内容会被打包进可执行文件
// novncFS 是一个嵌入的文件系统，包含 noVNC-1.6.0 目录的所有文件
//
//go:embed noVNC-1.6.0
var novncFS embed.FS

// WebStart 函数用于启动 Web 服务，设置 HTTP 路由和 WebSocket 处理逻辑
// virtClient 是 KubeVirt 客户端，用于与 KubeVirt 集群交互
// namespace 是操作的命名空间
func WebStart() {

	//	gin.DisableConsoleColor()
	//	gin.SetMode(gin.ReleaseMode)
	// 创建一个新的 Gin 路由实例
	route := gin.New()
	//route.Use(gin.Logger())
	// 使用 Gin 的 Recovery 中间件，捕获并恢复处理请求时的 panic
	route.Use(gin.Recovery())
	Route(route)

	// 启动 Gin 服务器，监听 8080 端口
	route.Run(":8080")
}
