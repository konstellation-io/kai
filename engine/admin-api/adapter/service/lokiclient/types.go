package lokiclient

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Response is the top-level structure returned by Loki's API.
type Response struct {
	Status string `json:"status"`
	Data   Data   `json:"data"`
}

// Data is the data structure returned by Loki's API.
type Data struct {
	ResultType string      `json:"resultType"`
	Result     Streams     `json:"result"`
	Stats      interface{} `json:"stats"`
}

type Streams []Stream

// Stream represents a log stream.  It includes a set of log entries and their labels.
type Stream struct {
	Labels  map[string]string `json:"stream"`
	Entries []Entry           `json:"values"`
}

// Entry represents a log entry. It includes a log message and the time it occurred at.
type Entry struct {
	Timestamp time.Time
	Line      string
}

func (e *Entry) UnmarshalJSON(data []byte) error {
	var unmarshal []string

	err := json.Unmarshal(data, &unmarshal)
	if err != nil {
		return err
	}

	t, err := strconv.ParseInt(unmarshal[0], 10, 64)
	if err != nil {
		return err
	}

	e.Timestamp = time.Unix(0, t)
	e.Line = unmarshal[1]

	return nil
}

type logJSON struct {
	Message string `json:"msg"`
	Level   string `json:"level"`
	Logger  string `json:"logger"`
}

func (logData logJSON) formatLog(timestamp time.Time) string {
	return fmt.Sprintf("%s %s %s %s", timestamp, logData.Level, logData.Logger, logData.Message)
}
