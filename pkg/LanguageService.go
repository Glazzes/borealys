package pkg

import "github.com/gin-gonic/gin"

type LanguageService interface {
	GetAll(c *gin.Context)
	GetByName(name string)
}

type SimpleLanguageService struct {}

func (c *SimpleLanguageService) GetAll(context *gin.Context) {

}

func (c *SimpleLanguageService) GetByName(name string) SupportedLanguageDTO {
	return SupportedLanguageDTO{}
}