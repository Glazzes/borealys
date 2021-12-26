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

type LanguageService interface {
	SaveAll()
	SaveRunnableLanguages()
	GetAll(c *gin.Context)
	GetByKey(key string)
}

type SimpleLanguageService struct {}

func (c *SimpleLanguageService) SaveAll() {
	languages := make(map[string]SupportedLanguageDTO, 3)

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

func (c *SimpleLanguageService) GetAll(context *gin.Context) {
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

	context.JSON(http.StatusOK, store)
}

func (c *SimpleLanguageService) GetByName(key string) SupportedLanguageDTO {
	savedLang := config.RedisClient.Get(key)
	bytes, err := savedLang.Bytes()

	checkNilError(err)

	supportedLanguage := SupportedLanguage{}
	if err = json.Unmarshal(bytes, &supportedLanguage); err != nil {
		log.Fatal(err)
	}

	return SupportedLanguageDTO{
		Name: supportedLanguage.Name,
		Versions:[]string{},
	}
}

func (c *SimpleLanguageService) SaveRunnableLanguages() {
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

func checkNilError(err error){
	if err != nil {
		log.Fatal(err)
	}
}