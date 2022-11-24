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

func CompileCpp(filename string, input string) types.CompileCodeResponse {
	compilationPromise := promise.New(func(resolve func(isCreated bool), reject func(error)) {
		execCommand := exec.Command("g++", "-o", strings.Split(filename, ".")[0], filename)

		time.AfterFunc(8*time.Second, func() {
			if processKillError := execCommand.Process.Kill(); processKillError != nil {
				log.Println(processKillError.Error())
			}

			reject(errors.New("TLE"))
		})

		stdError, _ := execCommand.StderrPipe()

		if commandStartError := execCommand.Start(); commandStartError != nil {
			reject(errors.New(commandStartError.Error()))
		}

		stdErrorBytes, _ := io.ReadAll(stdError)

		if stdErrorBytes != nil && len(string(stdErrorBytes)) > 0 {
			reject(errors.New(string(stdErrorBytes)))
		}

		waitError := execCommand.Wait()

		if waitError != nil {
			reject(errors.New(waitError.Error()))
		}

		resolve(true)
	})

	compilationResult, compilationError := compilationPromise.Await()

	var compilationWarningAndError string

	if compilationError != nil {
		compilationWarningAndError = compilationError.Error()
	}

	if !compilationResult {
		if _, fileFoundError := os.Stat(strings.Split(filename, ".")[0]); errors.Is(fileFoundError, os.ErrNotExist) {
			return types.CompileCodeResponse{
				Success: false,
				Error:   compilationWarningAndError,
			}
		}
	}

	runtimePromise := promise.New(func(resolve func(result string), reject func(error)) {
		execCommand := exec.Command(strings.Split(filename, ".")[0])

		time.AfterFunc(8*time.Second, func() {
			if processKilError := execCommand.Process.Kill(); processKilError != nil {
				log.Println(processKilError.Error())
			}

			reject(errors.New("TLE"))
		})

		stdInput, _ := execCommand.StdinPipe()
		stdError, _ := execCommand.StderrPipe()
		stdOutput, _ := execCommand.StdoutPipe()

		if commandStartError := execCommand.Start(); commandStartError != nil {
			reject(errors.New(commandStartError.Error()))
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

	// TODO: Refactor all TLE Errors with escape sequence of \n

	helpers.DeleteFile(strings.Split(filename, ".")[0])

	if runtimeError != nil {
		if runtimeError.Error() == "TLE" {
			return types.CompileCodeResponse{
				Success: false,
				Error:   compilationWarningAndError + "Time Limit Exceeded...\n",
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
