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

func CompileJavascript(filename string, stdInput string) types.CompileCodeResponse {
	compilationPromise := promise.New(func(resolve func(result string), reject func(error)) {
		execCommand := exec.Command("node", "--no-warnings", filename)

		time.AfterFunc(8*time.Second, func() {
			if processKillError := execCommand.Process.Kill(); processKillError != nil {
				log.Println(processKillError)
			}

			reject(errors.New("TLE"))
		})

		input, _ := execCommand.StdinPipe()
		output, _ := execCommand.StdoutPipe()
		stdErr, _ := execCommand.StderrPipe()

		commandStartError := execCommand.Start()

		if commandStartError != nil {
			reject(errors.New(commandStartError.Error()))
		}

		if stdInput != "" {
			_, inputErr := input.Write([]byte(stdInput))
			inputErr = input.Close()

			if inputErr != nil {
				reject(errors.New(inputErr.Error()))
			}
		}

		stdErrBytes, _ := io.ReadAll(stdErr)

		if stdErrBytes != nil && len(string(stdErrBytes)) > 0 {
			reject(errors.New(string(stdErrBytes)))
		}

		stdOutputBytes, _ := io.ReadAll(output)

		waitErr := execCommand.Wait()

		if waitErr != nil {
			reject(errors.New(waitErr.Error()))
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
