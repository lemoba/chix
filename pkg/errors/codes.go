package errors

var (
	ErrSuccess             = NewError(0, "success")
	ErrFailure             = NewError(-1, "failure")
	ErrInvalidRequestParam = NewError(2000, "invalid request parameters")
	ErrInternalServerError = NewError(5000, "internal server error")
)
