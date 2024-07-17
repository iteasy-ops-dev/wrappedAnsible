package model

import (
	"fmt"
)

type NoDocError struct{}

func (e *NoDocError) Error() string {
	return "no exist Documents"
}

func NewNoDocError() error {
	return &NoDocError{}
}

type AlreadyExistsError struct {
	Email string
}

func (e *AlreadyExistsError) Error() string {
	return fmt.Sprintf("already exist user: %s", e.Email)
}

func NewAlreadyExistsError(email string) error {
	return &AlreadyExistsError{Email: email}
}

type UserNotFoundError struct {
	Email string
}

func (e *UserNotFoundError) Error() string {
	return fmt.Sprintf("user not found: %s", e.Email)
}

func NewUserNotFoundError(email string) error {
	return &UserNotFoundError{Email: email}
}

type IncorrectPasswordError struct{}

func (e *IncorrectPasswordError) Error() string {
	return "incorrect password"
}

func NewIncorrectPasswordError() error {
	return &IncorrectPasswordError{}
}
