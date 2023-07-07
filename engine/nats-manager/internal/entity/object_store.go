package entity

import (
	"regexp"

	"github.com/konstellation-io/kai/engine/nats-manager/internal"
)

type ObjectStoreScope int

const (
	ScopeUndefined = iota
	ScopeWorkflow
	ScopeProject
)

type ObjectStore struct {
	Name  string
	Scope ObjectStoreScope
}

func (o *ObjectStore) Validate() error {
	isValidName, _ := regexp.MatchString("^[a-z0-9]([-a-z0-9]*[a-z0-9])?$", o.Name)

	if !isValidName {
		return internal.ErrInvalidObjectStoreName
	}

	switch o.Scope {
	case ScopeProject, ScopeWorkflow:
		return nil
	default:
		return internal.ErrInvalidObjectStoreScope
	}
}
