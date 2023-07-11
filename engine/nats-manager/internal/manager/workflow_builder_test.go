//go:build unit

package manager_test

import (
	"github.com/konstellation-io/kai/engine/nats-manager/internal/entity"
)

type WorkflowBuilder struct {
	workflow *entity.Workflow
}

func NewWorkflowBuilder() *WorkflowBuilder {
	return &WorkflowBuilder{
		&entity.Workflow{
			Name: "test-workflow",
			Processes: []entity.Process{
				{
					Name: "defaultProcess",
				},
			},
		},
	}
}

func (w *WorkflowBuilder) Build() entity.Workflow {
	return *w.workflow
}

func (w *WorkflowBuilder) WithID(name string) *WorkflowBuilder {
	w.workflow.Name = name
	return w
}

func (w *WorkflowBuilder) WithProcessName(name string) *WorkflowBuilder {
	w.workflow.Processes[0].Name = name
	return w
}

func (w *WorkflowBuilder) WithProcessSubscriptions(subscriptions []string) *WorkflowBuilder {
	w.workflow.Processes[0].Subscriptions = subscriptions
	return w
}

func (w *WorkflowBuilder) WithProcessObjectStore(objectStore *entity.ObjectStore) *WorkflowBuilder {
	w.workflow.Processes[0].ObjectStore = objectStore
	return w
}

func (w *WorkflowBuilder) WithProcesses(processes []entity.Process) *WorkflowBuilder {
	w.workflow.Processes = processes
	return w
}
