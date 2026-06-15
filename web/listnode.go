package web

import (
	"vmcontroller/vm"

	"github.com/gin-gonic/gin"
)

func ListNode(c *gin.Context) {
	nodes, err := vm.ListNode()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, nodes)
}
