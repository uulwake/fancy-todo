package libs

import "net/http"

func DefaultInternalServerError(err error) CustomError {
	return CustomError{
		HTTPCode: http.StatusInternalServerError,
		Message: err.Error(),
	}
}

type CustomError struct {
	HTTPCode int 
	BusinessCode int
	Message string
}

func (ce CustomError) Error() string {
	return ce.Message
}