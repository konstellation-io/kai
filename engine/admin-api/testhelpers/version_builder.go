package testhelpers

import "github.com/konstellation-io/kai/engine/admin-api/domain/entity"

type VersionBuilder struct {
	version entity.Version
}

func NewVersionBuilder() *VersionBuilder {
	return &VersionBuilder{
		version: entity.Version{
			ID:          "v1.0.0-id",
			Name:        "v1.0.0-name",
			Description: "test description",
			Version:     "v1.0.0",
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
