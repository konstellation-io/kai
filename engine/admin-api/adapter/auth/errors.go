package auth

import (
	"errors"
	"fmt"
)

var ErrInvalidAccessControlAction = errors.New("invalid AccessControlAction")
var ErrInvalidNumberOfArguments = errors.New("invalid number of arguments")
var ErrNonAuthorized = errors.New("non authorized")
var ErrNonAdminAccess = errors.New("you are not allowed to")
var ErrNonAuthorizedForProduct = errors.New("you are not allowed to")

func NonAdminAccess(action string) error {
	return fmt.Errorf("%w %q", ErrNonAdminAccess, action)
}

func NonAuthorizedForProductError(action, product string) error {
	return fmt.Errorf("%w %q in product %q", ErrNonAuthorizedForProduct, action, product)
}
