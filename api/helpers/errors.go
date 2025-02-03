package helpers

import "errors"

var ErrInvalidEmailAddress = errors.New("email address syntax is invalid")
var ErrInvaliPhoneNumber = errors.New("invalid phone number")
var ErrMismatchedPassword = errors.New("password and confirm password must match")
