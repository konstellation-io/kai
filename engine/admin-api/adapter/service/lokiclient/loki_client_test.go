//go:build integration

package lokiclient_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/lokiclient"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLokiClientGetLogs(t *testing.T) {
	const responseBody = `{"status":"success","data":{"resultType":"streams","result":[{"stream":{"process_name":"processName","product_id":"productID","service":"kai-product-version","version_tag":"versionID","workflow_name":"workflowName"},"values":[["1700753452257683259","{\"log\":\"[GIN] 2023/11/23 - 15:30:52 | 200 |    1.651831ms |    192.168.49.1 | GET      \\\"/trigger\\\"\"}"]]},{"stream":{"level":"info","logger":"[TRIGGER].[SUBSCRIBER]","process_name":"processName","product_id":"productID","service":"kai-product-version","version_tag":"versionID","workflow_name":"workflowName"},"values":[["1700753452257536602","{\"level\":\"info\",\"ts\":1700753452.257505,\"logger\":\"[TRIGGER].[SUBSCRIBER]\",\"caller\":\"trigger/subscriber.go:97\",\"msg\":\"New message received with subject productID_v1_0_0_workflowName.processName\"}"]]},{"stream":{"level":"info","logger":"[TRIGGER].[RESPONSE HANDLER]","process_name":"processName","product_id":"productID","request_id":"2ba35f04-3de0-42da-870d-40a03b5f704b","service":"kai-product-version","version_tag":"versionID","workflow_name":"workflowName"},"values":[["1700753452257537507","{\"level\":\"info\",\"ts\":1700753452.257514,\"logger\":\"[TRIGGER].[RESPONSE HANDLER]\",\"caller\":\"trigger/helpers.go:78\",\"msg\":\"Message received with request id 2ba35f04-3de0-42da-870d-40a03b5f704b\",\"request_id\":\"2ba35f04-3de0-42da-870d-40a03b5f704b\"}"]]},{"stream":{"level":"info","logger":"[TRIGGER]","process_name":"processName","product_id":"productID","service":"kai-product-version","version_tag":"versionID","workflow_name":"workflowName"},"values":[["1700753452257652864","{\"level\":\"info\",\"ts\":1700753452.257543,\"logger\":\"[TRIGGER]\",\"caller\":\"app/main.go:119\",\"msg\":\"response recieved\",\"response\":\"[type.googleapis.com/google.protobuf.Value]:{struct_value:{fields:{key:\\\"message\\\" value:{string_value:\\\"OK\\\"}} fields:{key:\\\"status_code\\\" value:{string_value:\\\"200\\\"}}}}\"}"]]},{"stream":{"level":"debug","logger":"[TRIGGER].[SUBSCRIBER]","process_name":"processName","product_id":"productID","service":"kai-product-version","version_tag":"versionID","workflow_name":"workflowName"},"values":[["1700753452257533697","{\"level\":\"debug\",\"ts\":1700753452.257471,\"logger\":\"[TRIGGER].[SUBSCRIBER]\",\"caller\":\"trigger/subscriber.go:87\",\"msg\":\"New message received\"}"]]}],"stats":{}}}`

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		expectedURL := `/loki/api/v1/query_range?end=1701388800000000000&limit=100&query={product_id="productID", version_tag="versionID"}&start=1672531200000000000`

		actualQuery, err := url.QueryUnescape(req.URL.String())
		require.NoError(t, err)

		assert.Equal(t, expectedURL, actualQuery)

		rw.Write([]byte(responseBody))
	}))
	defer server.Close()

	cfg := &config.Config{}
	cfg.Loki.Address = server.URL

	fromTime, err := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
	require.NoError(t, err)
	toTime, err := time.Parse(time.RFC3339, "2023-12-01T00:00:00Z")
	require.NoError(t, err)

	logFilters := entity.LogFilters{
		ProductID:  "productID",
		VersionTag: "versionID",
		From:       fromTime,
		To:         toTime,
		Limit:      100,
	}

	lokiClient := lokiclient.NewClient(cfg)

	logs, err := lokiClient.GetLogs(logFilters)
	require.NoError(t, err)

	assert.Len(t, logs, 5)
	for _, log := range logs {
		assert.NotEmpty(t, log.FormatedLog)
		assert.NotEmpty(t, log.Labels)

		labelsMap := make(map[string]string)
		for _, label := range log.Labels {
			assert.NotEmpty(t, label.Key)
			assert.NotEmpty(t, label.Value)
			labelsMap[label.Key] = label.Value
		}

		assert.Equal(t, "productID", labelsMap["product_id"])
		assert.Equal(t, "versionID", labelsMap["version_tag"])
	}
}

func TestLokiClientGetFullQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		expectedURL := `/loki/api/v1/query_range?end=1701388800000000000&limit=100&query={product_id="productID", version_tag="versionID", workflow_name="workflowName", process_name="processName", request_id="requestID", level="info", logger="[LOGGER]"}&start=1672531200000000000`

		actualQuery, err := url.QueryUnescape(req.URL.String())
		require.NoError(t, err)

		assert.Equal(t, expectedURL, actualQuery)

		rw.Write([]byte(`OK`))
	}))
	defer server.Close()

	cfg := &config.Config{}
	cfg.Loki.Address = server.URL

	fromTime, err := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
	require.NoError(t, err)
	toTime, err := time.Parse(time.RFC3339, "2023-12-01T00:00:00Z")
	require.NoError(t, err)

	logFilters := entity.LogFilters{
		ProductID:    "productID",
		VersionTag:   "versionID",
		From:         fromTime,
		To:           toTime,
		Limit:        100,
		WorkflowName: "workflowName",
		ProcessName:  "processName",
		RequestID:    "requestID",
		Level:        "info",
		Logger:       "[LOGGER]",
	}

	lokiClient := lokiclient.NewClient(cfg)

	_, err = lokiClient.GetLogs(logFilters)
	require.Error(t, err) //unmarhsall error
}
