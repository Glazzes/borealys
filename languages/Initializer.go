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
	regex = regexp.MustCompile("(?P<binary>\\w+)/(?P<version>\\d{1,2}(\\.\\d{1,2})?(\\.\\d{1,2})?)/?\\w+\\.sh/?$")
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
		if strings.HasSuffix(path, "setup.sh"){
			info, err := GetBinaryInfoFromPath(path)
			if err == nil {
				infoLogger.Printf("Downloading %s %s binaries", info[binaryName], info[binaryVersion])
			}

			DownloadBinary(path, info[binaryName], info[binaryVersion])
		}
		return nil
	})

	err = filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, "environment.sh") {
			info, err := GetBinaryInfoFromPath(path)
			if err != nil {
				return fmt.Errorf("could not expor environment variable")
			}

			ExportEnvironmentVariables(path, info[binaryName], info[binaryVersion])
			infoLogger.Printf("Exported variable for %s %s", info[binaryName], info[binaryVersion])
		}

		return nil
	})

	checkNilErrorFatal(err)

	infoLogger.Println("Binaries and environment have been downloaded and exported successfully")
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

func DownloadBinary(path, binary, version string){
	err := exec.Command("bash", path).Run()
	if err != nil{
		log.Fatalf("could not download binary %s %s", binary, version)
	}

	message := fmt.Sprintf("%s %s binaries downloaded successfully!!", binary, version)
	infoLogger.Println(strings.Title(message))
}

func ExportEnvironmentVariables(path, binary, version string){
	command := fmt.Sprintf("source %s", path)
	err := exec.Command("bash", "-c", command).Run()
	if err != nil{
		log.Fatal(err)
	}

	infoLogger.Printf("Set up environment variable for %s %s", binary, version)
}

func checkNilErrorFatal(err error){
	if err != nil{
		log.Fatal(err)
	}
}