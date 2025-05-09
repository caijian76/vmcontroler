package web

import (
	"vmcontroller/vm"

	"github.com/gin-gonic/gin"
)

var Listvm = func(c *gin.Context) {
	// 调用 vm 包的 ListVm 函数获取虚拟机列表
	vml, err := vm.ListVm()
	if err != nil {
		// 若出错，返回 500 状态码和错误信息
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// 成功获取列表，返回 200 状态码和虚拟机列表
	c.JSON(200, vml)
}
