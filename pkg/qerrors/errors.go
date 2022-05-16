package qerrors

import (
	"errors"
	"fmt"
)

type Error struct {
	State     int               `json:"state"`     // 状态值
	StateInfo string            `json:"stateInfo"` // 状态描述
	Metadata  map[string]string `json:"metadata"`  // 自定义描述
}

func (e *Error) Error() string {
	return fmt.Sprintf("error:state = %d stateInfo = %s metadata = %v", e.State, e.StateInfo, e.Metadata)
}

// Is matches each error in the chain with the target value.
func (e *Error) Is(err error) bool {
	if se := new(Error); errors.As(err, &se) {
		return se.State == e.State
	}
	return false
}

// WithMetadata with an MD formed by the mapping of key, value.
func (e *Error) WithMetadata(md map[string]string) *Error {
	e.Metadata = md
	return e
}

// New returns an error object for the ss, message.
func New(state int, stateInfo string) *Error {
	return &Error{
		State:     state,
		StateInfo: stateInfo,
	}
}

// Newf New(code fmt.Sprintf(format, a...))
func Newf(state int, format string, a ...interface{}) *Error {
	return New(state, fmt.Sprintf(format, a...))
}

// Errorf returns an error object for the code, message and error info.
func Errorf(state int, format string, a ...interface{}) error {
	return New(state, fmt.Sprintf(format, a...))
}

// State returns the state for a particular error.
// It supports wrapped errors.
func State(err error) int {
	if err == nil {
		return 0
	}
	if se := FromError(err); err != nil {
		return se.State
	}
	return -1
}

// FromError try to convert an error to *Error.
// It supports wrapped errors.
func FromError(err error) *Error {
	if err == nil {
		return nil
	}
	if se := new(Error); errors.As(err, &se) {
		return se
	}

	return New(2, err.Error())
}
