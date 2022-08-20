package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ram2104/go-authentication/controller"
	"github.com/ram2104/go-authentication/initializer"
)

func init() {
	initializer.LoadENVData()
	initializer.CreateDBConnection()
	initializer.SyncDatabase()
}

func main() {
	r := gin.Default()
	r.POST("/signup", controller.Signup)
	r.POST("/signin", controller.Login)
	r.Run()
}
