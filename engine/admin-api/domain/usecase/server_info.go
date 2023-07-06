package usecase

import (
	"context"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/spf13/viper"
)

//var ErrServerVersionNotFound = errors.New("server version not found")

func GetKAIServerInfo(_ context.Context) (*entity.ServerInfo, error) {
	return &entity.ServerInfo{
		Components: []entity.ComponentInfo{
			{
				Component: "AIO",
				Version:   viper.GetString(config.ReleaseVersionKey),
				Status:    entity.ComponentStatusOK,
			},
		},
	}, nil
}

func GetKAIComponentsInfo(_ context.Context) (*entity.ServerInfo, error) {
	return &entity.ServerInfo{
		Components: []entity.ComponentInfo{
			{
				Component: "AIO",
				Version:   viper.GetString(config.ReleaseVersionKey),
				Status:    entity.ComponentStatusOK,
			},
			{
				Component: "MongoDB",
				Version:   viper.GetString(config.ReleaseVersionKey),
				Status:    entity.ComponentStatusOK,
			},
		},
	}, nil
}
