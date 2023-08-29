package entity

import "time"

type RegisteredProcess struct {
	ID         string
	Name       string
	Version    string
	Type       string
	Image      string
	UploadDate time.Time
	Owner      string
}
