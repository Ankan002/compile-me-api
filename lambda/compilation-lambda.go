package lambda

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Ankan002/compiler-api/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/go-playground/validator/v10"
)

type requestBody struct {
	Code     string `json:"code" validate:"required"`
	Language string `json:"language" validate:"required,eq=js|eq=ts|eq=py|eq=go|eq=java|eq=rs|eq=kt|eq=cpp|eq=c|eq=cs"`
	StdInput string `json:"stdInput"`
}

type SuccessResponseBody struct {
	Success   bool      `json:"success"`
	Output    string    `json:"output"`
	Timestamp time.Time `json:"timestamp"`
}

type ErrorResponseBody struct {
	Success   bool      `json:"success"`
	Error     string    `json:"error"`
	Timestamp time.Time `json:"timestamp"`
}

func CompilationLambda(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	request := requestBody{}

	if err := json.Unmarshal([]byte(req.Body), &request); err != nil {
		response := ErrorResponseBody{
			Success:   false,
			Error:     err.Error(),
			Timestamp: time.Now(),
		}

		responseJson, _ := json.Marshal(response)

		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       string(responseJson),
		}, nil
	}

	if err := validator.New().Struct(request); err != nil {
		var validationError string

		if err.Error() == "Key: 'CompRequest.Language' Error:Field validation for 'Language' failed on the 'eq=js|eq=ts|eq=py|eq=go' tag" {
			validationError = "Provide us with a valid language extension..."
		} else {
			validationError = err.Error()
		}

		response := ErrorResponseBody{
			Success:   false,
			Error:     validationError,
			Timestamp: time.Now(),
		}

		responseJson, _ := json.Marshal(response)

		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       string(responseJson),
		}, nil
	}

	compileResponse := utils.Compile(utils.CompilationArgs(request))

	if compileResponse.Error != "" {
		response := ErrorResponseBody{
			Success:   false,
			Error:     compileResponse.Error,
			Timestamp: time.Now(),
		}

		responseJson, _ := json.Marshal(response)

		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       string(responseJson),
		}, nil
	}

	if compileResponse.StdErr != "" {
		response := ErrorResponseBody{
			Success:   false,
			Error:     compileResponse.StdErr,
			Timestamp: time.Now(),
		}

		responseJson, _ := json.Marshal(response)

		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       string(responseJson),
		}, nil
	}

	response := SuccessResponseBody{
		Success:   true,
		Output:    compileResponse.StdOutput,
		Timestamp: time.Now(),
	}

	responseJson, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseJson),
	}, nil
}
