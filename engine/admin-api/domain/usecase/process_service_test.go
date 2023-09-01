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
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase"
	"github.com/konstellation-io/kai/engine/admin-api/mocks"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type registeredProcessMatcher struct {
	expectedRegisteredProcess *entity.RegisteredProcess
}

func newregisteredProcessMatcher(expectedStreamConfig *entity.RegisteredProcess) *registeredProcessMatcher {
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
	ctrl              *gomock.Controller
	processRepo       *mocks.MockProcessRepository
	versionService    *mocks.MockVersionService
	processInteractor *usecase.ProcessService

	registryHost string
}

const (
	userID       = "userID"
	productID    = "productID"
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

func (s *ProcessServiceTestSuite) getRegisteredProcessImage(registryHost string, expectedProcessID string) string {
	return fmt.Sprintf("%s/%s", registryHost, expectedProcessID)
}

func (s *ProcessServiceTestSuite) getRegisteredProcessID(product string, process string, version string) string {
	return fmt.Sprintf("%s_%s:%s", product, process, version)
}

func (s *ProcessServiceTestSuite) getTestProcess() *entity.RegisteredProcess {
	var (
		expectedProcessID    = s.getRegisteredProcessID(productID, processName, version)
		expectedProcessImage = s.getRegisteredProcessImage(s.registryHost, expectedProcessID)
	)

	return &entity.RegisteredProcess{
		ID:         expectedProcessID,
		Name:       processName,
		Version:    version,
		Type:       processType,
		Image:      expectedProcessImage,
		UploadDate: time.Now().Truncate(time.Millisecond).UTC(),
		Owner:      userID,
	}
}

func TestProcessTestSuite(t *testing.T) {
	suite.Run(t, new(ProcessServiceTestSuite))
}

func (s *ProcessServiceTestSuite) SetupSuite() {
	logger := zapr.NewLogger(zap.NewNop())
	s.ctrl = gomock.NewController(s.T())
	s.processRepo = mocks.NewMockProcessRepository(s.ctrl)
	s.versionService = mocks.NewMockVersionService(s.ctrl)
	s.processInteractor = usecase.NewProcessService(logger, s.versionService, s.processRepo)

	s.registryHost = "test.registry"

	viper.Set(config.RegistryURLKey, "http://"+s.registryHost)

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

	expectedRegisteredProcess := s.getTestProcess()
	expectedRegisteredProcess.Status = entity.RegisterProcessStatusCreating
	customMatcher := newregisteredProcessMatcher(expectedRegisteredProcess)

	expectedUpdatedProcess := s.getTestProcess()
	expectedUpdatedProcess.Status = entity.RegisterProcessStatusCreated
	customMatcherUpdate := newregisteredProcessMatcher(expectedUpdatedProcess)

	s.processRepo.EXPECT().GetByID(ctx, productID, expectedRegisteredProcess.ID).Return(nil, usecase.ErrRegisteredProcessNotFound)
	s.processRepo.EXPECT().Create(productID, customMatcher).Return(nil, nil)
	s.versionService.EXPECT().RegisterProcess(gomock.Any(), expectedRegisteredProcess.ID, expectedRegisteredProcess.Image, expectedBytes).Return("", nil)
	s.processRepo.EXPECT().Update(gomock.Any(), productID, customMatcherUpdate).Return(nil)

	returnedID, notifyCh, err := s.processInteractor.RegisterProcess(
		ctx, user, productID, version, processName, processType, testFile,
	)
	s.Require().NoError(err)

	s.Equal(expectedRegisteredProcess, returnedID)

	registeredProcess := <-notifyCh
	s.Equal(entity.RegisterProcessStatusCreated, registeredProcess.Status)
}

func (s *ProcessServiceTestSuite) TestRegisterProcessNoFileError() {
	ctx := context.Background()

	testFile, err := os.Open("no-file")
	s.Require().Error(err)

	expectedRegisteredProcess := s.getTestProcess()
	expectedRegisteredProcess.Status = entity.RegisterProcessStatusCreating
	customMatcher := newregisteredProcessMatcher(expectedRegisteredProcess)

	expectedUpdatedProcess := s.getTestProcess()
	expectedUpdatedProcess.Status = entity.RegisterProcessStatusFailed
	expectedUpdatedProcess.Logs = "copying temp file for version: invalid argument"
	customMatcherUpdate := newregisteredProcessMatcher(expectedUpdatedProcess)

	s.processRepo.EXPECT().GetByID(ctx, productID, expectedRegisteredProcess.ID).Return(nil, usecase.ErrRegisteredProcessNotFound)
	s.processRepo.EXPECT().Create(productID, customMatcher).Return(nil, nil)
	s.processRepo.EXPECT().Update(gomock.Any(), productID, customMatcherUpdate).Return(nil)

	returnedRef, notifyCh, err := s.processInteractor.RegisterProcess(
		ctx, user, productID, version, processName, processType, testFile,
	)
	s.Require().NoError(err)

	s.Equal(expectedRegisteredProcess, returnedRef)

	registeredProcess := <-notifyCh
	s.Equal(entity.RegisterProcessStatusFailed, registeredProcess.Status)
}

func (s *ProcessServiceTestSuite) TestRegisterProcess_K8sServiceError() {
	ctx := context.Background()

	testFile, err := os.Open(testFileAddr)
	s.Require().NoError(err)
	expectedBytes, err := os.ReadFile(testFileAddr)
	s.Require().NoError(err)

	expectedRegisteredProcess := s.getTestProcess()
	expectedRegisteredProcess.Status = entity.RegisterProcessStatusCreating
	customMatcher := newregisteredProcessMatcher(expectedRegisteredProcess)

	expectedUpdatedProcess := s.getTestProcess()
	expectedUpdatedProcess.Status = entity.RegisterProcessStatusFailed
	expectedUpdatedProcess.Logs = "registering process: mocked error"
	customMatcherUpdate := newregisteredProcessMatcher(expectedUpdatedProcess)

	s.processRepo.EXPECT().GetByID(ctx, productID, expectedRegisteredProcess.ID).Return(nil, usecase.ErrRegisteredProcessNotFound)
	s.processRepo.EXPECT().Create(productID, customMatcher).Return(nil, nil)
	s.versionService.EXPECT().
		RegisterProcess(gomock.Any(), expectedRegisteredProcess.ID, expectedRegisteredProcess.Image, expectedBytes).
		Return("", fmt.Errorf("mocked error"))
	s.processRepo.EXPECT().Update(gomock.Any(), productID, customMatcherUpdate).Return(nil)

	returnedRef, notifyCh, err := s.processInteractor.RegisterProcess(
		ctx, user, productID, version, processName, processType, testFile,
	)
	s.Require().NoError(err)

	s.Equal(expectedRegisteredProcess, returnedRef)

	registeredProcess := <-notifyCh
	s.Equal(entity.RegisterProcessStatusFailed, registeredProcess.Status)
}

func (s *ProcessServiceTestSuite) TestRegisterProcess_RepositoryError() {
	ctx := context.Background()

	testFile, err := os.Open(testFileAddr)
	s.Require().NoError(err)

	expectedRegisteredProcess := s.getTestProcess()
	expectedRegisteredProcess.Status = entity.RegisterProcessStatusCreating
	customMatcher := newregisteredProcessMatcher(expectedRegisteredProcess)

	s.processRepo.EXPECT().GetByID(ctx, productID, expectedRegisteredProcess.ID).Return(nil, usecase.ErrRegisteredProcessNotFound)
	s.processRepo.EXPECT().Create(productID, customMatcher).Return(nil, fmt.Errorf("mocked error"))

	_, _, err = s.processInteractor.RegisterProcess(
		ctx, user, productID, version, processName, processType, testFile,
	)
	s.Require().Error(err)
}

func (s *ProcessServiceTestSuite) TestListByProduct_WithTypeFilter() {
	ctx := context.Background()

	typeFilter := "trigger"
	expectedRegisteredProcess := []*entity.RegisteredProcess{
		{
			ID:         "test-id",
			Name:       processName,
			Version:    version,
			Type:       processType,
			Image:      "image",
			UploadDate: time.Now(),
			Owner:      userID,
		},
	}

	s.processRepo.EXPECT().ListByProductAndType(ctx, productID, typeFilter).Return(expectedRegisteredProcess, nil)

	returnedRegisteredProcess, err := s.processInteractor.ListByProductAndType(ctx, user, productID, typeFilter)
	s.Require().NoError(err)

	s.Equal(expectedRegisteredProcess, returnedRegisteredProcess)
}

func (s *ProcessServiceTestSuite) TestListByProduct_NoTypeFilter() {
	ctx := context.Background()

	typeFilter := ""
	expectedRegisteredProcess := []*entity.RegisteredProcess{
		{
			ID:         "test-id",
			Name:       processName,
			Version:    version,
			Type:       processType,
			Image:      "image",
			UploadDate: time.Now(),
			Owner:      userID,
		},
	}

	s.processRepo.EXPECT().ListByProductAndType(ctx, productID, typeFilter).Return(expectedRegisteredProcess, nil)

	returnedRegisteredProcess, err := s.processInteractor.ListByProductAndType(ctx, user, productID, "")
	s.Require().NoError(err)

	s.Equal(expectedRegisteredProcess, returnedRegisteredProcess)
}

func (s *ProcessServiceTestSuite) TestListByProduct_InvalidTypeFilterFilter() {
	ctx := context.Background()

	typeFilter := "Kazuma Kiryu"

	_, err := s.processInteractor.ListByProductAndType(ctx, user, productID, typeFilter)
	s.Require().Error(err)
}
