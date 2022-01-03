package pkg

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type FileService interface {
	CreateTemporaryFolder(user string)
	CreateTemporaryFile(user, folderName string) string
	DeleteTemporaryFolder(user, folderName string)
	WriteCodeToFile(filename string, code []string)
}

type SimpleFileService struct {}

func (context *SimpleFileService) CreateTemporaryFolder(user, folderName string) {
	temporaryFolderName := fmt.Sprintf("/tmp/%s/%s", user, folderName)
	if err := exec.Command("runuser", "-u", user, "--", "mkdir", temporaryFolderName).Run(); err != nil {
		infoLogger.Println("could not create temp folder")
		log.Fatalln(err)
	}
}

func (context *SimpleFileService) CreateTemporaryFile(user, folderName, extension string) string {
	filename := fmt.Sprintf("/tmp/%s/%s/code%s", user, folderName, extension)

	if err := exec.Command(
		"runuser",
		"-u",
		user,
		"--",
		"touch",
		filename).Run(); err != nil {
		infoLogger.Println("Could not create temporary file")
		log.Fatal(err)
	}

	return filename
}

func (context *SimpleFileService) DeleteTemporaryFolder(user, folderName string)  {
	temporaryFolderName := fmt.Sprintf("/tmp/%s/%s", user, folderName)
	err := exec.Command("rm", "-rf", temporaryFolderName).Run()
	if err != nil{
		infoLogger.Println("could not delete temp folder")
		log.Fatal(err)
	}
}

func (context *SimpleFileService) WriteCodeToFile(filename string, code []string)  {
	file, err := os.OpenFile(filename, os.O_APPEND | os.O_RDWR, 0644)
	if err != nil {
		infoLogger.Println("could not write code to file")
		log.Fatal(err)
	}

	defer file.Close()

	for _, line := range code {
		if _ , err := file.WriteString(line + "\n"); err != nil {
			log.Fatalln("could not write line to file status")
		}
	}
}
