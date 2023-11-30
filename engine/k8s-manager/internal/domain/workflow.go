package domain

const (
	WorkflowTypeUnknown  WorkflowType = "unknown"
	WorkflowTypeTraining WorkflowType = "training"
	WorkflowTypeData     WorkflowType = "data"
	WorkflowTypeServing  WorkflowType = "serving"
	WorkflowTypeFeedback WorkflowType = "feedback"
)

type Workflow struct {
	Name          string
	Stream        string
	KeyValueStore string
	Type          WorkflowType
	Processes     []*Process
}

type WorkflowType string

func (w WorkflowType) ToString() string {
	return string(w)
}
