//go:build unit

package http_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-logr/logr/testr"
	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	httpapp "github.com/konstellation-io/kai/engine/admin-api/delivery/http"
	"github.com/konstellation-io/kai/engine/admin-api/mocks"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestServerCall_GraphqlFailsIfUnauthorized(t *testing.T) {
	viper.Set(config.ApplicationPortKey, 4000)
	ctrl := gomock.NewController(t)
	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})

	gqlController := mocks.NewMockGraphQL(ctrl)

	app := httpapp.NewApp(
		logger,
		gqlController,
	)

	go app.Start()

	url := fmt.Sprintf("http://localhost:%s/graphql", viper.GetString(config.ApplicationPortKey))

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
	require.NoError(t, err)

	req.Header.Add("Authorization", "Bearer invalid-token")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer res.Body.Close()
	require.Equal(t, res.StatusCode, http.StatusUnauthorized)
}
