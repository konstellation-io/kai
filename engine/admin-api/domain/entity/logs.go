package entity

import (
	"time"
)

type Label struct {
	Key   string
	Value string
}

type Log struct {
	FormatedLog string
	Labels      []Label
}

type LogFilters struct {
	ProductID    string
	VersionTag   string
	From         time.Time
	To           time.Time
	WorkflowName string
	ProcessName  string
	RequestID    string
	Level        string
	Logger       string
	Limit        int
}
