package compiler

import (
	"github.com/Ankan002/compiler-api/helpers"
	execute_code "github.com/Ankan002/compiler-api/helpers/execute-code"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"time"
)

type CompRequest struct {
	Code     string `json:"code" validate:"required"`
	Language string `json:"language" validate:"required,eq=js|eq=ts|eq=py|eq=go"`
	StdInput string `json:"stdInput"`
}

func Compiler(c *fiber.Ctx) error {
	request := CompRequest{}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Parsing Error Occurred...",
		})
	}

	if err := validator.New().Struct(request); err != nil {
		var validationError string

		if err.Error() == "Key: 'CompRequest.Language' Error:Field validation for 'Language' failed on the 'eq=js|eq=ts|eq=py|eq=go' tag" {
			validationError = "Provide us with a valid language extension..."
		} else {
			validationError = err.Error()
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   validationError,
		})
	}

	createFileResponse := helpers.CreateFile(request.Code, request.Language)

	if !createFileResponse.Success {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   createFileResponse.Error,
		})
	}

	var stdOutput string
	var stdErr string

	switch request.Language {
	case "js":
		jsCompileResponse := execute_code.CompileJavascript("code/"+createFileResponse.FileName, request.StdInput)

		if !jsCompileResponse.Success {
			stdErr = jsCompileResponse.Error
		} else {
			stdOutput = jsCompileResponse.Output
		}
		break
	case "ts":
		tsCompileResponse := execute_code.CompileTypescript("code/"+createFileResponse.FileName, request.StdInput)

		if !tsCompileResponse.Success {
			stdErr = tsCompileResponse.Error
		} else {
			stdOutput = tsCompileResponse.Output
		}
		break
	case "py":
		pyCompileResponse := execute_code.CompilePython("code/"+createFileResponse.FileName, request.StdInput)

		if !pyCompileResponse.Success {
			stdErr = pyCompileResponse.Error
		} else {
			stdOutput = pyCompileResponse.Output
		}
		break
	case "go":
		goCompileResponse := execute_code.CompileGo("code/"+createFileResponse.FileName, request.StdInput)

		if !goCompileResponse.Success {
			stdErr = goCompileResponse.Error
		} else {
			stdOutput = goCompileResponse.Output
		}
		break
	default:
		stdErr = "Please provide us with a valid language"
	}

	helpers.DeleteFile("code/" + createFileResponse.FileName)

	if stdErr != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success":   true,
			"error":     stdErr,
			"timestamp": time.Now(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success":   true,
		"output":    stdOutput,
		"timestamp": time.Now(),
	})
}
