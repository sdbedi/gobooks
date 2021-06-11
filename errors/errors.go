package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	// ErrInternal HTTP 500
	ErrInternal = &Error{
		Code:    http.StatusInternalServerError,
		Message: "Something went wrong",
	}
	// ErrUnprocessableEntity HTTP 422
	ErrUnprocessableEntity = &Error{
		Code:    http.StatusUnprocessableEntity,
		Message: "Unprocessable Entity",
	}
	// ErrBadRequest HTTP 400
	ErrBadRequest = &Error{
		Code:    http.StatusBadRequest,
		Message: "Error invalid argument",
	}
	// ErrStatusIsRequired HTTP 400
	ErrStatusIsRequired = &Error{
		Code:    http.StatusBadRequest,
		Message: "Please a provide Status of CheckedIn or CheckedOut",
	}
	// ErrRatingIsRequired HTTP 400
	ErrRatingIsRequired = &Error{
		Code:    http.StatusBadRequest,
		Message: "Rating must be 1-3",
	}
	// ErrNotFound HTTP 404
	ErrBookNotFound = &Error{
		Code:    http.StatusNotFound,
		Message: "Book not found",
	}
	// ErrObjectIsRequired HTTP 400
	ErrObjectIsRequired = &Error{
		Code:    http.StatusBadRequest,
		Message: "Request object should be provided",
	}
	// ErrValidBookIDIsRequired HTTP 400
	ErrValidBookIdIsRequired = &Error{
		Code:    http.StatusBadRequest,
		Message: "A valid book id is required",
	}
	// ErrValidBookIDIsRequired HTTP 400
	ErrTitleandAuthorIsRequired = &Error{
		Code:    http.StatusBadRequest,
		Message: "A title and author are required",
	}
	// ErrInvalidLimit HTTP 400
	ErrInvalidLimit = &Error{
		Code:    http.StatusBadRequest,
		Message: "Limit should be an integral value",
	}
)

// Error main object for error
type Error struct {
	Code    int
	Message string
}

func (err *Error) Error() string {
	return err.String()
}

func (err *Error) String() string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("error: code=%s message=%s", http.StatusText(err.Code), err.Message)
}

// JSON convert Error in json
func (err *Error) JSON() []byte {
	if err == nil {
		return []byte("{}")
	}
	res, _ := json.Marshal(err)
	return res
}

// StatusCode get status code
func (err *Error) StatusCode() int {
	if err == nil {
		return http.StatusOK
	}
	return err.Code
}
