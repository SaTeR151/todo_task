package entity

import "fmt"

type AppError struct {
	Err       error
	ErrStatus AppErrorStatus
}

type AppErrors []AppError
type AppErrorStatus string

var (
	ErrBadAuth  AppErrorStatus = "bad_password"
	ErrInternal AppErrorStatus = "internal"
	ErrBadINput AppErrorStatus = "bad_input"
)

var (
	ErrNotFound = fmt.Errorf("not found")
)

func (a *AppError) Error() string {
	return a.Err.Error()
}

func (a *AppError) IsNotFound() bool {
	if a.Err == ErrNotFound {
		return true
	}
	return false
}

func (a *AppError) IsBadAuth() bool {
	if a.ErrStatus == ErrBadAuth {
		return true
	}
	return false
}

func (a *AppError) IsInternal() bool {
	if a.ErrStatus == ErrInternal {
		return true
	}
	return false
}

func (a *AppError) IsBadInput() bool {
	if a.ErrStatus == ErrBadINput {
		return true
	}
	return false
}
