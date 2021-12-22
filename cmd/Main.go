package main

import (
	"github.com/gin-gonic/gin"
	pkg2 "github.com/glazzes/borealys/pkg"
)

func main() {
	router := gin.Default()

	router.GET("/api/languages", pkg2.SimpleLanguageService{}.GetAll)
	router.POST("/api/run", pkg2.SimpleCodeRunnerService{}.RunCode)

	err := router.Run(":5000")

	if err != nil{
		panic(err)
	}
}
