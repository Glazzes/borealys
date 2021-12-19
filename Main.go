package main

import (
	"./handler"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/docker", handler.HandleInDockerCode)
	router.GET("/system", handler.HandleInSystemCode)
	router.GET("/langs", handler.GetSupportedLanguages)

	err := router.Run(":5000")
	if err != nil{
		panic(err)
	}
}
