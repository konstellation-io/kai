//go:build unit

package usecase_test

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/go-logr/zapr"
	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase"
	"github.com/konstellation-io/kai/engine/admin-api/mocks"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type processRegistryMatcher struct {
	expectedprocessRegistry *entity.ProcessRegistry
}

func newprocessRegistryMatcher(expectedStreamConfig *entity.ProcessRegistry) *processRegistryMatcher {
	return &processRegistryMatcher{
		expectedprocessRegistry: expectedStreamConfig,
	}
}

func (m processRegistryMatcher) String() string {
	return fmt.Sprintf("is equal to %v", m.expectedprocessRegistry)
}

func (m processRegistryMatcher) Matches(actual interface{}) bool {
	actualCfg, ok := actual.(*entity.ProcessRegistry)
	if !ok {
		return false
	}

	return reflect.DeepEqual(actualCfg, m.expectedprocessRegistry)
}

type ProcessServiceTestSuite struct {
	suite.Suite
	ctrl                *gomock.Controller
	processRegistryRepo *mocks.MockProcessRegistryRepo
	k8sService          *mocks.MockK8sService
	processInteractor   *usecase.ProcessService
}

const (
	userID    = "userID"
	productID = "productID"
)

func TestProcessTestSuite(t *testing.T) {
	suite.Run(t, new(ProcessServiceTestSuite))
}

func (s *ProcessServiceTestSuite) SetupSuite() {
	s.ctrl = gomock.NewController(s.T())

	logger := zapr.NewLogger(zap.NewNop())
	s.processRegistryRepo = mocks.NewMockProcessRegistryRepo(s.ctrl)
	s.k8sService = mocks.NewMockK8sService(s.ctrl)
	s.processInteractor = usecase.NewProcessService(logger, s.k8sService, s.processRegistryRepo)
}

func (s *ProcessServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *ProcessServiceTestSuite) TestRegisterProcess() {
	monkey.Patch(time.Now, func() time.Time {
		return time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	})
	defer monkey.UnpatchAll()

	ctx := context.Background()
	user := &entity.User{
		ID:    userID,
		Roles: []string{"admin"},
		ProductGrants: entity.ProductGrants{
			productID: {"admin"},
		},
	}
	version := "v1.0.0"
	process := "test-process"
	processType := entity.ProcessTypeTrigger.String()

	testFile, err := os.Open("testdata/fake_compressed_process.txt")
	s.Require().NoError(err)
	expectedBytes, err := os.ReadFile("testdata/fake_compressed_process.txt")
	s.Require().NoError(err)

	mockedRef := "mocked-ref"
	expectedRegisteredProcess := &entity.ProcessRegistry{
		ID:         mockedRef,
		Name:       process,
		Version:    version,
		Type:       processType,
		UploadDate: time.Now(),
		Owner:      userID,
	}
	customMatcher := newprocessRegistryMatcher(expectedRegisteredProcess)

	s.k8sService.EXPECT().RegisterProcess(ctx, productID, version, process, expectedBytes).Return(mockedRef, nil)
	s.processRegistryRepo.EXPECT().Create(productID, customMatcher).Return(nil, nil)

	returnedRef, err := s.processInteractor.RegisterProcess(
		ctx, user, productID, version, process, processType, testFile,
	)
	s.Require().NoError(err)

	s.Equal(mockedRef, returnedRef)
}
