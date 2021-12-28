package pkg

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"os/exec"
	"strings"
)

var (
	user = 1
	cacheService = &SimpleCacheService{}
	languageService = &SimpleLanguageService{}
	fileService = &SimpleFileService{}
)

type CodeRunnerService interface {
	RunCode() []string
}

type SimpleCodeRunnerService struct {}

type ExecutableCode struct {
	Language string `json:"language" binding:"required"`
	Version string `json:"version" binding:"required"`
	Code []string `json:"code" binding:"required"`
}

type Response struct {
	Trimmed bool
	Output []string
}

func (c *SimpleCodeRunnerService) RunCode(context *gin.Context) {
	body:= ExecutableCode{}

	if err := context.BindJSON(&body); err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
		return
	}

	language := languageService.GetRunnableLanguageByKey(body.Language + "-" + body.Version)
	currentUser := fmt.Sprintf("user%d", user)
	tempFolderName := uuid.New().String()

	fileService.CreateTemporaryFolder(currentUser, tempFolderName)
	filename := fileService.CreateTemporaryFile(currentUser, tempFolderName, language.Extension)
	fileService.WriteCodeToFile(filename, body.Code)
	output := executeFile(currentUser, filename, language)

	context.JSON(http.StatusOK, output)
	fileService.DeleteTemporaryFolder(currentUser, tempFolderName)

	/*
	if err = cleanUpProcesses(currentUser); err != nil {
		log.Fatalln("could not kill user processes")
	}
	 */
}

func executeFile(currentUser, file string, lang SupportedLanguage) Response {
	script := fmt.Sprintf("/borealys/languages/%s/%s/run.sh", strings.ToLower(lang.Name), lang.Version)

	output, _ := exec.Command(
		"runuser",
		"-u",
		currentUser,
		"--",
		"/bin/bash",
		script,
		file).Output()

	strOutput := strings.Split(string(output[:48000]), "\n")

	return Response{
		Trimmed: len(strOutput) > 200,
		Output: strOutput,
	}
}

func cleanUpProcesses(currentUser string) error {
	return exec.Command("pkill", "-u", currentUser).Run()
}
