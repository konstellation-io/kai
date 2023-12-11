package process

import (
	"errors"
	"io"

	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/repository"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
)

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/${GOFILE} -package=mocks

var (
	ErrInvalidProcessType        = errors.New("invalid process type")
	ErrRegisteredProcessNotFound = errors.New("registered process not found")
	ErrProcessAlreadyRegistered  = errors.New("process already registered")
)

type Service struct {
	logger            logr.Logger
	processRepository repository.ProcessRepository
	versionService    service.VersionService
	objectStorage     repository.ObjectStorage
	accessControl     auth.AccessControl
}

type RegisterProcessOpts struct {
	Product     string
	Version     string
	Process     string
	ProcessType entity.ProcessType
	IsPublic    bool
	Sources     io.Reader
}

var (
	ErrMissingProductInParams = errors.New("missing product in params")
	ErrMissingVersionInParams = errors.New("missing version in params")
	ErrMissingSourcesInParams = errors.New("missing sources in params")
)

func (o RegisterProcessOpts) Validate() error {
	if o.Product == "" && !o.IsPublic {
		return ErrMissingProductInParams
	}

	if o.Version == "" {
		return ErrMissingVersionInParams
	}

	if err := o.ProcessType.Validate(); err != nil {
		return err
	}

	if o.Sources == nil {
		return ErrMissingSourcesInParams
	}

	return nil
}

func NewProcessService(
	logger logr.Logger,
	k8sService service.VersionService,
	processRepository repository.ProcessRepository,
	objectStorage repository.ObjectStorage,
	accessControl auth.AccessControl,
) *Service {
	return &Service{
		logger:            logger,
		versionService:    k8sService,
		processRepository: processRepository,
		objectStorage:     objectStorage,
		accessControl:     accessControl,
	}
}
