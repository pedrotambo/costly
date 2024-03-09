package handlers

import "fmt"

type ErrorResponse struct {
	APIError *APIError `json:"error"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewErrorResponse(code string, message string) *ErrorResponse {
	return &ErrorResponse{
		APIError: &APIError{
			Code:    code,
			Message: message,
		},
	}
}

func NewInvalidInputResponseError(message string) *ErrorResponse {
	return &ErrorResponse{
		APIError: &APIError{
			Code:    "INVALID_INPUT",
			Message: message,
		},
	}
}

func (re *ErrorResponse) Error() string {
	return fmt.Sprintf("error code: %s, message: %s", re.APIError.Code, re.APIError.Message)
}

var ErrBadName = NewInvalidInputResponseError("name is invalid")
var ErrBadUnit = NewInvalidInputResponseError("unit is invalid")
var ErrBadPrice = NewInvalidInputResponseError("price is invalid")
var ErrBadID = NewInvalidInputResponseError("id is invalid")
var ErrBadNewUnits = NewInvalidInputResponseError("new_units should be more than 0")
var ErrBadJson = NewErrorResponse("INVALID_JSON", "error unmarshalling request body")
var ErrBadIngrs = NewInvalidInputResponseError("recipe must have at least one ingredient")
