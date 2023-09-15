package testhelpers

import "github.com/konstellation-io/kai/engine/admin-api/domain/entity"

type VersionBuilder struct {
	version *entity.Version
}

func NewVersionBuilder() *VersionBuilder {
	return &VersionBuilder{
		version: &entity.Version{
			ID:          "version-id",
			Tag:         "v1.0.0",
			Description: "test description",
			Workflows: []entity.Workflow{
				NewWorkflowBuilder().Build(),
			},
		},
	}
}

func (vb *VersionBuilder) Build() *entity.Version {
	return vb.version
}

func (vb *VersionBuilder) WithID(versionID string) *VersionBuilder {
	vb.version.ID = versionID
	return vb
}

func (vb *VersionBuilder) WithTag(tag string) *VersionBuilder {
	vb.version.Tag = tag
	return vb
}

func (vb *VersionBuilder) WithStatus(status entity.VersionStatus) *VersionBuilder {
	vb.version.Status = status
	return vb
}

func (vb *VersionBuilder) WithWorkflows(workflows []entity.Workflow) *VersionBuilder {
	vb.version.Workflows = workflows
	return vb
}
