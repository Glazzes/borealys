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
	CreateRunnersGroup()
	CreateExecutorUsers()
	SetUpBinaries()
	HardenScripts()
}

type SimpleInitializerService struct {}

const (
	binaryName string = "binary"
)

var (
	infoLogger *log.Logger
	regex = regexp.MustCompile("(?P<binary>\\w+)/\\w+\\.sh/?$")
)

func init()  {
	infoLogger = log.New(os.Stdout, "Info: ", log.LstdFlags | log.Lshortfile)
}

func (context *SimpleInitializerService) CreateRunnersGroup() {
	infoLogger.Println("Creating runners group...")
	if err := exec.Command("/bin/bash", "/borealys/scripts/create-group.sh").Run(); err != nil {
		if strings.Contains(err.Error(), "already exists"){
			infoLogger.Println("Runners group already exists")
		}else{
			log.Fatal(err)
		}
	}

	infoLogger.Println("Created runners group successfully!!!")
}

func (context *SimpleInitializerService) CreateExecutorUsers(){
	infoLogger.Println("Creating executor users...")
	if err := exec.Command("/bin/bash", "/borealys/scripts/create-users.sh").Run(); err != nil {
		log.Fatal(err)
	}

	infoLogger.Println("Created executor users successfully!!!")
}

func (context *SimpleInitializerService)SetUpBinaries(){
	infoLogger.Print("Setting up binaries")
	_ = filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, "setup.sh"){
			info, err := GetBinaryInfoFromPath(path)
			if err == nil {
				infoLogger.Printf("Downloading %s binaries", info[binaryName])
			}

			DownloadBinary(path, info[binaryName])
		}
		return nil
	})

	//checkNilErrorFatal(err)
	infoLogger.Println("Binaries and environment have been downloaded and exported successfully")
}

func GetBinaryInfoFromPath(filePath string) (map[string]string, error) {
	info := make(map[string]string)
	matches := regex.FindStringSubmatch(filePath)

	if len(matches) > 0 {
		binaryIndex := regex.SubexpIndex(binaryName)
		info[binaryName] = matches[binaryIndex]

		return info, nil
	}

	return nil, fmt.Errorf("no matches were found on the given path")
}

func DownloadBinary(path, binary string){
	err := exec.Command("bash", path).Run()
	if err != nil{
		log.Fatal(err)
	}

	message := fmt.Sprintf("%s binaries downloaded successfully!!", binary)
	infoLogger.Println(strings.Title(message))
}

func (context *SimpleInitializerService) HardenScripts(){
	infoLogger.Println("Hardening scripts and binaries...")
	if err := exec.Command("/bin/bash", "/borealys/scripts/harden-scripts.sh").Run(); err != nil {
		log.Fatal(err)
	}

	infoLogger.Println("Hardened scripts and binaries successfully!!!")
}

func checkNilErrorFatal(err error){
	if err != nil{
		log.Fatal(err)
	}
}