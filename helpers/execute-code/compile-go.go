package execute_code

import (
	"errors"
	"github.com/Ankan002/compiler-api/types"
	"github.com/chebyrash/promise"
	"io"
	"log"
	"os/exec"
	"time"
)

func CompileGo(filename string, stdInput string) types.CompileCodeResponse {
	compilationPromise := promise.New(func(resolve func(result string), reject func(error)) {
		execCommand := exec.Command("go", "run", filename)

		time.AfterFunc(8*time.Second, func() {
			if processKillingError := execCommand.Process.Kill(); processKillingError != nil {
				log.Println(processKillingError.Error())
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

		waitErr := execCommand.Wait()

		if waitErr != nil {
			reject(errors.New(waitErr.Error()))
		}

		resolve(string(stdOutputBytes))
	})

	compilationResult, compilationError := compilationPromise.Await()

	if compilationError != nil {
		if compilationError.Error() == "TLE" {
			return types.CompileCodeResponse{
				Success: false,
				Error:   "Time Limit Exceeded",
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
