//go:build unit

package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/mocks"
	"github.com/stretchr/testify/assert"
)

type versionDashboardsSuite struct {
	ctrl              *gomock.Controller
	versionInteractor *VersionInteractor
	mocks             versionDashboardsSuiteMocks
}

type versionDashboardsSuiteMocks struct {
	logger           logr.Logger
	dashboardService *mocks.MockDashboardService
}

func newVersionDashboardsSuite(t *testing.T) *versionDashboardsSuite {
	ctrl := gomock.NewController(t)

	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})
	dashboardService := mocks.NewMockDashboardService(ctrl)
	versionRepo := mocks.NewMockVersionRepo(ctrl)
	runtimeRepo := mocks.NewMockProductRepo(ctrl)
	versionService := mocks.NewMockVersionService(ctrl)
	natsManagerService := mocks.NewMockNatsManagerService(ctrl)
	userActivityInteractor := mocks.NewMockUserActivityInteracter(ctrl)
	accessControl := mocks.NewMockAccessControl(ctrl)
	processLogRepo := mocks.NewMockProcessLogRepository(ctrl)

	versionInteractor := NewVersionInteractor(logger, versionRepo, runtimeRepo, versionService,
		natsManagerService, userActivityInteractor, accessControl, dashboardService, processLogRepo)

	return &versionDashboardsSuite{ctrl: ctrl,
		versionInteractor: versionInteractor,
		mocks: versionDashboardsSuiteMocks{
			logger:           logger,
			dashboardService: dashboardService,
		},
	}
}

func TestStoreDashboard(t *testing.T) {
	s := newVersionDashboardsSuite(t)
	defer s.ctrl.Finish()

	version := "test"
	runtimeID := "test-runtime-id"
	dashboardsFolder := "../../testdata/dashboards"
	dashboardPath := fmt.Sprintf("%s/models.json", dashboardsFolder)
	s.mocks.dashboardService.EXPECT().Create(context.Background(), runtimeID, version, dashboardPath).Return(nil)

	err := s.versionInteractor.storeDashboards(context.Background(), dashboardsFolder, runtimeID, version)
	assert.Nil(t, err)
}

func TestStoreDashboardWrongFolderPath(t *testing.T) {
	s := newVersionDashboardsSuite(t)
	defer s.ctrl.Finish()

	version := "test"
	runtimeID := "test-runtime-id"
	dashboardsFolder := "../../testdata/dashboard"

	err := s.versionInteractor.storeDashboards(context.Background(), dashboardsFolder, runtimeID, version)
	assert.NotNil(t, err)
	assert.Contains(t, err[0].Error(), "error listing dashboards files: open ../../testdata/dashboard: no such file or directory")
}
