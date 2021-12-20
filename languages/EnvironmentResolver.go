package main

import (
	"fmt"
	"os"
)

func main(){
	/*
	files, err := ioutil.ReadDir("/home/glaze/")
	if os.IsNotExist(err){
		fmt.Println("This directory does not exists")
		log.Fatal(err)
	}

	for _, f := range files{
		fmt.Println(f.Name(), f.IsDir())
	}

	 */
	createUserFolders()
}

func createUserFolders(){
	for i := 1; i <= 100; i++ {
		userFolder := fmt.Sprintf("/tmp/user%d", i)
		os.Mkdir(userFolder, 0700)
	}
}

func findSetupFiles()  {

}