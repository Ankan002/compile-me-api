package execute_code

import (
	"errors"
	"github.com/Ankan002/compiler-api/helpers"
	"github.com/Ankan002/compiler-api/types"
	"github.com/chebyrash/promise"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func CompileKotlin(filename string, input string) types.CompileCodeResponse {
	compilationPromise := promise.New(func(resolve func(isCreated bool), reject func(error)) {
		execCommand := exec.Command("kotlinc", filename, "-include-runtime", "-d", strings.Split(filename, ".")[0]+".jar")

		time.AfterFunc(20*time.Second, func() {
			if processKillError := execCommand.Process.Kill(); processKillError != nil {
				log.Println(processKillError.Error())
			}

			reject(errors.New("TLE"))
		})

		stdError, _ := execCommand.StderrPipe()

		if startCommandError := execCommand.Start(); startCommandError != nil {
			reject(errors.New(startCommandError.Error()))
		}

		stdErrorBytes, _ := io.ReadAll(stdError)

		if stdErrorBytes != nil && len(string(stdErrorBytes)) > 0 {
			reject(errors.New(string(stdErrorBytes)))
		}

		waitErr := execCommand.Wait()

		if waitErr != nil {
			reject(errors.New(waitErr.Error()))
		}

		resolve(true)
	})

	compilationResult, compilationError := compilationPromise.Await()

	var compilationWarningAndError string

	if compilationError != nil {
		compilationWarningAndError = compilationError.Error()
	}

	if !compilationResult {
		if _, fileFoundError := os.Stat(strings.Split(filename, ".")[0] + ".jar"); errors.Is(fileFoundError, os.ErrNotExist) {
			return types.CompileCodeResponse{
				Success: false,
				Error:   compilationWarningAndError,
			}
		}
	}

	runtimePromise := promise.New(func(resolve func(result string), reject func(error)) {
		execCommand := exec.Command("java", "-jar", strings.Split(filename, ".")[0]+".jar")

		time.AfterFunc(8*time.Second, func() {
			if processKillError := execCommand.Process.Kill(); processKillError != nil {
				log.Println(processKillError.Error())
			}

			reject(errors.New("TLE"))
		})

		stdInput, _ := execCommand.StdinPipe()
		stdError, _ := execCommand.StderrPipe()
		stdOutput, _ := execCommand.StdoutPipe()

		if startCommandError := execCommand.Start(); startCommandError != nil {
			reject(errors.New(startCommandError.Error()))
		}

		if input != "" {
			_, inputError := stdInput.Write([]byte(input))
			inputError = stdInput.Close()

			if inputError != nil {
				reject(errors.New(inputError.Error()))
			}
		}

		stdErrorBytes, _ := io.ReadAll(stdError)

		if stdErrorBytes != nil && len(string(stdErrorBytes)) > 0 {
			reject(errors.New(string(stdErrorBytes)))
		}

		stdOutputBytes, _ := io.ReadAll(stdOutput)

		waitErr := execCommand.Wait()

		if waitErr != nil {
			reject(errors.New(waitErr.Error()))
		}

		resolve(string(stdOutputBytes))
	})

	runtimeResult, runtimeError := runtimePromise.Await()

	helpers.DeleteFile(strings.Split(filename, ".")[0] + ".jar")

	if runtimeError != nil {
		if runtimeError.Error() == "TLE" {
			return types.CompileCodeResponse{
				Success: false,
				Error:   compilationWarningAndError + "Time Limit Exceeded....",
			}
		}

		return types.CompileCodeResponse{
			Success: false,
			Error:   compilationWarningAndError + runtimeError.Error(),
		}
	}

	return types.CompileCodeResponse{
		Success: true,
		Output:  compilationWarningAndError + runtimeResult,
	}
}
