package entity

import (
	"time"
)

type ConfigurationVariable struct {
	Key   string
	Value string
}

type Version struct {
	Tag         string
	Description string
	Config      []ConfigurationVariable
	Workflows   []Workflow

	CreationDate   time.Time
	CreationAuthor string

	PublicationDate   *time.Time
	PublicationAuthor *string

	Status VersionStatus
	Error  string
}

type VersionStatus string

const (
	VersionStatusCreated   VersionStatus = "CREATED"
	VersionStatusStarting  VersionStatus = "STARTING"
	VersionStatusStarted   VersionStatus = "STARTED"
	VersionStatusPublished VersionStatus = "PUBLISHED"
	VersionStatusStopping  VersionStatus = "STOPPING"
	VersionStatusStopped   VersionStatus = "STOPPED"
	VersionStatusError     VersionStatus = "ERROR"
)

func (e VersionStatus) String() string {
	return string(e)
}

func (v Version) CanBeStarted() bool {
	switch v.Status {
	case VersionStatusCreated, VersionStatusStopped, VersionStatusError:
		return true
	default:
		return false
	}
}

func (v Version) CanBeStopped() bool {
	return v.Status == VersionStatusStarted
}

type Workflow struct {
	Name      string
	Type      WorkflowType
	Config    []ConfigurationVariable
	Processes []Process
	Stream    string
}

type WorkflowType string

const (
	WorkflowTypeData     WorkflowType = "data"
	WorkflowTypeTraining WorkflowType = "training"
	WorkflowTypeFeedback WorkflowType = "feedback"
	WorkflowTypeServing  WorkflowType = "serving"
)

func (wt WorkflowType) String() string {
	return string(wt)
}

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

func (pt ProcessType) IsValid() bool {
	var processTypeMap = map[string]ProcessType{
		string(ProcessTypeTrigger): ProcessTypeTrigger,
		string(ProcessTypeTask):    ProcessTypeTask,
		string(ProcessTypeExit):    ProcessTypeExit,
	}

	_, ok := processTypeMap[string(pt)]

	return ok
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
	var processStatusMap = map[string]ProcessStatus{
		string(ProcessStatusStarting): ProcessStatusStarting,
		string(ProcessStatusStarted):  ProcessStatusStarted,
		string(ProcessStatusStopped):  ProcessStatusStopped,
		string(ProcessStatusError):    ProcessStatusError,
	}

	_, ok := processStatusMap[string(ps)]

	return ok
}

func (ps ProcessStatus) String() string {
	return string(ps)
}
