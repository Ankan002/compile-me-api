package execute_code

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Ankan002/compiler-api/helpers"
	"github.com/Ankan002/compiler-api/types"
	"github.com/chebyrash/promise"
)

// TODO: Think of how to make the standard command as csc [upgrade from mcs (maybe to mcs-mono by only one file)]

func CompileCSharp(filename string, input string) types.CompileCodeResponse {
	compilationPromise := promise.New(func(resolve func(isCreated bool), reject func(error)) {
		execCommand := exec.Command("mcs", "/target:exe", "/out:"+strings.Split(filename, ".")[0]+".exe", filename)

		time.AfterFunc(10*time.Second, func() {
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

	compilationResultRef, compilationError := compilationPromise.Await(context.TODO())
	compilationResult := *compilationResultRef

	var compilationWarningAndError string

	if compilationError != nil {
		if compilationError.Error() == "TLE" {
			compilationWarningAndError = "Time Limit Exceeded...\n"
		} else {
			compilationWarningAndError = compilationError.Error()
		}
	}

	if !compilationResult {
		if _, fileFoundError := os.Stat(strings.Split(filename, ".")[0] + ".exe"); errors.Is(fileFoundError, os.ErrNotExist) {
			return types.CompileCodeResponse{
				Success: false,
				Error:   compilationWarningAndError,
			}
		}
	}

	runtimePromise := promise.New(func(resolve func(result string), reject func(error)) {
		execCommand := exec.Command("mono", strings.Split(filename, ".")[0]+".exe")

		time.AfterFunc(8*time.Second, func() {
			if processKillError := execCommand.Process.Kill(); processKillError != nil {
				log.Println(processKillError.Error())
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

	runtimeResultRef, runtimeError := runtimePromise.Await(context.TODO())
	runtimeResult := *runtimeResultRef

	helpers.DeleteFile(strings.Split(filename, ".")[0] + ".exe")

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
		Output:  runtimeResult,
	}
}
