package pkg

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os/exec"
	"strings"
)

var currentRunner = 1

type CodeRunnerService interface {
	RunCode() []string
	GetSupportedLanguages() []SupportedLanguage
}

type SimpleCodeRunnerService struct {}

type SupportedLanguage struct {
	Name string
	SupportedVersions []string
	Extension string
	BinaryDiffers bool
	BinaryHints []string
}

type SupportedLanguageDTO struct {
	Name string
	Version string
	SupportedVersions []string
}

type ExecutableCode struct {
	Language string `json:"language" binding:"required"`
	Version string `json:"version" binding:"required"`
	Code []string `json:"code" binding:"required"`
}

func (c *SimpleCodeRunnerService) GetSupportedLanguages(context *gin.Context){
	supportedLanguages := []SupportedLanguage{
		{Name: "Java", SupportedVersions: []string{"17"}},
		{Name: "node", SupportedVersions: []string{"16.13.0", "17.3.0"}},
		{Name: "python", SupportedVersions: []string{"3.10.1"}},
		{Name: "go", SupportedVersions: []string{"1.17.5"}},
	}

	context.JSON(http.StatusOK, supportedLanguages)
}

func (c *SimpleCodeRunnerService) RunCode(context *gin.Context) {
	body := ExecutableCode{}
	err := context.BindJSON(&body)

	if err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
		return
	}

	binary := fmt.Sprintf("/binaries/%s/%s/bin/%s", body.Language, body.Version, body.Language)
	output, err := exec.Command(binary, "").Output()

	if err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	prettifiedOutput := strings.Split(string(output), "\n")
	context.JSON(http.StatusOK, prettifiedOutput)
}
