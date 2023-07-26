package domain

//nolint:maligned
type Process struct {
	Name          string
	Image         string
	EnableGpu     bool
	Type          ProcessType
	Subject       string
	ObjectStore   *string
	Subscriptions []string
	KeyValueStore string
	Config        map[string]string

	Replicas   int32
	Networking *Networking
	CPU        *ProcessCPU
	Memory     *ProcessMemory
}

func (p *Process) IsTrigger() bool {
	return p.Type == TriggerProcessType
}

type Networking struct {
	SourcePort int
	Protocol   string
	TargetPort int
}

type ProcessType int

const (
	UnknownProcessType = iota
	TriggerProcessType
	TaskProcessType
	ExitProcessType
)

func (p ProcessType) ToString() string {
	switch p {
	case TriggerProcessType:
		return "trigger"
	case TaskProcessType:
		return "task"
	case ExitProcessType:
		return "exit"
	default:
		return "unknown"
	}
}

type ProcessCPU struct {
	Request string
	Limit   string
}

type ProcessMemory struct {
	Request string
	Limit   string
}
