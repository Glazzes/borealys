package pkg

import (
	"encoding/json"
	"fmt"
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

type SupportedLanguageDTO struct {
	Name string `json:"name" binding:"required"`
	Timeout int `json:"timeout" binding:"required"`
	Versions []string `json:"versions" binding:"required"`
}

type SupportedLanguage struct {
	Name string `json:"name" binding:"required"`
	Version string `json:"version" binding:"required"`
	Extension string `json:"extension" binding:"required"`
	Timeout int `json:"timeout" binding:"required"`
}

type LanguageService interface {
	SaveAll()
	GetAll(c *gin.Context)
	SaveRunnableLanguages()
	GetRunnableLanguageByKey(key string) SupportedLanguage
}

type SimpleLanguageService struct {}

func (context *SimpleLanguageService) SaveAll() {
	languages := make(map[string]SupportedLanguageDTO, 0)

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

			value, exists := languages[language.Name]
			if exists {
				value.Versions = append(value.Versions, language.Version)
			}else {
				languages[language.Name] = SupportedLanguageDTO{
					Name: language.Name,
					Timeout: language.Timeout,
					Versions: []string{language.Version},
				}
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	saveToRedis(languages)
}

func saveToRedis(languages map[string]SupportedLanguageDTO){
	for key, value := range languages {
		serialized, err := json.Marshal(value)

		if err != nil {
			log.Fatal(err)
		}else{
			currentKey := fmt.Sprintf("%s-%s", lang, key)
			config.RedisClient.Set(currentKey, string(serialized), -1)
		}
	}
}

func (context *SimpleLanguageService) GetAll(ginC *gin.Context) {
	result, err := config.RedisClient.Keys(lang + "-*").Result()
	checkNilError(err)

	store := make([]SupportedLanguageDTO, 0)
	for _, key := range result {
		lang, err := config.RedisClient.Get(key).Result()
		checkNilError(err)

		serialized := []byte(lang)
		holder := SupportedLanguageDTO{}
		err = json.Unmarshal(serialized, &holder)

		checkNilError(err)

		store = append(store, holder)
	}

	ginC.JSON(http.StatusOK, store)
}

func (context *SimpleLanguageService) SaveRunnableLanguages() {
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

			key := fmt.Sprintf("%s-%s", language.Name, language.Version)
			config.RedisClient.Set(key, string(fileBytes), -1)
		}
		return nil
	})

	checkNilError(err)
}

func (context *SimpleLanguageService) GetRunnableLanguageByKey(key string) SupportedLanguage {
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