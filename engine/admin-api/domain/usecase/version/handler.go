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

type HanlderParams struct {
	Logger                 logr.Logger
	VersionRepo            repository.VersionRepo
	ProductRepo            repository.ProductRepo
	K8sService             service.VersionService
	NatsManagerService     service.NatsManagerService
	UserActivityInteractor usecase.UserActivityInteracter
	AccessControl          auth.AccessControl
	DashboardService       service.DashboardService
	ProcessLogRepo         repository.ProcessLogRepository
}

// NewHandler creates a new interactor.
func NewHandler(p HanlderParams) *Handler {
	return &Handler{
		p.Logger,
		p.VersionRepo,
		p.ProductRepo,
		p.K8sService,
		p.NatsManagerService,
		p.UserActivityInteractor,
		p.AccessControl,
		p.DashboardService,
		p.ProcessLogRepo,
	}
}
