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

type ProcessServiceTestSuite struct {
	suite.Suite
	ctrl           *gomock.Controller
	processRepo    *mocks.MockProcessRepository
	versionService *mocks.MockVersionService
	objectStorage  *mocks.MockObjectStorage
	processService *process.Handler
	accessControl  *mocks.MockAccessControl

	registryHost string
}

const (
	_publicRegistry = "kai"
	userID          = "userID"
	userEmail       = "test@email.com"
	productID       = "productID"
	version         = "v1.0.0"
	processName     = "test-process"
	processType     = entity.ProcessTypeTrigger
	testFileAddr    = "testdata/fake_compressed_process.txt"
)

var (
	user = &entity.User{
		ID:    userID,
		Roles: []string{"admin"},
		ProductGrants: entity.ProductGrants{
			productID: {"admin"},
		},
		Email: userEmail,
	}
)

func TestProcessTestSuite(t *testing.T) {
	suite.Run(t, new(ProcessServiceTestSuite))
}

func (s *ProcessServiceTestSuite) SetupSuite() {
	logger := zapr.NewLogger(zap.NewNop())
	s.ctrl = gomock.NewController(s.T())
	s.processRepo = mocks.NewMockProcessRepository(s.ctrl)
	s.versionService = mocks.NewMockVersionService(s.ctrl)
	s.objectStorage = mocks.NewMockObjectStorage(s.T())
	s.accessControl = mocks.NewMockAccessControl(s.ctrl)
	s.processService = process.NewHandler(logger, s.versionService, s.processRepo, s.objectStorage, s.accessControl)

	s.registryHost = "test.registry"

	viper.Set(config.RegistryHostKey, s.registryHost)
	viper.Set(config.GlobalRegistryKey, _publicRegistry)

	monkey.Patch(time.Now, func() time.Time {
		return time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	})
}

func (s *ProcessServiceTestSuite) TearDownSuite() {
	monkey.UnpatchAll()
}

func (s *ProcessServiceTestSuite) TearDownTest() {
	// Clean mockery calls to avoid false positives
	s.objectStorage.ExpectedCalls = nil
}

func (s *ProcessServiceTestSuite) getTestProcess(registry, status string, isPublicOpt ...bool) *entity.RegisteredProcess {
	var isPublic bool
	if len(isPublicOpt) > 0 {
		isPublic = isPublicOpt[0]
	}

	return testhelpers.NewRegisteredProcessBuilder(registry).
		WithName(processName).
		WithVersion(version).
		WithType(processType).
		WithOwner(user.Email).
		WithStatus(status).
		WithIsPublic(isPublic).
		Build()
}
