package web

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"
	"vmcontroller/vm"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var WsVnc1 = func(c *gin.Context) {
	// VNC WebSocket 服务配置
	var upgader = websocket.Upgrader{
		// 读取缓冲区大小
		ReadBufferSize: 1024,
		// 写入缓冲区大小
		WriteBufferSize: 1024,
		// 检查 Origin 头，允许所有来源的连接
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	// 创建一个可取消的上下文，用于控制 goroutine 的生命周期
	ctx, cancel := context.WithCancel(context.Background())

	// 将 HTTP 连接升级为 WebSocket 连接
	conn, err := upgader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// 升级失败，记录错误日志并取消上下文
		log.Printf("Failed to upgrade connection: %s", err.Error())
		cancel()
		return
	}

	// 注册 WebSocket 关闭处理函数，当客户端关闭连接时触发
	conn.SetCloseHandler(func(code int, text string) error {
		// 记录客户端关闭连接的信息
		log.Printf("WebSocket connection closed by client: code %d, text %s", code, text)
		// 取消上下文
		cancel()
		return nil
	})

	// 设置 Pong 消息处理函数和读取截止时间，用于检测连接是否空闲
	conn.SetReadDeadline(time.Now().Add(idleTimeout))

	conn.SetPongHandler(func(string) error {
		// 收到 Pong 消息后，更新读取截止时间
		conn.SetReadDeadline(time.Now().Add(idleTimeout))
		return nil
	})

	// 获取指定虚拟机的 VNC 连接
	vnc, err := vm.Getvnc(c.Param("vmname"))
	if err != nil {
		// 获取 VNC 连接失败，记录错误日志，取消上下文并关闭 WebSocket 连接
		log.Printf("can't access VMI %s: %s\n", c.Param("vmname"), err.Error())
		cancel()
		conn.Close()
		return
	}

	// 创建两个管道，用于在 VNC 连接和 WebSocket 之间传输数据

	// 定义一个延迟执行的函数，确保在函数结束时关闭所有资源
	defer func() {

		conn.Close()
		cancel()
	}()

	// 创建一个错误通道，用于接收各个 goroutine 中的错误信息
	errChan := make(chan error, 3)

	vncws := vnc.AsConn()
	defer vncws.Close()
	// 记录 VNC 客户端连接成功的信息
	log.Printf("VNC Client connected to: %s", c.Param("vmname"))

	// 创建一个定时器，用于定期发送 Ping 消息
	pingTicker := time.NewTicker(pingInterval)
	// 确保在函数结束时停止定时器
	defer pingTicker.Stop()

	// 启动一个 goroutine 用于定期发送 Ping 消息
	go func() {
		for {
			select {
			case <-pingTicker.C:
				// 定时器触发，发送 Ping 消息
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					// 发送失败，记录错误日志，取消上下文并退出 goroutine
					log.Printf("Failed to send ping: %v", err)
					cancel()
					return
				}
			case <-ctx.Done():
				// 上下文取消，退出 goroutine

				return
			}
		}
	}()

	// 启动一个 goroutine 用于从 pipeOutReader 读取数据并写入 WebSocket
	go func() {
		// 定义一个延迟执行的函数，在 goroutine 结束时关闭资源并发送 EOF 到错误通道
		defer func() {

			conn.Close()
			errChan <- io.EOF

		}()
		// 创建一个缓冲区，用于存储从 pipeOutReader 读取的数据
		buf := make([]byte, 4096)
		for {
			select {
			case <-ctx.Done():
				// 上下文取消，记录日志并退出 goroutine
				log.Println("Context canceled, exiting goroutine for reading from pipeOutReader")
				return
			default:
				// 从 pipeOutReader 读取数据
				n, err := vncws.Read(buf)
				if err != nil {

					// 退出 goroutine
					return
				}
				// 将读取的数据以二进制消息的形式写入 WebSocket
				if err := conn.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
					// 写入失败，记录错误日志并退出 goroutine
					log.Printf("Error writing to WebSocket: %v", err)
					return
				}
			}
		}

	}()

	// 启动一个 goroutine 用于从 WebSocket 读取数据并写入 pipeInWriter
	go func() {
		// 定义一个延迟执行的函数，在 goroutine 结束时关闭资源并发送 EOF 到错误通道
		defer func() {

			conn.Close()
			errChan <- io.EOF

		}()
		for {
			select {
			case <-ctx.Done():
				// 上下文取消，记录日志并退出 goroutine
				log.Println("Context canceled, exiting goroutine for reading from WebSocket")
				return
			default:
				// 从 WebSocket 获取下一个消息的读取器
				_, reader, err := conn.NextReader()
				if err != nil {

					return
				}
				// 将读取器中的数据复制到 pipeInWriter
				if _, err := io.Copy(vncws, reader); err != nil {
					// 复制失败，记录错误日志并退出 goroutine
					log.Printf("Error writing to pipeInWriter: %v", err)
					return
				}
			}
		}
	}()

	// 等待错误通道接收到第一个错误信息
	<-errChan
	// 触发上下文取消，通知其他 goroutine 退出
	cancel()
	// 记录 VNC 客户端断开连接的信息
	log.Printf("VNC Client disconnected from: %s", c.Param("vmname"))
}
