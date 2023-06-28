//go:build unit

package config_test

import (
	"os"
	"testing"

	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
)

const configFile = "../../../config.yml"

func TestInitConfig_OverrideWithEnv(t *testing.T) {
	err := os.Setenv("KAI_SERVER_PORT", "9090")
	require.NoError(t, err)

	err = config.Init(configFile)
	require.NoError(t, err)

	port := viper.GetInt("server.port")
	assert.Equal(t, port, 9090)
}
