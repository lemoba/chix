package response

import (
	"github.com/lemoba/chix/pkg/errors"
)

var (
	ErrRecordsNotExist = errors.NewError(3000, "records not exist")
	ErrUnauthorized    = errors.NewError(4001, "unauthenticated")
	ErrBadRequest      = errors.NewError(4002, "bad request")
	ErrForbidden       = errors.NewError(4003, "forbidden")
	ErrNotFound        = errors.NewError(4004, "not found")
)
