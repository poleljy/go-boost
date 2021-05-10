package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrNotFound = errors.New("not found")
)

// Error implements the error interface.
type Error struct {
	Code    int32  `json:"code"`
	Detail  string `json:"error"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

// New generates a custom error.
func New(detail string, code int32) error {
	return &Error{
		Code:    code,
		Detail:  detail,
		Message: http.StatusText(int(code)),
	}
}

// Parse tries to parse a JSON string into an error. If that
// fails, it will set the given string as the error detail.
func Parse(err string) *Error {
	e := new(Error)
	errr := json.Unmarshal([]byte(err), e)
	if errr != nil {
		e.Detail = err
	}
	return e
}

// BadRequest/InvalidRequest generates a 400 error.
func BadRequest(format string, a ...interface{}) error {
	return &Error{
		Code:    400,
		Detail:  fmt.Sprintf(format, a...),
		Message: http.StatusText(400),
	}
}

// Unauthorized generates a 401 error.
func Unauthorized(format string, a ...interface{}) error {
	return &Error{
		Code:    401,
		Detail:  fmt.Sprintf(format, a...),
		Message: http.StatusText(401),
	}
}

// Forbidden generates a 403 error.
func Forbidden(format string, a ...interface{}) error {
	return &Error{
		Code:    403,
		Detail:  fmt.Sprintf(format, a...),
		Message: http.StatusText(403),
	}
}

// NotFound generates a 404 error.
func NotFound(format string, a ...interface{}) error {
	return &Error{
		Code:    404,
		Detail:  fmt.Sprintf(format, a...),
		Message: http.StatusText(404),
	}
}

// 406

// Conflict generates a 409 error.
func Conflict(format string, a ...interface{}) error {
	return &Error{
		Code:    409,
		Detail:  fmt.Sprintf(format, a...),
		Message: http.StatusText(409),
	}
}

// 410

// 422

// InternalServerError generates a 500 error.
func InternalServerError(format string, a ...interface{}) error {
	return &Error{
		Code:    500,
		Detail:  fmt.Sprintf(format, a...),
		Message: http.StatusText(500),
	}
}
