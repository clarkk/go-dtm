package dtm

import "errors"

var (
	ErrNotFound			= errors.New("Not found")
	ErrLimitOffsetRange	= errors.New("Limit offset out of bounds")
)

type Error struct {
	code 	int
	err 	error
}

func NewError(code int, err error) *Error {
	return &Error{code, err}
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Error() string {
	return e.err.Error()
}