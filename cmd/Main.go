package main

import (
	"github.com/gin-gonic/gin"
	"github.com/glazzes/borealys/languages"
	pkg2 "github.com/glazzes/borealys/pkg"
)

var (
	initializerService = &languages.SimpleInitializerService{}
	languageService = &pkg2.SimpleLanguageService{}
	codeRunnerService = &pkg2.SimpleCodeRunnerService{}
)

func init(){
	languageService.SaveAll()
	languageService.SaveRunnableLanguages()
	initializerService.CreateExecutorUsers()
	initializerService.SetUpBinaries()
}

func main() {
	router := gin.Default()

	router.GET("/api/languages", languageService.GetAll)
	router.POST("/api/run", codeRunnerService.RunCode)

	err := router.Run(":5000")

	if err != nil{
		panic(err)
	}
}
