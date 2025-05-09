package web

import (
	"vmcontroller/vm"

	"github.com/gin-gonic/gin"
)

var StopVM = func(c *gin.Context) {
	// 从请求参数中获取虚拟机名称
	vmname := c.Param("vmname")
	// 调用 vm 包的 StopVM 函数停止虚拟机
	err := vm.StopVM(vmname, 0)
	if err != nil {
		// 若出错，返回 500 状态码和错误信息
		c.JSON(500, gin.H{"error": err.Error(),
			"code": 500,
		})
		return
	}
	// 停止成功，返回 200 状态码和成功消息
	c.JSON(200, gin.H{"code": 0,
		"message": "VM:" + vmname + "已经关闭"})
}
