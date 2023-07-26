package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-logr/logr/testr"
	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/auth"
	"github.com/konstellation-io/kai/engine/admin-api/mocks"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetServerInfo(t *testing.T) {
	var (
		logger        = testr.NewWithOptions(t, testr.Options{Verbosity: -1})
		ctrl          = gomock.NewController(t)
		accessControl = mocks.NewMockAccessControl(ctrl)

		user = &entity.User{}
		ctx  = context.Background()
	)

	err := initTestConfig("testdata/versions.yaml")
	require.NoError(t, err)

	expectedServersInfo := &entity.ServerInfo{
		Components: []entity.ComponentInfo{
			{
				Name:    "aio",
				Version: "latest",
			},
			{
				Name:    "mongodb",
				Version: "latest",
			},
		},
	}

	accessControl.EXPECT().CheckRoleGrants(user, auth.ActViewServerInfo).Return(nil)

	serverInfoCollector := usecase.NewServerInfoGetter(logger, accessControl)
	serverInfo, err := serverInfoCollector.GetKAIServerInfo(ctx, user)
	assert.NoError(t, err)

	for _, component := range serverInfo.Components {
		assert.Contains(t, expectedServersInfo.Components, component)
	}
}

func TestGetServerInfo_Unauthorized(t *testing.T) {
	var (
		logger        = testr.NewWithOptions(t, testr.Options{Verbosity: -1})
		ctrl          = gomock.NewController(t)
		accessControl = mocks.NewMockAccessControl(ctrl)

		user = &entity.User{}
		ctx  = context.Background()
	)

	expectedError := errors.New("unauthorized")
	accessControl.EXPECT().CheckRoleGrants(user, auth.ActViewServerInfo).Return(expectedError)

	serverInfoCollector := usecase.NewServerInfoGetter(logger, accessControl)

	serverInfo, err := serverInfoCollector.GetKAIServerInfo(ctx, user)
	assert.ErrorIs(t, err, expectedError)
	assert.Nil(t, serverInfo)
}

func initTestConfig(configPath string) error {
	viper.Set(config.CfgFilePathKey, configPath)

	return config.InitConfig()
}
