// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gql

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

type ConfigurationVariable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CreateProductInput struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateVersionInput struct {
	File      graphql.Upload `json:"file"`
	ProductID string         `json:"productID"`
}

type LogPage struct {
	Cursor *string              `json:"cursor,omitempty"`
	Items  []*entity.ProcessLog `json:"items"`
}

type Process struct {
	ID            string                   `json:"id"`
	Name          string                   `json:"name"`
	ProcessType   ProcessType              `json:"processType"`
	Image         string                   `json:"image"`
	Replicas      int                      `json:"replicas"`
	Config        []*ConfigurationVariable `json:"config,omitempty"`
	Subscriptions []string                 `json:"subscriptions"`
	Status        ProcessStatus            `json:"status"`
}

type PublishVersionInput struct {
	VersionName string `json:"versionName"`
	Comment     string `json:"comment"`
	ProductID   string `json:"productID"`
}

type RevokeUserProductGrantsInput struct {
	TargetID string  `json:"targetID"`
	Product  string  `json:"product"`
	Comment  *string `json:"comment,omitempty"`
}

type Settings struct {
	AuthAllowedDomains    []string `json:"authAllowedDomains"`
	SessionLifetimeInDays int      `json:"sessionLifetimeInDays"`
}

type StartVersionInput struct {
	VersionName string `json:"versionName"`
	Comment     string `json:"comment"`
	ProductID   string `json:"productID"`
}

type StopVersionInput struct {
	VersionName string `json:"versionName"`
	Comment     string `json:"comment"`
	ProductID   string `json:"productID"`
}

type UnpublishVersionInput struct {
	VersionName string `json:"versionName"`
	Comment     string `json:"comment"`
	ProductID   string `json:"productID"`
}

type UpdateUserProductGrantsInput struct {
	TargetID string   `json:"targetID"`
	Product  string   `json:"product"`
	Grants   []string `json:"grants"`
	Comment  *string  `json:"comment,omitempty"`
}

type Workflow struct {
	ID           string                   `json:"id"`
	Name         string                   `json:"name"`
	WorkflowType WorkflowType             `json:"workflowType"`
	Config       []*ConfigurationVariable `json:"config,omitempty"`
	Processes    []*Process               `json:"processes"`
}

type ProcessStatus string

const (
	ProcessStatusStarting ProcessStatus = "STARTING"
	ProcessStatusStarted  ProcessStatus = "STARTED"
	ProcessStatusStopped  ProcessStatus = "STOPPED"
	ProcessStatusError    ProcessStatus = "ERROR"
)

var AllProcessStatus = []ProcessStatus{
	ProcessStatusStarting,
	ProcessStatusStarted,
	ProcessStatusStopped,
	ProcessStatusError,
}

func (e ProcessStatus) IsValid() bool {
	switch e {
	case ProcessStatusStarting, ProcessStatusStarted, ProcessStatusStopped, ProcessStatusError:
		return true
	}
	return false
}

func (e ProcessStatus) String() string {
	return string(e)
}

func (e *ProcessStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ProcessStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ProcessStatus", str)
	}
	return nil
}

func (e ProcessStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ProcessType string

const (
	ProcessTypeTrigger ProcessType = "TRIGGER"
	ProcessTypeTask    ProcessType = "TASK"
	ProcessTypeExit    ProcessType = "EXIT"
)

var AllProcessType = []ProcessType{
	ProcessTypeTrigger,
	ProcessTypeTask,
	ProcessTypeExit,
}

func (e ProcessType) IsValid() bool {
	switch e {
	case ProcessTypeTrigger, ProcessTypeTask, ProcessTypeExit:
		return true
	}
	return false
}

func (e ProcessType) String() string {
	return string(e)
}

func (e *ProcessType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ProcessType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ProcessType", str)
	}
	return nil
}

func (e ProcessType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type WorkflowType string

const (
	WorkflowTypeData     WorkflowType = "DATA"
	WorkflowTypeTraining WorkflowType = "TRAINING"
	WorkflowTypeFeedback WorkflowType = "FEEDBACK"
	WorkflowTypeServing  WorkflowType = "SERVING"
)

var AllWorkflowType = []WorkflowType{
	WorkflowTypeData,
	WorkflowTypeTraining,
	WorkflowTypeFeedback,
	WorkflowTypeServing,
}

func (e WorkflowType) IsValid() bool {
	switch e {
	case WorkflowTypeData, WorkflowTypeTraining, WorkflowTypeFeedback, WorkflowTypeServing:
		return true
	}
	return false
}

func (e WorkflowType) String() string {
	return string(e)
}

func (e *WorkflowType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = WorkflowType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid WorkflowType", str)
	}
	return nil
}

func (e WorkflowType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
