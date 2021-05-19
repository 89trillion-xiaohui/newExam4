package main

import (
	"test3/src/controller"
	"test3/src/service"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	service.GetMongoClient()

	r.POST("/createCode", controller.Create)

	r.GET("/inquire", controller.Inquire)

	r.GET("/client", controller.Client)

	r.GET("/User", controller.Log)

	r.Run()
}
