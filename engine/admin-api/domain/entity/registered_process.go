package entity

import "time"

const (
	RegisterProcessStatusUnknown  = "UNKNOWN"
	RegisterProcessStatusCreated  = "CREATED"
	RegisterProcessStatusCreating = "CREATING"
	RegisterProcessStatusFailed   = "FAILED"
)

type RegisteredProcess struct {
	ID         string
	Name       string
	Version    string
	Type       string
	Image      string
	UploadDate time.Time
	Owner      string
	Status     string
	Logs       string
}
