package helper

import (
	"errors"
)

var (
	ErrUserNotFound         = errors.New("User not found")
	ErrPasswordDoesNotMatch = errors.New("Password does not match")
	ErrUserInvalid          = errors.New("User invalid")
	ErrEmailDuplicate       = errors.New("Email already exists")
	ErrSessionNotFound      = errors.New("Session not found")
	ErrTokenNotFound        = errors.New("Token not found")
	ErrTokenNotMatch        = errors.New("Token doesn't match")
)

const (
	ErrAuthEmailMsg    = "Email entered does not exists"
	ErrAuthPasswordMsg = "Password is incorrect"
	ErrDuplicateEmail  = "Email already registered"
	ErrTokenExpired    = "Token has been expired"
	ErrTokenClaimed    = "Token has been claimed"
	ErrInvalidUserID   = "Invalid user ID"
	ErrInvalidToken    = "Invalid token"

	SuccessSignup = "Sign Up Successfully"
	SuccessLogin  = "Login Successfully"
)
