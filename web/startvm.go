package web

import (
	"vmcontroller/vm"

	"github.com/gin-gonic/gin"
)

var StartVm = func(c *gin.Context) {
	// 从请求参数中获取虚拟机名称
	vmname := c.Param("vmname")
	// 调用 vm 包的 StartVM 函数启动虚拟机
	err := vm.StartVM(vmname)
	if err != nil {
		// 若出错，返回 500 状态码和错误信息
		c.JSON(500, gin.H{
			"error": err.Error(),
			"code":  500,
		})
		return
	}
	// 启动成功，返回 200 状态码和成功消息
	c.JSON(200, gin.H{
		"message": "VM:" + vmname + "开机成功",
		"code":    0,
	})
}
