package web

import (
	"vmcontroller/vm"

	"github.com/gin-gonic/gin"
)

var StatusVm = func(c *gin.Context) {
	// 从请求参数中获取虚拟机名称
	vmname := c.Param("vmname")
	// 调用 vm 包的 StatusVM 函数获取虚拟机状态
	reportvm, err := vm.StatusVM(vmname)
	if err != nil {
		// 若出错，返回 500 状态码和错误信息
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// 成功获取状态，返回 200 状态码和虚拟机状态信息
	c.JSON(200, reportvm)
}
