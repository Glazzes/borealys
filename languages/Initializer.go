package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	binaryName string = "binary"
	binaryVersion string = "version"
)

var (
	infoLogger *log.Logger
	regex = regexp.MustCompile(".*/(?P<binary>\\w+)/(?P<version>\\d{1,2}(\\.\\d{1,2})?(\\.\\d{1,2})?)/\\w+\\.sh$")
)

func init()  {
	infoLogger = log.New(os.Stdout, "Info: ", log.LstdFlags | log.Lshortfile)
}

func main(){
	CreateCodeExecutors()
	SetUpBinaries()
}

func CreateCodeExecutors(){
	infoLogger.Println("Creating executor users")
	_, err := exec.Command("bash", "./scripts/create-users.sh").Output()
	if err != nil{
		if os.IsNotExist(err){
			log.Fatal(err)
		}
	}

	infoLogger.Println("Created executor users successfully")
}

func SetUpBinaries(){
	infoLogger.Print("Setting up binaries")
	err := filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, ".sh") && !strings.Contains(path, "scripts"){
			info, err := GetBinaryInfoFromPath(path)
			if err == nil {
				infoLogger.Printf("Downloading binary %s %s", info[binaryName], info[binaryVersion])
			}

			err = exec.Command("bash", path).Run()
			if err != nil{
				return err
			}

			infoLogger.Printf("Binary %s %s downloaded and linked successfully", info[binaryName], info[binaryVersion])
		}

		return nil
	})

	checkNilErrorFatal(err)

	infoLogger.Println("All binaries downloaded successfully!")
	infoLogger.Println("You're a ready to run some code :D !!!")
}

func GetBinaryInfoFromPath(filePath string) (map[string]string, error) {
	info := make(map[string]string)
	matches := regex.FindStringSubmatch(filePath)
	if len(matches) > 0 {
		binaryIndex := regex.SubexpIndex(binaryName)
		binaryVersionIndex := regex.SubexpIndex(binaryVersion)

		info[binaryName] = matches[binaryIndex]
		info[binaryVersion] = matches[binaryVersionIndex]

		return info, nil
	}

	return nil, fmt.Errorf("no matches were found on the given path")
}

func checkNilErrorFatal(err error){
	if err != nil{
		log.Fatal(err)
	}
}