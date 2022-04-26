package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pemistahl/lingua-go"
)

var detector lingua.LanguageDetector

type TextRequest struct {
	Text string `json:"text"`
}

type LanguageResult struct {
	Name     string  `json:"name"`
	Iso639_1 string  `json:"iso639_1"`
	Score    float64 `json:"score"`
}

type LanguageResponse struct {
	Results []LanguageResult `json:"results"`
}

func main() {
	engine := gin.Default()
	engine.POST("/api/language", postLanguage)

	languages := []lingua.Language{
		lingua.German,
		lingua.Spanish,
		lingua.Italian,
		lingua.French,
		lingua.English,
	}
	detector = lingua.
		NewLanguageDetectorBuilder().
		FromLanguages(languages...).
		WithPreloadedLanguageModels().
		Build()

	engine.Run(get_port())
}

func postLanguage(c *gin.Context) {
	var request TextRequest

	if err := c.BindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, "Bad request, could not get text from request")
		return
	}

	if languages := detector.ComputeLanguageConfidenceValues(request.Text); languages != nil {
		var response LanguageResponse

		for _, elem := range languages {
			res := LanguageResult{
				Name:     elem.Language().String(),
				Iso639_1: elem.Language().IsoCode639_1().String(),
				Score:    elem.Value(),
			}
			response.Results = append(response.Results, res)
		}
		c.IndentedJSON(http.StatusOK, response)
		return
	}

	c.String(http.StatusNotFound, "Language not detected")
}

func get_port() string {
	// if running on Azure, get the correct port
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		return ":" + val
	}
	return ":8080"
}
