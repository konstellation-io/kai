package testhelpers

import "github.com/konstellation-io/kai/engine/admin-api/domain/entity"

type WorkflowBuilder struct {
	workflow entity.Workflow
}

func NewWorkflowBuilder() *WorkflowBuilder {
	return &WorkflowBuilder{
		entity.Workflow{
			Name: "test-workflow-name",
			Type: entity.WorkflowTypeTraining,
			Processes: []entity.Process{
				NewProcessBuilder().Build(),
			},
		},
	}
}

func (wb *WorkflowBuilder) Build() entity.Workflow {
	return wb.workflow
}

func (wb *WorkflowBuilder) WithProcesses(processes []entity.Process) *WorkflowBuilder {
	wb.workflow.Processes = processes
	return wb
}

func (wb *WorkflowBuilder) WithConfig(config []entity.ConfigurationVariable) *WorkflowBuilder {
	wb.workflow.Config = config
	return wb
}
