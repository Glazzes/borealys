package pkg

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"os/exec"
	"strings"
)

var user = 1

type CodeRunnerService interface {
	RunCode() []string
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

func (c *SimpleCodeRunnerService) RunCode(context *gin.Context) {
	body := ExecutableCode{}
	err := context.BindJSON(&body)

	if err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
		return
	}

	currentUser := fmt.Sprintf("user%d", user)
	tempFolderName := uuid.New().String()
	tempFilename := uuid.New().String()

	createTempFolder(currentUser, tempFolderName)
}

func createTempFolder(currentUser, tempFolderName string) error {
	command := fmt.Sprintf("'mkdir /tmp/%s/%s'", currentUser, tempFolderName)
	return exec.Command("runuser", "-l", currentUser, "-c", command).Run()
}

func createTempFile(currentUser, tempFilename, tempFolderName, extension string) error {
	filename := fmt.Sprintf("/tmp/%s/%s/%s.%s", currentUser, tempFolderName, tempFilename, extension)
	command := fmt.Sprintf("'touch %s'", filename)
	return exec.Command("runuser", "-l", currentUser, "-c", command).Run()
}

func executeFile(currentUser, binary, file string) []string {
	command := fmt.Sprintf("'%s %s'", binary, file)
	output, _ := exec.Command("runuser", "-l", currentUser, "-c", command).Output()
	return strings.Split(string(output), "\n")
}

func cleanUpFiles(currentUser string) error {
	folderName := fmt.Sprint("/tmp/%s", currentUser)
	return exec.Command("rm", "-rf", folderName).Run()
}

func cleanUpProcesses(currentUser string) error {
	return exec.Command("pkill", "-u", currentUser).Run()
}
