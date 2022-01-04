package pkg

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"os"
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

	// clean up
	fileService.DeleteTemporaryFolder(currentUser, tempFolderName)
	_ = cleanUpProcesses(currentUser)

	// running rm -rf --no-preserve-root deletes user folder, so i needs to be restored
	restoreUserFolderIfDeleted(currentUser)

	// assign next user to run code
	if user >= 3{
		user = 1
	}else{
		user++
	}
}

func executeFile(currentUser, file string, lang SupportedLanguage) Response {
	script := fmt.Sprintf("/borealys/languages/%s/run.sh", strings.ToLower(lang.Name))

	run := exec.Command("/bin/bash", script, currentUser, file)
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

	if headOutput.Len() > 0 {
		result = headOutput.String()
	}else if headOutput.Len() == 0 && errBuf.Len() == 0 {
		result = headOutput.String()
	}else{
		result = errBuf.String()
	}

	return Response{
		Output: strings.Split(result, "\n"),
	}
}

func cleanUpProcesses(currentUser string) error {
	return exec.Command("pkill", "-9", "-u", currentUser).Run()
}

func restoreUserFolderIfDeleted(currentUser string){
	userFolder := "/tmp/" + currentUser
	if _, err := ioutil.ReadDir(userFolder); err != nil{
		if os.IsNotExist(err){
			_ = exec.Command("runuser", "-u", currentUser, "--", "mkdir", userFolder).Run()
		}
	}
}
