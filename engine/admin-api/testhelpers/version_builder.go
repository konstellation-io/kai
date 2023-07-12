package testhelpers

import "github.com/konstellation-io/kai/engine/admin-api/domain/entity"

type VersionBuilder struct {
	version entity.Version
}

func NewVersionBuilder() *VersionBuilder {
	return &VersionBuilder{
		version: entity.Version{
			ID:          "version-id",
			Tag:         "v1.0.0",
			Description: "test description",
			Workflows: []entity.Workflow{
				NewWorkflowBuilder().Build(),
			},
		},
	}
}

func (vb *VersionBuilder) Build() entity.Version {
	return vb.version
}

func (vb *VersionBuilder) WithWorkflows(workflows []entity.Workflow) *VersionBuilder {
	vb.version.Workflows = workflows
	return vb
}
