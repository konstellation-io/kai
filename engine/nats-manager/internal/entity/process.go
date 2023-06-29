package entity

import (
	"fmt"

	"github.com/konstellation-io/kai/engine/nats-manager/internal"
)

type Process struct {
	ID            string
	Subscriptions []string
	ObjectStore   *ObjectStore
}

func (n *Process) Validate() error {
	if n.ID == "" {
		return internal.ErrEmptyProcessID
	}

	if n.ObjectStore == nil {
		return nil
	}

	if err := n.ObjectStore.Validate(); err != nil {
		return fmt.Errorf("invalid process object store: %w", err)
	}

	return nil
}
