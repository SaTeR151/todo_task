package errorspkg

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrBadRequest    = errors.New("bad request")
	ErrInternalError = errors.New("internal error")
)

type ValidationError struct {
	Constructor string
	StructName  string
	ErrorMsg    error
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation failed for constructor [%s] struct [%s]: %s",
		e.Constructor,
		e.StructName,
		e.ErrorMsg,
	)
}

func NewValidationError(constructor string, obj any, validationErr error) error {
	var structName string

	if obj != nil {
		structName = reflect.Indirect(reflect.ValueOf(obj)).Type().Name()
	}

	return ValidationError{
		Constructor: constructor,
		StructName:  structName,
		ErrorMsg:    validationErr,
	}
}

type RepoFailedError struct {
	Method    string
	Operation string
	What      string
	ErrorMsg  error
}

func (e RepoFailedError) Error() string {
	return fmt.Sprintf("repo method [%s] failed to %s %s: %v", e.Method, e.Operation, e.What, e.ErrorMsg)
}

func NewRepoFailedError(method, operation, what string, err error) error {
	return RepoFailedError{
		Method:    method,
		Operation: operation,
		What:      what,
		ErrorMsg:  err,
	}
}

type StrconvError struct {
	Method         string
	ConvertedValue string
	ErrorMsg       error
}

func (e StrconvError) Error() string {
	return fmt.Sprintf("[%s] failed to convert [%s]: error=%s", e.Method, e.ConvertedValue, e.ErrorMsg.Error())
}

func NewStrconvError(method, convertedValue string, err error) error {
	return StrconvError{
		Method:         method,
		ConvertedValue: convertedValue,
		ErrorMsg:       err,
	}
}
