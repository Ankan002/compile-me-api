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

func CompileRust(filename string, input string) types.CompileCodeResponse {
	compilationPromise := promise.New(func(resolve func(isCreated bool), reject func(error)) {
		execCommand := exec.Command("rustc", filename, "-o", strings.Split(filename, ".")[0])

		time.AfterFunc(8*time.Second, func() {
			if processKillError := execCommand.Process.Kill(); processKillError != nil {
				log.Println(processKillError)
			}

			reject(errors.New("compilation time exceeded"))
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
		if compilationError.Error() == "TLE" {
			compilationWarningAndError = "Time Limit Exceeded...\n"
		} else {
			compilationWarningAndError = compilationError.Error()
		}
	}

	if !compilationResult {
		if _, fileFoundErr := os.Stat(strings.Split(filename, ".")[0]); errors.Is(fileFoundErr, os.ErrNotExist) {
			return types.CompileCodeResponse{
				Success: false,
				Error:   compilationWarningAndError,
			}
		}
	}

	runtimePromise := promise.New(func(resolve func(result string), reject func(error)) {
		execCommand := exec.Command(strings.Split(filename, ".")[0])

		time.AfterFunc(8*time.Second, func() {
			if processKillError := execCommand.Process.Kill(); processKillError != nil {
				log.Println(processKillError)
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

		waitError := execCommand.Wait()

		if waitError != nil {
			reject(errors.New(waitError.Error()))
		}

		resolve(string(stdOutputBytes))
	})

	runtimeResult, runtimeError := runtimePromise.Await()

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
