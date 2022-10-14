package execute_code

import (
	"errors"
	"github.com/chebyrash/promise"
	"io"
	"log"
	"os/exec"
	"time"
)

type CompileJavascriptResponse struct {
	Success bool
	Error   string
	Output  string
}

func CompileJavascript(filename string, stdInput string) CompileJavascriptResponse {
	compilationPromise := promise.New(func(resolve func(result string), reject func(error)) {
		execCommand := exec.Command("node", filename)

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

		if stdErrBytes != nil && len(string(stdErrBytes)) > 1 {
			reject(errors.New(string(stdErrBytes)))
		}

		stdOutputBytes, _ := io.ReadAll(output)

		waitErr := execCommand.Wait()

		if waitErr != nil {
			reject(errors.New(waitErr.Error()))
		}

		resolve(string(stdOutputBytes))
	})

	compilationResult, compilationError := compilationPromise.Await()

	if compilationError != nil {
		if compilationError.Error() == "TLE" {
			return CompileJavascriptResponse{
				Success: false,
				Error:   "Time Limit Exceeded...",
			}
		}

		return CompileJavascriptResponse{
			Success: false,
			Error:   compilationError.Error(),
		}
	}

	return CompileJavascriptResponse{
		Success: true,
		Output:  compilationResult,
	}
}
