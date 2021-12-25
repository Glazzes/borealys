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
	baseCapacity = 3
)

type LanguageService interface {
	SaveAll()
	GetAll(c *gin.Context)
	GetByKey(key string)
	Save(key string, language SupportedLanguage)
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
	keys := fmt.Sprintf("%s-*", lang)
	result, err := config.RedisClient.Keys(keys).Result()
	checkNilError(err)

	store := make([]SupportedLanguageDTO, baseCapacity)
	for _, r := range result {
		serialized := []byte(r)
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

func (c *SimpleLanguageService) Save(key string, language SupportedLanguage) {
	jsonLanguage, err := json.Marshal(language)
	if err != nil {
		log.Fatal(err)
	}

	status := config.RedisClient.Set(key, jsonLanguage, -1)
	if status.Err() != nil{
		log.Fatal(status.Err())
	}
}

func checkNilError(err error){
	if err != nil {
		log.Fatal(err)
	}
}