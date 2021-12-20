package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os/exec"
	"strings"
)

type SupportedLanguage struct {
	Name string
	SupportedVersions []string
	Extension string
	RequiresCompilation bool
	IsInterpreted bool
}

type ExecutableCode struct {
	Language string `json:"language" binding:"required"`
	Version string `json:"version" binding:"required"`
	Code []string `json:"code" binding:"required"`
}

func GetSupportedLanguages(context *gin.Context){
	supportedLanguages := []SupportedLanguage{
		{Name: "Java", SupportedVersions: []string{"17"}},
		{Name: "node", SupportedVersions: []string{"16.13.0", "17.3.0"}},
		{Name: "python", SupportedVersions: []string{"3.10.1"}},
		{Name: "go", SupportedVersions: []string{"1.17.5"}},
	}

	context.JSON(http.StatusOK, supportedLanguages)
}

func HandleInDockerCode(context *gin.Context){
	stdout, err := exec.Command("docker", "exec", "random", "node", "/volume/loop.js").Output()
	if err != nil{
		context.JSON(http.StatusBadRequest, gin.H{"message": "failure"})
	}

	context.JSON(http.StatusOK, gin.H{
		"stdout": strings.Split(string(stdout), "\n"),
	})
}

func HandleInSystemCode(context *gin.Context){
	stdout, err := exec.Command("/opt/jdk-17/bin/java", "languages/runner.java").Output()
	if err != nil{
		context.JSON(http.StatusBadRequest, gin.H{
			"stdout": strings.Split(err.Error(), "\n"),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"stdout": strings.Split(string(stdout), "\n"),
	})
}

func HandleCodeUpload(context *gin.Context){
	userCode := ExecutableCode{}
	err := context.BindJSON(&userCode)

	if err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
	}

	fmt.Println(userCode)
}
