package pkg

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

const (
	maxOutputBufferCapacity = "65332" // Max output size 64kb
)

var (
	user = 1
	languageService = &SimpleLanguageService{}
	fileService = &SimpleFileService{}
)

type CodeRunnerService interface {
	RunCode() []string
}

type SimpleCodeRunnerService struct {}

type ExecutableCode struct {
	Language string `json:"language" binding:"required"`
	Code []string `json:"code" binding:"required"`
}

type Response struct {
	Status int
	Output []string
}

func (c *SimpleCodeRunnerService) RunCode(context *gin.Context) {
	body:= ExecutableCode{}

	if err := context.BindJSON(&body); err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
		return
	}

	language := languageService.GetLanguageByName(body.Language)
	currentUser := fmt.Sprintf("user%d", user)
	tempFolderName := uuid.New().String()

	fileService.CreateTemporaryFolder(currentUser, tempFolderName)
	filename := fileService.CreateTemporaryFile(currentUser, tempFolderName, language.Extension)
	fileService.WriteCodeToFile(filename, body.Code)
	output := executeFile(currentUser, filename, language)

	context.JSON(http.StatusOK, output)
	fileService.DeleteTemporaryFolder(currentUser, tempFolderName)

	if err := cleanUpProcesses(currentUser); err != nil {
		log.Fatalln("could not kill user processes")
	}
}

func executeFile(currentUser, file string, lang SupportedLanguage) Response {
	script := fmt.Sprintf("/borealys/languages/%s/run.sh", strings.ToLower(lang.Name))
	run := exec.Command("runuser", "-u", currentUser, "--", "/bin/bash", script, file)
	head := exec.Command("head", "--bytes", maxOutputBufferCapacity)

	errBuf := bytes.Buffer{}
	run.Stderr = &errBuf

	head.Stdin, _ = run.StdoutPipe()
	headOutput := bytes.Buffer{}
	head.Stdout = &headOutput

	_ = run.Start()
	_ = head.Start()
	_ = run.Wait()
	_ = head.Wait()

	var result string
	var status int
	if headOutput.Len() > 0 {
		result = headOutput.String()
		status = http.StatusOK
	}else{
		result = errBuf.String()
		status = http.StatusInternalServerError
	}

	return Response{
		Status: status,
		Output: strings.Split(result, "\n"),
	}
}

func cleanUpProcesses(currentUser string) error {
	return exec.Command("pkill", "-u", currentUser).Run()
}
