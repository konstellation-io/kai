package entity

import "errors"

var (
	ErrInvalidProcessType = errors.New("invalid process type")
)

type Process struct {
	Name           string
	Type           ProcessType
	Image          string
	Replicas       int32
	GPU            bool
	Config         []ConfigurationVariable
	ObjectStore    *ProcessObjectStore
	Secrets        []ConfigurationVariable
	Subscriptions  []string
	Networking     *ProcessNetworking
	ResourceLimits *ProcessResourceLimits
	Status         ProcessStatus
}

type ProcessType string

const (
	ProcessTypeTrigger ProcessType = "trigger"
	ProcessTypeTask    ProcessType = "task"
	ProcessTypeExit    ProcessType = "exit"
)

func (pt ProcessType) Validate() error {
	switch pt {
	case ProcessTypeTrigger, ProcessTypeTask, ProcessTypeExit:
		return nil
	default:
		return ErrInvalidProcessType
	}
}

func (pt ProcessType) String() string {
	return string(pt)
}

type ProcessObjectStore struct {
	Name  string
	Scope ObjectStoreScope
}

type ObjectStoreScope string

const (
	ObjectStoreScopeProduct  ObjectStoreScope = "product"
	ObjectStoreScopeWorkflow ObjectStoreScope = "workflow"
)

func (s ObjectStoreScope) String() string {
	return string(s)
}

type ProcessNetworking struct {
	TargetPort      int
	DestinationPort int
	Protocol        NetworkingProtocol
}

type NetworkingProtocol string

const (
	NetworkingProtocolHTTP NetworkingProtocol = "HTTP"
	NetworkingProtocolGRPC NetworkingProtocol = "GRPC"
)

type ResourceLimit struct {
	Request string
	Limit   string
}

type ProcessResourceLimits struct {
	CPU    *ResourceLimit
	Memory *ResourceLimit
}

type ProcessStatus string

const (
	ProcessStatusStarting ProcessStatus = "STARTING"
	ProcessStatusStarted  ProcessStatus = "STARTED"
	ProcessStatusStopped  ProcessStatus = "STOPPED"
	ProcessStatusError    ProcessStatus = "ERROR"
)

func (ps ProcessStatus) IsValid() bool {
	switch ps {
	case ProcessStatusStarting, ProcessStatusStopped, ProcessStatusStarted, ProcessStatusError:
		return true
	default:
		return false
	}
}

func (ps ProcessStatus) String() string {
	return string(ps)
}
