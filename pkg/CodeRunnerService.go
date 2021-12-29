package pkg

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os/exec"
	"strings"
)

const (
	maxOutputByteSliceLength = 65332 // Max output size 64kb
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
	Status int
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

	cmd := exec.Command(
		"runuser",
		"-u",
		currentUser,
		"--",
		"/bin/bash",
		script,
		file)

	stdoutBuf := bytes.Buffer{}
	stderrBuf := bytes.Buffer{}

	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	cmd.Run()
	var currentOutput []byte
	var status int
	if len(stdoutBuf.Bytes()) > 0 {
		currentOutput = stdoutBuf.Bytes()
		status = http.StatusOK
	}else{
		currentOutput = stderrBuf.Bytes()
		status = http.StatusBadRequest
	}

	strOutput := getOutput(currentOutput)

	return Response{
		Status: status,
		Trimmed: len(currentOutput) > maxOutputByteSliceLength,
		Output: strOutput[:len(strOutput) - 1],
	}
}

func getOutput(byteOutPut []byte) []string{
	byteReader := bytes.NewReader(byteOutPut)
	limitReader := io.LimitReader(byteReader, maxOutputByteSliceLength)
	dest := make([]byte, maxOutputByteSliceLength)

	limitReader.Read(dest)
	return strings.Split(string(dest[:len(dest) - 1]), "\n")
}

func cleanUpProcesses(currentUser string) error {
	return exec.Command("pkill", "-u", currentUser).Run()
}
