package config_test

import (
	"os"
	"testing"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitConfig_ValueSetFromConfigFile(t *testing.T) {
	expectedValue := "kai"

	viper.Set(config.CfgFilePathKey, "../../config.yml")

	err := config.InitConfig()
	require.NoError(t, err)

	actualValue := viper.GetString(config.MongoDBKaiDatabaseKey)

	assert.Equal(t, expectedValue, actualValue)

	viper.Reset()
}

func TestInitConfig_ValueSetFromEnv(t *testing.T) {
	expectedValue := "test"

	err := os.Setenv("KAI_MONGODB_DATABASE", expectedValue)
	require.NoError(t, err)

	viper.Set(config.CfgFilePathKey, "../../config.yml")

	err = config.InitConfig()
	require.NoError(t, err)

	actualValue := viper.GetString(config.MongoDBKaiDatabaseKey)

	assert.Equal(t, expectedValue, actualValue)
	viper.Reset()
}

func TestInitConfig_DefaultValue(t *testing.T) {
	expectedValue := "predictionsIdx"

	viper.Set(config.CfgFilePathKey, "../../config.yml")

	err := config.InitConfig()
	require.NoError(t, err)

	actualValue := viper.GetString(config.RedisPredictionsIndexKey)

	assert.Equal(t, expectedValue, actualValue)
	viper.Reset()
}

func TestInitConfig_DefaultValueOverridenByEnv(t *testing.T) {
	expectedValue := "notDefaultPredictionsIndex"

	err := os.Setenv("KAI_REDIS_PREDICTIONS_INDEX", expectedValue)
	require.NoError(t, err)

	viper.Set(config.CfgFilePathKey, "../../config.yml")

	err = config.InitConfig()
	require.NoError(t, err)

	actualValue := viper.GetString(config.RedisPredictionsIndexKey)

	assert.Equal(t, expectedValue, actualValue)
	viper.Reset()
}

func TestInitConfig_IgnoreConfigFileIfDoesNotExist(t *testing.T) {
	viper.Set(config.CfgFilePathKey, "notexists.yaml")
	err := config.InitConfig()
	assert.NoError(t, err)
}

func TestInitConfig_FailsIfInvalidConfigFile(t *testing.T) {
	viper.Set(config.CfgFilePathKey, "notvalid")
	err := config.InitConfig()
	require.Error(t, err)
}
