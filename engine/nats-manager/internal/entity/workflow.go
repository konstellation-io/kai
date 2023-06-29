package entity

import (
	"fmt"

	"github.com/konstellation-io/kai/engine/nats-manager/internal"
)

type Workflow struct {
	ID        string
	Processes []Process
}

func (w Workflow) Validate() error {
	if w.ID == "" {
		return internal.ErrEmptyWorkflow
	}

	for _, process := range w.Processes {
		if err := process.Validate(); err != nil {
			return fmt.Errorf("invalid process: %w", err)
		}
	}

	return nil
}
