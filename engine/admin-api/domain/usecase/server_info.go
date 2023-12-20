package usecase

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/spf13/viper"
)

type ComponentInfoDTO struct {
	Version string `yaml:"version"`
}

type ServerInfoGetter struct {
	logger        logr.Logger
	accessControl auth.AccessControl
}

func NewServerInfoGetter(logger logr.Logger, accessControl auth.AccessControl) *ServerInfoGetter {
	return &ServerInfoGetter{
		logger:        logger,
		accessControl: accessControl,
	}
}

func (ig *ServerInfoGetter) GetKAIServerInfo(_ context.Context, user *entity.User) (*entity.ServerInfo, error) {
	if err := ig.accessControl.CheckRoleGrants(user, auth.ActViewServerInfo); err != nil {
		return nil, err
	}

	return &entity.ServerInfo{
		Components: ig.collectServerInfo(),
	}, nil
}

func (ig *ServerInfoGetter) collectServerInfo() []entity.ComponentInfo {
	// viper transform map keys in lowercase
	components := viper.GetStringMapString(config.ComponentsKey)

	componentsInfo := make([]entity.ComponentInfo, 0, len(components))
	for component, version := range components {
		componentsInfo = append(componentsInfo, entity.ComponentInfo{
			Name:    component,
			Version: version,
		})
	}

	return componentsInfo
}
