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
	VersionStatusCritical  VersionStatus = "CRITICAL"
)

func (e VersionStatus) String() string {
	return string(e)
}

func (v *Version) SetPublishStatus(publicationAuthor string) {
	now := time.Now()

	v.Status = VersionStatusPublished
	v.PublicationAuthor = &publicationAuthor
	v.PublicationDate = &now
}

func (v *Version) UnsetPublishStatus() {
	v.Status = VersionStatusStarted
	v.PublicationAuthor = nil
	v.PublicationAuthor = nil
	v.PublicationDate = nil
}

func (v *Version) CanBeStarted() bool {
	switch v.Status {
	case VersionStatusCreated, VersionStatusStopped, VersionStatusError, VersionStatusCritical:
		return true
	default:
		return false
	}
}

func (v *Version) CanBeStopped() bool {
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
