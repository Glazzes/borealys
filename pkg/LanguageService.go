package pkg

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/glazzes/borealys/pkg/config"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

var (
	lang = "lang"
)

type SupportedLanguage struct {
	Name string `json:"name" binding:"required"`
	Version string `json:"version" binding:"required"`
	Extension string `json:"extension" binding:"required"`
	Timeout int `json:"timeout" binding:"required"`
}

type LanguageService interface {
	SaveAll()
	GetAll(c *gin.Context)
	GetLanguageByName(key string) SupportedLanguage
}

type SimpleLanguageService struct {}

func (context *SimpleLanguageService) SaveAll() {
	err := filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, "metadata.json") {
			fileBytes, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			language := SupportedLanguage{}
			if err = json.Unmarshal(fileBytes, &language); err != nil {
				return err
			}

			config.RedisClient.Set(language.Name, string(fileBytes), -1)
		}
		return nil
	})

	checkNilError(err)
}

func (context *SimpleLanguageService) GetAll(ginC *gin.Context) {
	result, err := config.RedisClient.Keys("*").Result()
	checkNilError(err)

	store := make([]SupportedLanguage, 0)
	for _, key := range result {
		lang, err := config.RedisClient.Get(key).Result()
		checkNilError(err)

		serialized := []byte(lang)
		holder := SupportedLanguage{}
		err = json.Unmarshal(serialized, &holder)

		checkNilError(err)

		store = append(store, holder)
	}

	ginC.JSON(http.StatusOK, store)
}

func (context *SimpleLanguageService) GetLanguageByName(key string) SupportedLanguage {
	find, err := config.RedisClient.Get(key).Result()
	if err != nil {
		log.Fatal("Could not find language with key " + key)
	}

	language := SupportedLanguage{}
	if err = json.Unmarshal([]byte(find), &language); err != nil {
		log.Fatal(err)
	}

	return language
}

func checkNilError(err error){
	if err != nil {
		log.Fatal(err)
	}
}