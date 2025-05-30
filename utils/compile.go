package utils

import (
	"github.com/Ankan002/compiler-api/helpers"
	execute_code "github.com/Ankan002/compiler-api/helpers/execute-code"
)

type CompilationArgs struct {
	Code     string
	Language string
	StdInput string
}

type CompilationResponse struct {
	StdErr    string
	StdOutput string
	Error     string
}

func Compile(args CompilationArgs) CompilationResponse {
	createFileResponse := helpers.CreateFile(args.Code, args.Language)

	if !createFileResponse.Success {
		return CompilationResponse{
			Error: createFileResponse.Error,
		}
	}

	var stdOutput string
	var stdErr string

	switch args.Language {
	case "js":
		jsCompileResponse := execute_code.CompileJavascript("code/"+createFileResponse.FileName, args.StdInput)

		if !jsCompileResponse.Success {
			stdErr = jsCompileResponse.Error
		} else {
			stdOutput = jsCompileResponse.Output
		}
		break
	case "ts":
		tsCompileResponse := execute_code.CompileTypescript("code/"+createFileResponse.FileName, args.StdInput)

		if !tsCompileResponse.Success {
			stdErr = tsCompileResponse.Error
		} else {
			stdOutput = tsCompileResponse.Output
		}
		break
	case "py":
		pyCompileResponse := execute_code.CompilePython("code/"+createFileResponse.FileName, args.StdInput)

		if !pyCompileResponse.Success {
			stdErr = pyCompileResponse.Error
		} else {
			stdOutput = pyCompileResponse.Output
		}
		break
	case "go":
		goCompileResponse := execute_code.CompileGo("code/"+createFileResponse.FileName, args.StdInput)

		if !goCompileResponse.Success {
			stdErr = goCompileResponse.Error
		} else {
			stdOutput = goCompileResponse.Output
		}
		break
	case "java":
		javaCompileResponse := execute_code.CompileJava("code/"+createFileResponse.FileName, args.StdInput)

		if !javaCompileResponse.Success {
			stdErr = javaCompileResponse.Error
		} else {
			stdOutput = javaCompileResponse.Output
		}
		break
	case "rs":
		rustCompileResponse := execute_code.CompileRust("code/"+createFileResponse.FileName, args.StdInput)

		if !rustCompileResponse.Success {
			stdErr = rustCompileResponse.Error
		} else {
			stdOutput = rustCompileResponse.Output
		}
		break
	case "kt":
		kotlinCompileResponse := execute_code.CompileKotlin("code/"+createFileResponse.FileName, args.StdInput)

		if !kotlinCompileResponse.Success {
			stdErr = kotlinCompileResponse.Error
		} else {
			stdOutput = kotlinCompileResponse.Output
		}
		break
	case "cpp":
		cppCompileResponse := execute_code.CompileCpp("code/"+createFileResponse.FileName, args.StdInput)

		if !cppCompileResponse.Success {
			stdErr = cppCompileResponse.Error
		} else {
			stdOutput = cppCompileResponse.Output
		}
		break
	case "c":
		cCompileResponse := execute_code.CompileC("code/"+createFileResponse.FileName, args.StdInput)

		if !cCompileResponse.Success {
			stdErr = cCompileResponse.Error
		} else {
			stdOutput = cCompileResponse.Output
		}
		break
	case "cs":
		cSharpCompileResponse := execute_code.CompileCSharp("code/"+createFileResponse.FileName, args.StdInput)

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

	return CompilationResponse{
		StdErr:    stdErr,
		StdOutput: stdOutput,
	}
}
