package http

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/go-playground/validator/v10"
)

var InternalServerError = GetFailedResponseFromMessage("Internal Server Error")

var statues = [...]string{
	"ok",
	"fail",
}

const (
	OK ResponseStatus = iota
	Fail
)

type ResponseStatus uint8

func (rs ResponseStatus) String() string {
	return statues[rs]
}

func (rs ResponseStatus) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(rs.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (rs *ResponseStatus) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	*rs = 0

	for i, v := range statues {
		if v == s {
			*rs = ResponseStatus(i)
			break
		}
	}

	return nil
}

type StandardResponse struct {
	Status  ResponseStatus `json:"status"`
	Message string         `json:"message"`
	Data    interface{}    `json:"data"`
	Errors  []string       `json:"errors,omitempty"`
}

func GetSuccessResponse(data interface{}) StandardResponse {
	return StandardResponse{
		Status:  OK,
		Message: "success",
		Data:    data,
	}
}

func GetFailedResponseFromMessage(message string) StandardResponse {
	return StandardResponse{
		Status:  Fail,
		Message: message,
		Data:    nil,
	}
}

func GetFailedResponseFromError(err error) StandardResponse {
	var errMsg []string
	errMsg = append(errMsg, err.Error())

	return StandardResponse{
		Status:  Fail,
		Message: err.Error(),
		Data:    nil,
		Errors:  errMsg,
	}
}

func GetFailedResponseFromMessageAndErrors(message string, errors []error) StandardResponse {
	var errMsg []string

	for _, e := range errors {
		errMsg = append(errMsg, e.Error())
	}

	return StandardResponse{
		Status:  Fail,
		Message: message,
		Data:    nil,
		Errors:  errMsg,
	}
}

func GetFailedValidationResponse(err error) StandardResponse {
	var errorsList []error
	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		for _, e := range ve {
			errorsList = append(errorsList, e)
		}
	} else {
		errorsList = append(errorsList, err)
	}

	return GetFailedResponseFromMessageAndErrors("validation error", errorsList)
}
