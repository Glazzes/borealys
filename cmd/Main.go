package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	pkg2 "github.com/glazzes/borealys/pkg"
	"time"
)

var (
	initializerService = &pkg2.SimpleInitializerService{}
	languageService = &pkg2.SimpleLanguageService{}
	codeRunnerService = &pkg2.SimpleCodeRunnerService{}
)

func init(){
	initializerService.CreateRunnersGroup()
	initializerService.CreateExecutorUsers()
	initializerService.SetUpBinaries()
	languageService.SaveAll()
}

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:19006"},
		AllowMethods: []string{"POST", "GET", "OPTIONS"},
		MaxAge: 12 * time.Hour,
	}))

	router.GET("/api/languages", languageService.GetAll)
	router.POST("/api/run", codeRunnerService.RunCode)

	err := router.Run(":5000")

	if err != nil{
		panic(err)
	}
}
