package version

import (
	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/admin-api/domain/repository"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase"
)

// Handler contains app logic about Version entities.
type Handler struct {
	logger                 logr.Logger
	versionRepo            repository.VersionRepo
	productRepo            repository.ProductRepo
	k8sService             service.VersionService
	natsManagerService     service.NatsManagerService
	userActivityInteractor usecase.UserActivityInteracter
	accessControl          auth.AccessControl
	dashboardService       service.DashboardService
	processLogRepo         repository.ProcessLogRepository
}

// NewHandler creates a new interactor.
func NewHandler(
	logger logr.Logger,
	versionRepo repository.VersionRepo,
	productRepo repository.ProductRepo,
	k8sService service.VersionService,
	natsManagerService service.NatsManagerService,
	userActivityInteractor usecase.UserActivityInteracter,
	accessControl auth.AccessControl,
	dashboardService service.DashboardService,
	processLogRepo repository.ProcessLogRepository,
) *Handler {
	return &Handler{
		logger,
		versionRepo,
		productRepo,
		k8sService,
		natsManagerService,
		userActivityInteractor,
		accessControl,
		dashboardService,
		processLogRepo,
	}
}
