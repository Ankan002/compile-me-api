package languages

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"io/ioutil"
)

type supportedLanguagesObject struct {
	Language              string `json:"language"`
	CompileQueryParameter string `json:"compile_query_parameter"`
}

type supportedLanguagesResponse struct {
	Languages []supportedLanguagesObject `json:"languages"`
}

func GetLanguages(c *fiber.Ctx) error {
	var currentSupportedLanguages supportedLanguagesResponse

	currentSupportedLanguagesBytes, _ := ioutil.ReadFile("constants/supported-languages.json")

	if unmarshalError := json.Unmarshal(currentSupportedLanguagesBytes, &currentSupportedLanguages); unmarshalError != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   unmarshalError.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"data":    currentSupportedLanguages,
	})
}
