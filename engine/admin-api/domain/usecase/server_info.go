package usecase

import (
	"context"
	"fmt"
	"os"

	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/auth"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

//var ErrServerVersionNotFound = errors.New("server version not found")

type ServerInfo struct {
	Components []ComponentInfo
	Status     string
}

type ComponentInfo struct {
	Name    string
	Version string
}

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
	if err := ig.accessControl.CheckProductGrants(user, "", auth.ActViewServerInfo); err != nil {
		return nil, err
	}

	componentsInfo, err := ig.collectServerInfo()
	if err != nil {
		return nil, err
	}

	return &entity.ServerInfo{
		Components: componentsInfo,
		Status:     "OK",
	}, nil
}

func (ig *ServerInfoGetter) collectServerInfo() ([]entity.ComponentInfo, error) {
	versionsFilePath := viper.GetString(config.VersionsFilePathKey)
	file, err := os.ReadFile(versionsFilePath)
	if err != nil {
		return nil, fmt.Errorf("read %q file: %w", versionsFilePath, err)
	}

	var componentsInfoDTO map[string]ComponentInfoDTO

	err = yaml.Unmarshal(file, &componentsInfoDTO)
	if err != nil {
		return nil, fmt.Errorf("parsing versions file: %w", err)
	}

	componentsInfo := make([]entity.ComponentInfo, 0, len(componentsInfoDTO))
	for component, info := range componentsInfoDTO {
		componentsInfo = append(componentsInfo, entity.ComponentInfo{
			Name:    component,
			Version: info.Version,
		})
	}

	return componentsInfo, nil
}
