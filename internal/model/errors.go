package model

import (
	"fmt"
)

type NoDocError struct{}

func NewNoDocError() error          { return &NoDocError{} }
func (e *NoDocError) Error() string { return "no exist Documents" }

type AlreadyExistsError struct {
	Email string
}

func NewAlreadyExistsError(email string) error { return &AlreadyExistsError{Email: email} }
func (e *AlreadyExistsError) Error() string    { return fmt.Sprintf("already exist user: %s", e.Email) }

type UserNotFoundError struct {
	Email string
}

func NewUserNotFoundError(email string) error { return &UserNotFoundError{Email: email} }
func (e *UserNotFoundError) Error() string    { return fmt.Sprintf("user not found: %s", e.Email) }

type IncorrectPasswordError struct{}

func NewIncorrectPasswordError() error          { return &IncorrectPasswordError{} }
func (e *IncorrectPasswordError) Error() string { return "incorrect password" }

type UserNotActiveError struct {
	Email string
}

func NewUserNotActiveError(email string) error { return &UserNotActiveError{Email: email} }
func (e *UserNotActiveError) Error() string    { return fmt.Sprintf("user %s is not active", e.Email) }

type UserNotVerifiedError struct {
	Email string
}

func NewUserNotVerifiedError(email string) error { return &UserNotVerifiedError{Email: email} }
func (e *UserNotVerifiedError) Error() string    { return fmt.Sprintf("user %s is not verified", e.Email) }

type UserLogoutError struct {
	Email string
}

func NewUserLogoutError(email string) error { return &UserLogoutError{Email: email} }
func (e *UserLogoutError) Error() string    { return fmt.Sprintf("user %s: logout Error", e.Email) }

type AlreadyLoginError struct {
	Email string
}

func NewAlreadyLoginError(email string) error { return &AlreadyLoginError{Email: email} }
func (e *AlreadyLoginError) Error() string {
	return fmt.Sprintf("user %s is already logged in", e.Email)
}
