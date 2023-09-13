package utils

import (
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/krt"
)

type TestVersion struct {
	version *entity.Version
}

func InitTestVersion() *TestVersion {
	commonObjectStore := &entity.ProcessObjectStore{
		Name:  "emails",
		Scope: "workflow",
	}

	return &TestVersion{
		&entity.Version{
			Description: "Play Golden Sun II",
			Config: []entity.ConfigurationVariable{
				{
					Key:   "keyA",
					Value: "value1",
				},
				{
					Key:   "keyB",
					Value: "value2",
				},
			},
			Workflows: []entity.Workflow{
				{
					Name: "go-classificator",
					Type: "data",
					Config: []entity.ConfigurationVariable{
						{
							Key:   "keyA",
							Value: "value1",
						},
						{
							Key:   "keyB",
							Value: "value2",
						},
					},
					Processes: []entity.Process{
						{
							Name:          "entrypoint",
							Type:          "trigger",
							Image:         "konstellation/kai-grpc-trigger:latest",
							Replicas:      krt.DefaultNumberOfReplicas,
							GPU:           krt.DefaultGPUValue,
							Subscriptions: []string{"exitpoint"},
							Networking: &entity.ProcessNetworking{
								TargetPort:      9000,
								DestinationPort: 9000,
								Protocol:        "TCP",
							},
						},
						{
							Name:          "etl",
							Type:          "task",
							Image:         "konstellation/kai-etl-task:latest",
							Replicas:      krt.DefaultNumberOfReplicas,
							GPU:           krt.DefaultGPUValue,
							ObjectStore:   commonObjectStore,
							Subscriptions: []string{"entrypoint"},
						},
						{
							Name:          "email-classificator",
							Type:          "task",
							Image:         "konstellation/kai-ec-task:latest",
							Replicas:      krt.DefaultNumberOfReplicas,
							GPU:           krt.DefaultGPUValue,
							ObjectStore:   commonObjectStore,
							Subscriptions: []string{"etl"},
						},
						{
							Name:          "exitpoint",
							Type:          "exit",
							Image:         "konstellation/kai-exitpoint:latest",
							Replicas:      krt.DefaultNumberOfReplicas,
							GPU:           krt.DefaultGPUValue,
							ObjectStore:   commonObjectStore,
							Subscriptions: []string{"etl", "stats-storer"},
						},
					},
				},
			},
		},
	}
}

func (v *TestVersion) WithVersionID(versionID string) *TestVersion {
	v.version.ID = versionID
	return v
}

func (v *TestVersion) WithTag(tag string) *TestVersion {
	v.version.Tag = tag
	return v
}

func (v *TestVersion) WithStatus(status entity.VersionStatus) *TestVersion {
	v.version.Status = status
	return v
}

func (v *TestVersion) GetVersion() *entity.Version {
	return v.version
}
