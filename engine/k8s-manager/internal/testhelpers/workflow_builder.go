package testhelpers

import "github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"

type WorkflowBuilder struct {
	workflow *domain.Workflow
}

func NewWorkflowBuilder() *WorkflowBuilder {
	return &WorkflowBuilder{
		&domain.Workflow{
			ID:            "test-workflow",
			Stream:        "test-stream",
			KeyValueStore: "test-workflow-kv-store",
			Processes: []*domain.Process{
				NewProcessBuilder().Build(),
			},
		},
	}
}

func (wb *WorkflowBuilder) Build() *domain.Workflow {
	return wb.workflow
}

func (wb *WorkflowBuilder) WithProcesses(processes []*domain.Process) *WorkflowBuilder {
	wb.workflow.Processes = processes
	return wb
}
