package entity

import (
	"regexp"

	"github.com/konstellation-io/kai/engine/nats-manager/internal"
)

type ObjectStoreScope string

const (
	ObjStoreScopeUndefined ObjectStoreScope = "undefined"
	ObjStoreScopeWorkflow  ObjectStoreScope = "workflow"
	ObjStoreScopeProject   ObjectStoreScope = "project"
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
	case ObjStoreScopeProject, ObjStoreScopeWorkflow:
		return nil
	case ObjStoreScopeUndefined:
		return internal.ErrInvalidObjectStoreScope
	default:
		return internal.ErrInvalidObjectStoreScope
	}
}
