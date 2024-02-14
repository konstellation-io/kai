package process

import (
	"errors"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/repository"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/spf13/viper"
)

var (
	ErrRegisteredProcessNotFound = errors.New("registered process not found")
	ErrProcessAlreadyRegistered  = errors.New("process already registered")
	ErrMissingProductInParams    = errors.New("missing product in params")
	ErrMissingVersionInParams    = errors.New("missing version in params")
	ErrMissingProcessInParams    = errors.New("missing process in params")
	ErrMissingSourcesInParams    = errors.New("missing sources in params")
	ErrIsPublicAndHasProduct     = errors.New("a process cannot be public and come from a product at the same time")
)

type Handler struct {
	logger            logr.Logger
	processRepository repository.ProcessRepository
	versionService    service.VersionService
	objectStorage     repository.ObjectStorage
	accessControl     auth.AccessControl
	processRegistry   service.ProcessRegistry
}

type HandlerParams struct {
	Logger            logr.Logger
	VersionService    service.VersionService
	ProcessRepository repository.ProcessRepository
	ObjectStorage     repository.ObjectStorage
	AccessControl     auth.AccessControl
	ProcessRegistry   service.ProcessRegistry
}

func NewHandler(
	params *HandlerParams,
) *Handler {
	return &Handler{
		logger:            params.Logger,
		processRepository: params.ProcessRepository,
		versionService:    params.VersionService,
		objectStorage:     params.ObjectStorage,
		accessControl:     params.AccessControl,
		processRegistry:   params.ProcessRegistry,
	}
}

func (ps *Handler) getProcessID(scope, process, version string) string {
	return fmt.Sprintf("%s_%s:%s", scope, process, version)
}

func (ps *Handler) getImageName(scope, process string) string {
	return fmt.Sprintf("%s_%s", scope, process)
}

func (ps *Handler) getProcessImage(processID string) string {
	return fmt.Sprintf("%s/%s", viper.GetString(config.RegistryHostKey), processID)
}

func (ps *Handler) getProcessRegisterScope(isPublic bool, product string) string {
	if isPublic {
		return viper.GetString(config.GlobalRegistryKey)
	}

	return product
}
