package main

import (
	"../handler"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/api/languages")
	router.POST("/execute", handler.HandleCodeUpload)

	err := router.Run(":5000")

	if err != nil{
		panic(err)
	}
}
