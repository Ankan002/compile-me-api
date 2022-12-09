package download_code

import (
	"github.com/Ankan002/compiler-api/helpers"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type downloadCodeFileRequest struct {
	Code     string `json:"code" validate:"required"`
	Language string `json:"language" validate:"required,eq=js|eq=ts|eq=py|eq=go|eq=java|eq=rs|eq=kt|eq=cpp|eq=c|eq=cs"`
}

func DownloadCodeFile(c *fiber.Ctx) error {
	requestBody := downloadCodeFileRequest{}

	if bodyParserError := c.BodyParser(&requestBody); bodyParserError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Parsing Error Occurred",
		})
	}

	if validatorError := validator.New().Struct(requestBody); validatorError != nil {
		var validationError string

		if validatorError.Error() == "Key: 'CompRequest.Language' Error:Field validation for 'Language' failed on the 'eq=js|eq=ts|eq=py|eq=go|eq=java|eq=rs|eq=kt|eq=cpp|eq=c|eq=cs' tag" {
			validationError = "Provide us with a valid language extension..."
		} else {
			validationError = validatorError.Error()
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   validationError,
		})
	}

	createFileResponse := helpers.CreateFile(requestBody.Code, requestBody.Language)

	if !createFileResponse.Success {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   createFileResponse.Error,
		})
	}

	fileDownloadError := c.Download("code/"+createFileResponse.FileName, createFileResponse.FileName)

	if fileDownloadError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   fileDownloadError.Error(),
		})
	}

	helpers.DeleteFile("code/" + createFileResponse.FileName)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}
