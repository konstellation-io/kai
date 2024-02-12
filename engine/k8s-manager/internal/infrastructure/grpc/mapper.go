package grpc

import (
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/grpc/proto/versionpb"
)

func mapRequestToVersion(req *versionpb.StartRequest) *domain.Version {
	return &domain.Version{
		Product:              req.ProductId,
		Tag:                  req.VersionTag,
		GlobalKeyValueStore:  req.GlobalKeyValueStore,
		VersionKeyValueStore: req.VersionKeyValueStore,
		Workflows:            mapReqWorkflowsToWorkflows(req.Workflows),
		MinioConfiguration: domain.MinioConfiguration{
			Bucket: req.MinioConfiguration.Bucket,
		},
		ServiceAccount: domain.ServiceAccount{
			Username: req.ServiceAccount.Username,
			Password: req.ServiceAccount.Password,
		},
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
			Type:          mapReqWorkflowTypeToDomain(workflow.Type),
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
			NodeSelectors: process.NodeSelectors,
		}

		if process.Networking != nil {
			p.Networking = &domain.Networking{
				SourcePort: int(process.Networking.SourcePort),
				TargetPort: int(process.Networking.TargetPort),
				Protocol:   domain.NetworkingProtocol(process.Networking.Protocol),
			}
		}

		if process.ResourceLimits != nil {
			p.ResourceLimits = &domain.ProcessResourceLimits{
				CPU: &domain.ResourceLimit{
					Request: process.ResourceLimits.Cpu.Request,
					Limit:   process.ResourceLimits.Cpu.Limit,
				},
				Memory: &domain.ResourceLimit{
					Request: process.ResourceLimits.Memory.Request,
					Limit:   process.ResourceLimits.Memory.Limit,
				},
			}
		}

		processes = append(processes, p)
	}

	return processes
}

func mapReqWorkflowTypeToDomain(workflowType versionpb.WorkflowType) domain.WorkflowType {
	switch workflowType {
	case versionpb.WorkflowType_WorkflowTypeTraining:
		return domain.WorkflowTypeTraining
	case versionpb.WorkflowType_WorkflowTypeData:
		return domain.WorkflowTypeData
	case versionpb.WorkflowType_WorkflowTypeServing:
		return domain.WorkflowTypeServing
	case versionpb.WorkflowType_WorkflowTypeFeedback:
		return domain.WorkflowTypeFeedback
	default:
		return domain.WorkflowTypeUnknown
	}
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
