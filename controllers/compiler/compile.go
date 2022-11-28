package compiler

import (
	"github.com/Ankan002/compiler-api/helpers"
	execute_code "github.com/Ankan002/compiler-api/helpers/execute-code"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"time"
)

type compRequest struct {
	Code     string `json:"code" validate:"required"`
	Language string `json:"language" validate:"required,eq=js|eq=ts|eq=py|eq=go|eq=java|eq=rs|eq=kt|eq=cpp|eq=c|eq=cs"`
	StdInput string `json:"stdInput"`
}

func Compiler(c *fiber.Ctx) error {
	request := compRequest{}

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
	case "java":
		javaCompileResponse := execute_code.CompileJava("code/"+createFileResponse.FileName, request.StdInput)

		if !javaCompileResponse.Success {
			stdErr = javaCompileResponse.Error
		} else {
			stdOutput = javaCompileResponse.Output
		}
		break
	case "rs":
		rustCompileResponse := execute_code.CompileRust("code/"+createFileResponse.FileName, request.StdInput)

		if !rustCompileResponse.Success {
			stdErr = rustCompileResponse.Error
		} else {
			stdOutput = rustCompileResponse.Output
		}
		break
	case "kt":
		kotlinCompileResponse := execute_code.CompileKotlin("code/"+createFileResponse.FileName, request.StdInput)

		if !kotlinCompileResponse.Success {
			stdErr = kotlinCompileResponse.Error
		} else {
			stdOutput = kotlinCompileResponse.Output
		}
		break
	case "cpp":
		cppCompileResponse := execute_code.CompileCpp("code/"+createFileResponse.FileName, request.StdInput)

		if !cppCompileResponse.Success {
			stdErr = cppCompileResponse.Error
		} else {
			stdOutput = cppCompileResponse.Output
		}
		break
	case "c":
		cCompileResponse := execute_code.CompileC("code/"+createFileResponse.FileName, request.StdInput)

		if !cCompileResponse.Success {
			stdErr = cCompileResponse.Error
		} else {
			stdOutput = cCompileResponse.Output
		}
		break
	case "cs":
		cSharpCompileResponse := execute_code.CompileCSharp("code/"+createFileResponse.FileName, request.StdInput)

		if !cSharpCompileResponse.Success {
			stdErr = cSharpCompileResponse.Error
		} else {
			stdOutput = cSharpCompileResponse.Output
		}
		break
	default:
		stdErr = "Please provide us with a valid language"
	}

	helpers.DeleteFile("code/" + createFileResponse.FileName)

	if stdErr != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success":   false,
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
