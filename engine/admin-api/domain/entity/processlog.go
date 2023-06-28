package entity

import "time"

type ProcessLog struct {
	ID           string   `bson:"_id"`
	Date         string   `bson:"date"`
	Message      string   `bson:"message"`
	Level        LogLevel `bson:"level"`
	PodID        string   `bson:"podId"`
	ProcessID    string   `bson:"processId" gqlgen:"processId"`
	ProcessName  string   `bson:"processName"`
	VersionID    string   `bson:"versionId"`
	VersionName  string   `bson:"versionName"`
	WorkflowID   string   `bson:"workflowId" gqlgen:"workflowId"`
	WorkflowName string   `bson:"workflowName"`
}

type SearchLogsOptions struct {
	StartDate      time.Time
	EndDate        time.Time
	Search         *string
	Levels         []LogLevel
	ProcessIDs     []string
	Cursor         *string
	VersionsIDs    []string
	WorkflowsNames []string
}

type SearchLogsResult struct {
	Cursor string
	Logs   []*ProcessLog
}

type LogFilters struct {
	StartDate      string     `json:"startDate"`
	EndDate        *string    `json:"endDate"`
	Search         *string    `json:"search"`
	Levels         []LogLevel `json:"levels"`
	ProcessIDs     []string   `json:"processIds"`
	VersionsIDs    []string   `json:"versionsIds"`
	WorkflowsNames []string   `json:"workflowsNames"`
}

type LogLevel string

const (
	LogLevelError LogLevel = "ERROR"
	LogLevelWarn  LogLevel = "WARN"
	LogLevelInfo  LogLevel = "INFO"
	LogLevelDebug LogLevel = "DEBUG"
)

func (e LogLevel) IsValid() bool {
	return e == LogLevelError || e == LogLevelWarn || e == LogLevelInfo || e == LogLevelDebug
}

func (e LogLevel) String() string {
	return string(e)
}
