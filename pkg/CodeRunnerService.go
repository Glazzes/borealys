package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/glazzes/borealys/pkg/config"
	"github.com/google/uuid"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var (
	user = 1
	cacheService = &SimpleCacheService{}
	languageService = &SimpleLanguageService{}
)

type CodeRunnerService interface {
	RunCode() []string
}

type SimpleCodeRunnerService struct {}

type SupportedLanguage struct {
	Name string `json:"name" binding:"required"`
	Version string `json:"version" binding:"required"`
	Extension string `json:"extension" binding:"required"`
	Binary string `json:"binary" binding:"required"`
}

type SupportedLanguageDTO struct {
	Name string `json:"name" binding:"required"`
	Versions []string `json:"versions" binding:"required"`
}

type ExecutableCode struct {
	Language string `json:"language" binding:"required"`
	Version string `json:"version" binding:"required"`
	Code []string `json:"code" binding:"required"`
}

func (c *SimpleCodeRunnerService) RunCode(context *gin.Context) {
	body:= ExecutableCode{}

	if err := context.BindJSON(&body); err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
		return
	}

	find, err := config.RedisClient.Get(body.Language + "-" + body.Version).Result()
	checkNilError(err)

	language := SupportedLanguage{}
	checkNilError(json.Unmarshal([]byte(find), &lang))


	currentUser := fmt.Sprintf("user%d", user)
	tempFolderName := uuid.New().String()

	checkNilError(createTempFolder(currentUser, tempFolderName))
	filename, err := createTempFile(currentUser, tempFolderName, language.Extension)
	checkNilError(err)

	writeCodeToTempFile(filename, body.Code)
	executeFile(currentUser, language.Binary, filename)
	checkNilError(cleanUpFiles(currentUser, tempFolderName))
	checkNilError(cleanUpProcesses(currentUser))
}

func createTempFolder(currentUser, tempFolderName string) error {
	command := fmt.Sprintf("'mkdir /tmp/%s/%s'", currentUser, tempFolderName)
	return exec.Command("runuser", "-l", currentUser, "-c", command).Run()
}

func createTempFile(currentUser, tempFolderName, extension string) (string ,error) {
	filename := fmt.Sprintf("/tmp/%s/%s/code%s", currentUser, tempFolderName, extension)
	command := fmt.Sprintf("'touch %s'", filename)
	return filename, exec.Command("runuser", "-l", currentUser, "-c", command).Run()
}

func writeCodeToTempFile(filename string, code []string){
	file, err := os.OpenFile(filename, os.O_APPEND, 0644)
	checkNilError(err)
	defer file.Close()

	for _, line := range code {
		file.WriteString(line + "\n")
	}
}

func executeFile(currentUser, binary, file string) []string {
	command := fmt.Sprintf("'%s %s'", binary, file)
	output, _ := exec.Command("runuser", "-l", currentUser, "-c", command).Output()
	return strings.Split(string(output), "\n")
}

func cleanUpFiles(currentUser, tempFolder string) error {
	folderName := fmt.Sprintf("/tmp/%s/%s", currentUser, tempFolder)
	return exec.Command("rm", "-rf", folderName).Run()
}

func cleanUpProcesses(currentUser string) error {
	return exec.Command("pkill", "-u", currentUser).Run()
}
