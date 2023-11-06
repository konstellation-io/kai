package testhelpers

import (
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"
)

type VersionBuilder struct {
	version domain.Version
}

func defaultVersion() domain.Version {
	return domain.Version{
		Product:              "test-product",
		Tag:                  "v1.0.0",
		VersionKeyValueStore: "v1.0.0-kv-store",
		Workflows: []*domain.Workflow{
			NewWorkflowBuilder().Build(),
		},
		MinioConfiguration: domain.MinioConfiguration{
			User:     "test-minio-user",
			Password: "test-minio-password",
			Bucket:   "test-minio-bucket",
		},
	}
}

func NewVersionBuilder() *VersionBuilder {
	return &VersionBuilder{
		version: defaultVersion(),
	}
}

func (b *VersionBuilder) Build() domain.Version {
	return b.version
}

func (b *VersionBuilder) WithWorkflows(workflows []*domain.Workflow) *VersionBuilder {
	b.version.Workflows = workflows
	return b
}
