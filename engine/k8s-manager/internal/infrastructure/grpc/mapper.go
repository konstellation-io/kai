package grpc

import (
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/grpc/proto/versionpb"
)

func mapRequestToVersion(req *versionpb.StartRequest) domain.Version {
	return domain.Version{
		Product:       req.ProductId,
		ID:            req.VersionId,
		KeyValueStore: req.KeyValueStore,
		Workflows:     mapReqWorkflowsToWorkflows(req.Workflows),
	}
}

func mapReqWorkflowsToWorkflows(reqWorkflows []*versionpb.Workflow) []*domain.Workflow {
	workflows := make([]*domain.Workflow, 0, len(reqWorkflows))

	for _, workflow := range reqWorkflows {
		workflows = append(workflows, &domain.Workflow{
			ID:            workflow.Id,
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
			ID:            process.Id,
			Type:          mapReqProcessTypeTpProcessType(process.Type),
			Image:         process.Image,
			Subject:       process.Subject,
			Replicas:      process.Replicas,
			EnableGpu:     process.Gpu,
			Subscriptions: process.Subscriptions,
			KeyValueStore: process.KeyValueStore,
			ObjectStore:   process.ObjectStore,
		}

		if process.Networking != nil {
			p.Networking = &domain.Networking{
				SourcePort: int(process.Networking.SourcePort),
				TargetPort: int(process.Networking.TargetPort),
				Protocol:   process.Networking.Protocol,
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
