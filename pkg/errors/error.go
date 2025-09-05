package errors

import (
	"fmt"
	"sync"
)

type Error struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
}

func (e *Error) Error() string {
	return e.Msg
}

// registry store all registered errors
var (
	registry = map[int]*Error{}
	regMu    sync.RWMutex
)

// NewError creates and registers a new error
func NewError(code int, msg string) *Error {
	regMu.Lock()
	defer regMu.Unlock()
	if _, exists := registry[code]; exists {
		panic(fmt.Sprintf("error code %d already registered", code))
	}
	err := &Error{
		Code: code,
		Msg:  msg,
	}
	registry[code] = err
	return err
}

// GetError retrieves an error by code
func GetError(code int) (*Error, bool) {
	regMu.Lock()
	defer regMu.Unlock()
	e, ok := registry[code]
	return e, ok
}
