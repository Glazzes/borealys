package pkg

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

type InitializerService interface {
	CreateExecutorUsers()
	SetUpBinaries()
	CreateRunnersGroup()
}

type SimpleInitializerService struct {}

func New() *SimpleInitializerService {
	return &SimpleInitializerService{}
}

const (
	binaryName string = "binary"
	binaryVersion string = "version"
)

var (
	infoLogger *log.Logger
	initializer = New()
	regex = regexp.MustCompile("(?P<binary>\\w+)/(?P<version>\\d{1,2}(\\.\\d{1,2})?(\\.\\d{1,2})?)/?\\w+\\.sh/?$")
)

func init()  {
	infoLogger = log.New(os.Stdout, "Info: ", log.LstdFlags | log.Lshortfile)
	initializer = New()
}

func (context *SimpleInitializerService) CreateRunnersGroup() {
	infoLogger.Println("Creating runners group")
	if err := exec.Command("/bin/bash", "/borealys/languages/scripts/create-group.sh").Run(); err != nil {
		log.Fatal(err)
	}

	infoLogger.Println("Created runners group successfully")
}

func (context *SimpleInitializerService) CreateExecutorUsers(){
	infoLogger.Println("Creating executor users")
	if err := exec.Command("/bin/bash", "/borealys/languages/scripts/create-users.sh").Run(); err != nil {
		log.Fatal(err)
	}

	infoLogger.Println("Created executor users successfully")
}

func (context *SimpleInitializerService)SetUpBinaries(){
	infoLogger.Print("Setting up binaries")
	err := filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, "/16.3.1/setup.sh"){
			info, err := GetBinaryInfoFromPath(path)
			if err == nil {
				infoLogger.Printf("Downloading %s %s binaries", info[binaryName], info[binaryVersion])
			}

			DownloadBinary(path, info[binaryName], info[binaryVersion])
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
		log.Fatal(err)
	}

	message := fmt.Sprintf("%s %s binaries downloaded successfully!!", binary, version)
	infoLogger.Println(strings.Title(message))
}

func checkNilErrorFatal(err error){
	if err != nil{
		log.Fatal(err)
	}
}