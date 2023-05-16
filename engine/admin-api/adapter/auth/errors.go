package auth

import "errors"

//nolint:gochecknoglobals,stylecheck // needs to be global
var InvalidAccessControlActionError = errors.New("invalid AccessControlAction")
