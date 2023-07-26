package grpc

import (
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/grpc/proto/versionpb"
)

func mapRequestToVersion(req *versionpb.StartRequest) domain.Version {
	return domain.Version{
		Product:       req.ProductId,
		Tag:           req.VersionTag,
		KeyValueStore: req.KeyValueStore,
		Workflows:     mapReqWorkflowsToWorkflows(req.Workflows),
	}
}

func mapReqWorkflowsToWorkflows(reqWorkflows []*versionpb.Workflow) []*domain.Workflow {
	workflows := make([]*domain.Workflow, 0, len(reqWorkflows))

	for _, workflow := range reqWorkflows {
		workflows = append(workflows, &domain.Workflow{
			Name:          workflow.Name,
			Stream:        workflow.Stream,
			KeyValueStore: workflow.KeyValueStore,
			Processes:     mapReqProcessToProcess(workflow.Processes),
		})
	}

	return workflows
}

func mapReqProcessToProcess(reqProcesses []*versionpb.Process) []*domain.Process {
	processes := make([]*domain.Process, 0, len(reqProcesses))

	for _, process := range reqProcesses {
		p := &domain.Process{
			Name:          process.Name,
			Type:          mapReqProcessTypeTpProcessType(process.Type),
			Image:         process.Image,
			Subject:       process.Subject,
			Replicas:      process.Replicas,
			EnableGpu:     process.Gpu,
			Subscriptions: process.Subscriptions,
			KeyValueStore: process.KeyValueStore,
			ObjectStore:   process.ObjectStore,
			Config:        process.Config,
		}

		if process.Networking != nil {
			p.Networking = &domain.Networking{
				SourcePort: int(process.Networking.SourcePort),
				TargetPort: int(process.Networking.TargetPort),
				Protocol:   process.Networking.Protocol,
			}
		}

		if process.Cpu != nil {
			p.CPU = &domain.ProcessCPU{
				Request: process.Cpu.Request,
				Limit:   process.Cpu.Limit,
			}
		}

		if process.Memory != nil {
			p.Memory = &domain.ProcessMemory{
				Request: process.Memory.Request,
				Limit:   process.Memory.Limit,
			}
		}

		processes = append(processes, p)
	}

	return processes
}

func mapReqProcessTypeTpProcessType(processType versionpb.ProcessType) domain.ProcessType {
	switch processType {
	case versionpb.ProcessType_ProcessTypeTrigger:
		return domain.TriggerProcessType
	case versionpb.ProcessType_ProcessTypeTask:
		return domain.TaskProcessType
	case versionpb.ProcessType_ProcessTypeExit:
		return domain.ExitProcessType
	case versionpb.ProcessType_ProcessTypeUnknown:
		return domain.UnknownProcessType
	default:
		return domain.UnknownProcessType
	}
}
