package casbinauth

import (
	"errors"
)

var ErrInvalidAccessControlAction = errors.New("invalid action")
var ErrInvalidNumberOfArguments = errors.New("invalid number of arguments")
var ErrUnauthorized = errors.New("unauthorized")
