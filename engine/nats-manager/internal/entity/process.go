package entity

import (
	"fmt"

	"github.com/konstellation-io/kai/engine/nats-manager/internal"
)

type Process struct {
	Name          string
	Subscriptions []string
	ObjectStore   *ObjectStore
}

func (n *Process) Validate() error {
	if n.Name == "" {
		return internal.ErrEmptyProcessName
	}

	if n.ObjectStore == nil {
		return nil
	}

	if err := n.ObjectStore.Validate(); err != nil {
		return fmt.Errorf("invalid process object store: %w", err)
	}

	return nil
}
