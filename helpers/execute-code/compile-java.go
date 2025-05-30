package execute_code

import (
	"context"
	"errors"
	"io"
	"log"
	"os/exec"
	"time"

	"github.com/Ankan002/compiler-api/types"
	"github.com/chebyrash/promise"
)

func CompileJava(filename string, stdInput string) types.CompileCodeResponse {
	compilationPromise := promise.New(func(resolve func(result string), reject func(error)) {
		execCommand := exec.Command("java", filename)

		time.AfterFunc(8*time.Second, func() {
			if processKillError := execCommand.Process.Kill(); processKillError != nil {
				log.Println(processKillError.Error())
			}

			reject(errors.New("TLE"))
		})

		input, _ := execCommand.StdinPipe()
		stdOutput, _ := execCommand.StdoutPipe()
		stdError, _ := execCommand.StderrPipe()

		if startCommandError := execCommand.Start(); startCommandError != nil {
			reject(errors.New(startCommandError.Error()))
		}

		if stdInput != "" {
			_, inputErr := input.Write([]byte(stdInput))
			inputErr = input.Close()

			if inputErr != nil {
				reject(errors.New(inputErr.Error()))
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

	compilationResultRef, compilationError := compilationPromise.Await(context.TODO())
	compilationResult := *compilationResultRef

	if compilationError != nil {
		if compilationError.Error() == "TLE" {
			return types.CompileCodeResponse{
				Success: false,
				Error:   "Time Limit Exceeded...\n",
			}
		}

		return types.CompileCodeResponse{
			Success: false,
			Error:   compilationError.Error(),
		}
	}

	return types.CompileCodeResponse{
		Success: true,
		Output:  compilationResult,
	}
}
