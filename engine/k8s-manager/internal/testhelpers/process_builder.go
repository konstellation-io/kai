package testhelpers

import "github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"

type ProcessBuilder struct {
	process *domain.Process
}

func NewProcessBuilder() *ProcessBuilder {
	return &ProcessBuilder{
		&domain.Process{
			Name:          "test-process",
			Image:         "test-image@test",
			EnableGpu:     false,
			Type:          domain.TaskProcessType,
			Subject:       "test-subject",
			Replicas:      1,
			Subscriptions: []string{"other-process"},
			KeyValueStore: "test-process-kv-store",
		},
	}
}

func (pb *ProcessBuilder) WithID(id string) *ProcessBuilder {
	pb.process.Name = id
	return pb
}

func (pb *ProcessBuilder) WithNetworking(networking domain.Networking) *ProcessBuilder {
	pb.process.Networking = &networking
	return pb
}

func (pb *ProcessBuilder) WithType(processType domain.ProcessType) *ProcessBuilder {
	pb.process.Type = processType
	return pb
}

func (pb *ProcessBuilder) WithObjectStore(objectStore string) *ProcessBuilder {
	pb.process.ObjectStore = &objectStore
	return pb
}

func (pb *ProcessBuilder) WithReplicas(replicas int32) *ProcessBuilder {
	pb.process.Replicas = replicas
	return pb
}

func (pb *ProcessBuilder) WithEnableGpu(enableGpu bool) *ProcessBuilder {
	pb.process.EnableGpu = enableGpu
	return pb
}

func (pb *ProcessBuilder) WithCPU(cpu *domain.CPUConfig) *ProcessBuilder {
	pb.process.CPU = cpu
	return pb
}

func (pb *ProcessBuilder) WithMemory(memory *domain.MemoryConfig) *ProcessBuilder {
	pb.process.Memory = memory
	return pb
}

func (pb *ProcessBuilder) Build() *domain.Process {
	return pb.process
}
