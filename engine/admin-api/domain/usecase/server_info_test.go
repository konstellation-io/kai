package usecase_test

import (
	"context"
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
)

func TestGetServerInfo(t *testing.T) {
	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})
	ctrl := gomock.NewController(t)
	accessControl := mocks.NewMockAccessControl(ctrl)
	user := &entity.User{}
	ctx := context.Background()
	viper.Set(config.VersionsFilePathKey, "testdata/versions.yaml")

	expectedServersInfo := &entity.ServerInfo{
		Components: []entity.ComponentInfo{
			{
				Name:    "aio",
				Version: "latest",
			},
			{
				Name:    "mongoDB",
				Version: "latest",
			},
		},
		Status: "OK",
	}

	accessControl.EXPECT().CheckProductGrants(user, "", auth.ActViewServerInfo).Return(nil)

	serverInfoCollector := usecase.NewServerInfoGetter(logger, accessControl)
	serverInfo, err := serverInfoCollector.GetKAIServerInfo(ctx, user)
	assert.NoError(t, err)
	assert.Equal(t, serverInfo, expectedServersInfo)
}
