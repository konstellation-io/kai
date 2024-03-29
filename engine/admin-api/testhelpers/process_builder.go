package testhelpers

import "github.com/konstellation-io/kai/engine/admin-api/domain/entity"

type ProcessBuilder struct {
	process entity.Process
}

func NewProcessBuilder() *ProcessBuilder {
	return &ProcessBuilder{
		entity.Process{
			Name:          "test-process-name",
			Type:          entity.ProcessTypeTask,
			Image:         "test-process-image",
			Replicas:      1,
			GPU:           false,
			Subscriptions: []string{"other-process"},
			ResourceLimits: &entity.ProcessResourceLimits{
				CPU: &entity.ResourceLimit{
					Request: "100m",
					Limit:   "200m",
				},
				Memory: &entity.ResourceLimit{
					Request: "100Mi",
					Limit:   "200Mi",
				},
			},
		},
	}
}

func (pb *ProcessBuilder) Build() entity.Process {
	return pb.process
}

func (pb *ProcessBuilder) WithObjectStore(objectStore *entity.ProcessObjectStore) *ProcessBuilder {
	pb.process.ObjectStore = objectStore
	return pb
}

func (pb *ProcessBuilder) WithNetworking(networking *entity.ProcessNetworking) *ProcessBuilder {
	pb.process.Networking = networking
	return pb
}

func (pb *ProcessBuilder) WithResourceLimits(resourceLimits *entity.ProcessResourceLimits) *ProcessBuilder {
	pb.process.ResourceLimits = resourceLimits
	return pb
}

func (pb *ProcessBuilder) WithConfig(config []entity.ConfigurationVariable) *ProcessBuilder {
	pb.process.Config = config
	return pb
}
