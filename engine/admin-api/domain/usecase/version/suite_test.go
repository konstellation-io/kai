//go:build unit

package version_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/go-logr/zapr"
	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/version"
	"github.com/konstellation-io/kai/engine/admin-api/mocks"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type versionMatcher struct {
	expectedVersion *entity.Version
}

func newVersionMatcher(expectedVersion *entity.Version) *versionMatcher {
	return &versionMatcher{
		expectedVersion: expectedVersion,
	}
}

func (m versionMatcher) String() string {
	return fmt.Sprintf("is equal to %v", m.expectedVersion)
}

func (m versionMatcher) Matches(actual interface{}) bool {
	actualCfg, ok := actual.(*entity.Version)
	if !ok {
		return false
	}

	return reflect.DeepEqual(actualCfg, m.expectedVersion)
}

type VersionUsecaseTestSuite struct {
	suite.Suite
	handler *version.Handler

	ctrl                   *gomock.Controller
	versionRepo            *mocks.MockVersionRepo
	productRepo            *mocks.MockProductRepo
	versionService         *mocks.MockVersionService
	natsManagerService     *mocks.MockNatsManagerService
	userActivityInteractor *mocks.MockUserActivityInteracter
	accessControl          *mocks.MockAccessControl
	dashboardService       *mocks.MockDashboardService
	processLogRepo         *mocks.MockProcessLogRepository
}

const (
	userID     = "userID"
	productID  = "productID"
	versionID  = "versionID"
	versionTag = "v1.0.0"
)

func (s *VersionUsecaseTestSuite) getTestUser() *entity.User {
	return &entity.User{
		ID:    userID,
		Roles: []string{"admin"},
		ProductGrants: entity.ProductGrants{
			productID: {"admin"},
		},
	}
}

func TestVersionUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(VersionUsecaseTestSuite))
}

func (s *VersionUsecaseTestSuite) SetupSuite() {
	logger := zapr.NewLogger(zap.NewNop())

	s.ctrl = gomock.NewController(s.T())
	s.versionRepo = mocks.NewMockVersionRepo(s.ctrl)
	s.productRepo = mocks.NewMockProductRepo(s.ctrl)
	s.versionService = mocks.NewMockVersionService(s.ctrl)
	s.natsManagerService = mocks.NewMockNatsManagerService(s.ctrl)
	s.userActivityInteractor = mocks.NewMockUserActivityInteracter(s.ctrl)
	s.accessControl = mocks.NewMockAccessControl(s.ctrl)
	s.dashboardService = mocks.NewMockDashboardService(s.ctrl)
	s.processLogRepo = mocks.NewMockProcessLogRepository(s.ctrl)

	s.handler = version.NewHandler(version.HandlerParams{
		Logger:                 logger,
		VersionRepo:            s.versionRepo,
		ProductRepo:            s.productRepo,
		K8sService:             s.versionService,
		NatsManagerService:     s.natsManagerService,
		UserActivityInteractor: s.userActivityInteractor,
		AccessControl:          s.accessControl,
		DashboardService:       s.dashboardService,
		ProcessLogRepo:         s.processLogRepo,
	})
}

func (s *VersionUsecaseTestSuite) TearDownTest() {
	s.ctrl.Finish()
}
