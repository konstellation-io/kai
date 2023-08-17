package entity

import "time"

type ProcessRegistry struct {
	ID         string
	Name       string
	Version    string
	Type       string
	UploadDate time.Time
	Owner      string
}
