package loki

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service"
)

var _ service.LogsService = (*Client)(nil)

type Client struct {
	queryRangeURL string
}

func NewClient(cfg *config.Config) *Client {
	queryRangeURL := fmt.Sprintf("%s/loki/api/v1/query_range", cfg.Loki.Address)

	return &Client{
		queryRangeURL: queryRangeURL,
	}
}

func (c Client) GetLogs(lf entity.LogFilters) ([]*entity.Log, error) {
	results := make([]*entity.Log, 0)

	params := url.Values{}
	params.Add("query", getQuery(lf))
	params.Add("limit", fmt.Sprintf("%d", lf.Limit))
	params.Add("start", strconv.FormatInt(lf.From.UnixNano(), 10))
	params.Add("end", strconv.FormatInt(lf.To.UnixNano(), 10))

	ctx := context.Background()
	fullQuery := c.queryRangeURL + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullQuery, http.NoBody)
	if err != nil {
		return results, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return results, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return results, err
	}

	lokiResponse := Response{}

	err = json.Unmarshal(body, &lokiResponse)
	if err != nil {
		return results, err
	}

	for _, s := range lokiResponse.Data.Result {
		for _, e := range s.Entries {
			logData := logJSON{}

			err = json.Unmarshal([]byte(e.Line), &logData)
			if err != nil {
				return results, err
			}

			results = append(results, &entity.Log{
				FormatedLog: logData.formatLog(e.Timestamp),
				Labels:      getLabels(s.Labels),
			})
		}
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
