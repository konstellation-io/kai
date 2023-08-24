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
	userID       = "userID"
	productID    = "productID"
	processID    = productID + "_" + processName + ":" + version
	version      = "v1.0.0"
	processName  = "test-process"
	processType  = "trigger"
	testFileAddr = "testdata/fake_compressed_process.txt"
)

var (
	user = &entity.User{
		ID:    userID,
		Roles: []string{"admin"},
		ProductGrants: entity.ProductGrants{
			productID: {"admin"},
		},
	}
)

func TestProcessTestSuite(t *testing.T) {
	suite.Run(t, new(ProcessServiceTestSuite))
}

func (s *ProcessServiceTestSuite) SetupSuite() {
	logger := zapr.NewLogger(zap.NewNop())
	s.ctrl = gomock.NewController(s.T())
	s.processRegistryRepo = mocks.NewMockProcessRegistryRepo(s.ctrl)
	s.k8sService = mocks.NewMockK8sService(s.ctrl)
	s.processInteractor = usecase.NewProcessService(logger, s.k8sService, s.processRegistryRepo)

	monkey.Patch(time.Now, func() time.Time {
		return time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	})
}

func (s *ProcessServiceTestSuite) TearDownSuite() {
	monkey.UnpatchAll()
}

func (s *ProcessServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *ProcessServiceTestSuite) TestRegisterProcess() {
	ctx := context.Background()

	testFile, err := os.Open(testFileAddr)
	s.Require().NoError(err)
	expectedBytes, err := os.ReadFile(testFileAddr)
	s.Require().NoError(err)

	mockedRef := fmt.Sprintf("%s/%s_%s:%s", "kai-local-registry:5000", productID, processName, version)
	expectedRegisteredProcess := &entity.ProcessRegistry{
		ID:         processID,
		Name:       processName,
		Version:    version,
		Type:       processType,
		Image:      mockedRef,
		UploadDate: time.Now(),
		Owner:      userID,
	}
	customMatcher := newprocessRegistryMatcher(expectedRegisteredProcess)

	s.k8sService.EXPECT().RegisterProcess(ctx, productID, version, processName, expectedBytes).Return(mockedRef, nil)
	s.processRegistryRepo.EXPECT().Create(productID, customMatcher).Return(nil, nil)

	returnedID, err := s.processInteractor.RegisterProcess(
		ctx, user, productID, version, processName, processType, testFile,
	)
	s.Require().NoError(err)

	s.Equal(processID, returnedID)
}

func (s *ProcessServiceTestSuite) TestRegisterProcessNoFileError() {
	ctx := context.Background()

	testFile, err := os.Open("no-file")
	s.Require().Error(err)

	returnedRef, err := s.processInteractor.RegisterProcess(
		ctx, user, productID, version, processName, processType, testFile,
	)
	s.Require().Error(err)

	s.Empty(returnedRef)
}

func (s *ProcessServiceTestSuite) TestRegisterProcessK8sManagerError() {
	ctx := context.Background()

	testFile, err := os.Open(testFileAddr)
	s.Require().NoError(err)
	expectedBytes, err := os.ReadFile(testFileAddr)
	s.Require().NoError(err)

	s.k8sService.EXPECT().
		RegisterProcess(ctx, productID, version, processName, expectedBytes).
		Return("", fmt.Errorf("mocked error"))

	returnedRef, err := s.processInteractor.RegisterProcess(
		ctx, user, productID, version, processName, processType, testFile,
	)
	s.Require().Error(err)

	s.Empty(returnedRef)
}

func (s *ProcessServiceTestSuite) TestRegisterProcessRepositoryError() {
	ctx := context.Background()

	testFile, err := os.Open(testFileAddr)
	s.Require().NoError(err)
	expectedBytes, err := os.ReadFile(testFileAddr)
	s.Require().NoError(err)

	mockedRef := fmt.Sprintf("%s/%s_%s:%s", "kai-local-registry:5000", productID, processName, version)
	expectedRegisteredProcess := &entity.ProcessRegistry{
		ID:         processID,
		Name:       processName,
		Version:    version,
		Type:       processType,
		Image:      mockedRef,
		UploadDate: time.Now(),
		Owner:      userID,
	}
	customMatcher := newprocessRegistryMatcher(expectedRegisteredProcess)

	s.k8sService.EXPECT().RegisterProcess(ctx, productID, version, processName, expectedBytes).Return(mockedRef, nil)
	s.processRegistryRepo.EXPECT().Create(productID, customMatcher).Return(nil, fmt.Errorf("mocked error"))

	returnedRef, err := s.processInteractor.RegisterProcess(
		ctx, user, productID, version, processName, processType, testFile,
	)
	s.Require().Error(err)

	s.Empty(returnedRef)
}

func (s *ProcessServiceTestSuite) TestListByProductWithTypeFilter() {
	ctx := context.Background()

	typeFilter := "trigger"
	expectedProcessRegistry := []*entity.ProcessRegistry{
		{
			ID:         processID,
			Name:       processName,
			Version:    version,
			Type:       processType,
			Image:      "image",
			UploadDate: time.Now(),
			Owner:      userID,
		},
	}

	s.processRegistryRepo.EXPECT().ListByProductWithTypeFilter(ctx, productID, typeFilter).Return(expectedProcessRegistry, nil)

	returnedProcessRegistry, err := s.processInteractor.ListByProductWithTypeFilter(ctx, productID, typeFilter)
	s.Require().NoError(err)

	s.Equal(expectedProcessRegistry, returnedProcessRegistry)
}

func (s *ProcessServiceTestSuite) TestListByProductWithTypeFilterInvalidFilter() {
	ctx := context.Background()

	typeFilter := "invlaid"

	_, err := s.processInteractor.ListByProductWithTypeFilter(ctx, productID, typeFilter)
	s.Require().Error(err)
}
