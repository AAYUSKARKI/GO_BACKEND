package response

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOk    = "OK"
	StatusError = "ERROR"
)

func WriteJSON(w http.ResponseWriter, status int, body interface{}) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(body)
}

func GeneralError(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, err.Field()+" is required")
		case "min":
			errMsgs = append(errMsgs, err.Field()+" must be greater than or equal to "+err.Param())
		case "max":
			errMsgs = append(errMsgs, err.Field()+" must be less than or equal to "+err.Param())
		case "email":
			errMsgs = append(errMsgs, err.Field()+" must be a valid email address")
		default:
			errMsgs = append(errMsgs, err.Error())
		}
	}
	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
