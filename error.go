package gopkg

type ErrorInterface interface {
	Error() string
}

type Error struct {
	code int
	msg  string
}

func NewError(code int, msg string) *Error {
	return &Error{
		code: code,
		msg:  msg,
	}
}

func (e *Error) Error() string {
	return e.msg
}

func (e *Error) Append(s string) {
	e.msg += s
}

func (e *Error) Set(s string) {
	e.msg = s
}

func (e *Error) GetMsg() string {
	return e.msg
}

func (e *Error) GetCode() int {
	return e.code
}
