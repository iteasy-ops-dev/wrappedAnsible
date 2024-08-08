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

type UserNotActiveError struct {
	Email string
}

func (e *UserNotActiveError) Error() string {
	return fmt.Sprintf("user %s is not active", e.Email)
}

func NewUserNotActiveError(email string) error {
	return &UserNotActiveError{Email: email}
}

type UserNotVerifiedError struct {
	Email string
}

func (e *UserNotVerifiedError) Error() string {
	return fmt.Sprintf("user %s is not verified", e.Email)
}

func NewUserNotVerifiedError(email string) error {
	return &UserNotVerifiedError{Email: email}
}

type UserLogoutError struct {
	Email string
}

func (e *UserLogoutError) Error() string {
	return fmt.Sprintf("user %s: logout Error", e.Email)
}

func NewUserLogoutError(email string) error {
	return &UserLogoutError{Email: email}
}

type AlreadyLoginError struct {
	Email string
}

func (e *AlreadyLoginError) Error() string {
	return fmt.Sprintf("user %s is already logged in", e.Email)
}

func NewAlreadyLoginError(email string) error {
	return &AlreadyLoginError{Email: email}
}
