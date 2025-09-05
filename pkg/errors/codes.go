package errors

var (
	ErrSuccess             = NewError(0, "Success")
	ErrFailure             = NewError(-1, "Failure")
	ErrInvalidRequestParam = NewError(4000, "Invalid Request Parameters")
	ErrUnauthorized        = NewError(4001, "Unauthenticated")
	ErrBadRequest          = NewError(4002, "Bad Request")
	ErrForbidden           = NewError(4003, "Forbidden")
	ErrNotFound            = NewError(4004, "Not Found")
	ErrInternalServerError = NewError(5000, "Internal Server Error")
)
