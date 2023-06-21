package entity

import (
	"time"

	"github.com/konstellation-io/krt/pkg/krt"
)

type Version struct {
	*krt.Krt
	ID string

	CreationDate   time.Time
	CreationAuthor string

	PublicationDate   *time.Time
	PublicationAuthor *string

	Status VersionStatus
	Errors []string
}

type VersionStatus string

const (
	VersionStatusCreating  VersionStatus = "CREATING"
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

func (v Version) PublishedOrStarted() bool {
	return v.Status == VersionStatusStarted || v.Status == VersionStatusPublished
}

func (v Version) CanBeStarted() bool {
	return v.Status == VersionStatusCreated || v.Status == VersionStatusStopped
}

func (v Version) CanBeStopped() bool {
	return v.Status == VersionStatusStarted
}
