//go:build unit

package process_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/go-logr/zapr"
	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/process"
	"github.com/konstellation-io/kai/engine/admin-api/mocks"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type registeredProcessMatcher struct {
	expectedRegisteredProcess *entity.RegisteredProcess
}

func newRegisteredProcessMatcher(expectedStreamConfig *entity.RegisteredProcess) *registeredProcessMatcher {
	return &registeredProcessMatcher{
		expectedRegisteredProcess: expectedStreamConfig,
	}
}

func (m registeredProcessMatcher) String() string {
	return fmt.Sprintf("is equal to %v", m.expectedRegisteredProcess)
}

func (m registeredProcessMatcher) Matches(actual interface{}) bool {
	actualCfg, ok := actual.(*entity.RegisteredProcess)
	if !ok {
		return false
	}

	return reflect.DeepEqual(actualCfg, m.expectedRegisteredProcess)
}

type ProcessHandlerTestSuite struct {
	suite.Suite
	ctrl            *gomock.Controller
	processRepo     *mocks.MockProcessRepository
	versionService  *mocks.MockVersionService
	objectStorage   *mocks.MockObjectStorage
	processHandler  *process.Handler
	accessControl   *mocks.MockAccessControl
	productRepo     *mocks.MockProductRepo
	processRegistry *mocks.MockProcessRegistry

	registryHost string
}

const (
	_publicRegistry = "kai"
	_userID         = "userID"
	_userEmail      = "test@email.com"
	_productID      = "productID"
	_version        = "v1.0.0"
	_processName    = "test-process"
	_processType    = entity.ProcessTypeTrigger
	_testFilePath   = "testdata/fake_compressed_process.txt"
)

var (
	user = &entity.User{
		ID:    _userID,
		Roles: []string{"admin"},
		ProductGrants: entity.ProductGrants{
			_productID: {"admin"},
		},
		Email: _userEmail,
	}
)

func TestProcessHandlerSuite(t *testing.T) {
	suite.Run(t, new(ProcessHandlerTestSuite))
}

func (s *ProcessHandlerTestSuite) SetupSuite() {
	logger := zapr.NewLogger(zap.NewNop())
	s.ctrl = gomock.NewController(s.T())
	s.processRepo = mocks.NewMockProcessRepository(s.ctrl)
	s.versionService = mocks.NewMockVersionService(s.ctrl)
	s.objectStorage = mocks.NewMockObjectStorage(s.T())
	s.accessControl = mocks.NewMockAccessControl(s.ctrl)
	s.processRegistry = mocks.NewMockProcessRegistry(s.ctrl)
	s.productRepo = mocks.NewMockProductRepo(s.ctrl)

	s.processHandler = process.NewHandler(
		&process.HandlerParams{
			Logger:            logger,
			VersionService:    s.versionService,
			ProcessRepository: s.processRepo,
			ObjectStorage:     s.objectStorage,
			AccessControl:     s.accessControl,
			ProcessRegistry:   s.processRegistry,
			ProductRepository: s.productRepo,
		},
	)

	s.registryHost = "test.registry"

	viper.Set(config.RegistryHostKey, s.registryHost)
	viper.Set(config.GlobalRegistryKey, _publicRegistry)

	monkey.Patch(time.Now, func() time.Time {
		return time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	})
}

func (s *ProcessHandlerTestSuite) TearDownSuite() {
	monkey.UnpatchAll()
}

func (s *ProcessHandlerTestSuite) TearDownTest() {
	// Clean mockery calls to avoid false positives
	s.ctrl.Finish()
}

func (s *ProcessHandlerTestSuite) getTestProcess(registry, status string, isPublicOpt ...bool) *entity.RegisteredProcess {
	var isPublic bool
	if len(isPublicOpt) > 0 {
		isPublic = isPublicOpt[0]
	}

	return testhelpers.NewRegisteredProcessBuilder(registry).
		WithName(_processName).
		WithVersion(_version).
		WithType(_processType).
		WithOwner(user.Email).
		WithStatus(status).
		WithIsPublic(isPublic).
		Build()
}
