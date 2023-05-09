package auth

import "errors"

//nolint:gochecknoglobals
var invalidAccessControlResourceError = errors.New("invalid AccessControlResource")
var invalidAccessControlActionError = errors.New("invalid AccessControlAction")
