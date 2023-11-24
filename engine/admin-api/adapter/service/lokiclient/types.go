package lokiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
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
	Labels  LabelSet `json:"stream"`
	Entries []Entry  `json:"values"`
}

type LabelSet map[string]string

// Map coerces LabelSet into a map[string]string. This is useful for working with adapter types.
func (l LabelSet) Map() map[string]string {
	return l
}

// String implements the String interface. It returns a formatted/sorted set of label key/value pairs.
func (l LabelSet) String() string {
	var b bytes.Buffer

	keys := make([]string, 0, len(l))
	for k := range l {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	b.WriteByte('{')

	for i, k := range keys {
		if i > 0 {
			b.WriteByte(',')
			b.WriteByte(' ')
		}

		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(strconv.Quote(l[k]))
	}

	b.WriteByte('}')

	return b.String()
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
