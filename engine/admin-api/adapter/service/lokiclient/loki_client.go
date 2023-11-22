package lokiclient

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/grafana/loki/pkg/logcli/client"
	"github.com/grafana/loki/pkg/loghttp"
	"github.com/grafana/loki/pkg/logproto"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service"
)

var _ service.LogsService = (*Client)(nil)

type Client struct {
	lokiClient client.Client
}

func NewClient(cfg *config.Config) *Client {
	defaultClient := client.DefaultClient{}
	defaultClient.Address = cfg.Loki.Address

	return &Client{
		lokiClient: &defaultClient,
	}
}

type logJSON struct {
	Message string `json:"msg"`
	Level   string `json:"level"`
	Logger  string `json:"logger"`
}

func formatLog(timestamp time.Time, logData logJSON) string {
	return fmt.Sprintf("%s %s %s %s", timestamp, logData.Level, logData.Logger, logData.Message)
}

func (c Client) GetLogs(lf entity.LogFilters) ([]*entity.Log, error) {
	results := make([]*entity.Log, 0)

	query := getQuery(lf)

	queryRes, err := c.lokiClient.QueryRange(
		query,
		lf.Limit,
		lf.From,
		lf.To,
		logproto.BACKWARD,
		0,
		0,
		false,
	)
	if err != nil {
		return results, err
	}

	dataResult := queryRes.Data.Result
	switch dataResult.Type() {
	case loghttp.ResultTypeStream:
		series := dataResult.(loghttp.Streams)
		for _, s := range series {
			for _, e := range s.Entries {
				logData := logJSON{}
				err = json.Unmarshal([]byte(e.Line), &logData)
				if err != nil {
					return results, err
				}
				results = append(results, &entity.Log{
					FormatedLog: formatLog(e.Timestamp, logData),
					Labels:      getLabels(s.Labels),
				})
			}
		}
	default:
		return results, fmt.Errorf(`result type %q is not expected "stream" type`, dataResult.Type())
	}

	return results, nil
}

func getLabels(labelsMap map[string]string) []entity.Label {
	labels := make([]entity.Label, 0)

	for k, v := range labelsMap {
		labels = append(labels, entity.Label{
			Key:   k,
			Value: v,
		})
	}

	return labels
}
