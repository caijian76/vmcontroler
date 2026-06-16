package web

import (
	"log"
	"vmcontroller/utils"
	"vmcontroller/vm"

	"github.com/gin-gonic/gin"
)

var Createvm = func(c *gin.Context) {
	var req utils.CreateVMRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	log.Printf("CreateVMRequest: %v", req)
	err := vm.CreateVM(&req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "VM created successfully"})
}
