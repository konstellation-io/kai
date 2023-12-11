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
	"github.com/konstellation-io/kai/engine/admin-api/domain/repository"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase"
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
	processService *usecase.ProcessService

	registryHost string
}

const (
	userID       = "userID"
	userEmail    = "test@email.com"
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
		Email: userEmail,
	}
)

func (s *ProcessServiceTestSuite) getTestProcess(status string) *entity.RegisteredProcess {
	return testhelpers.NewRegisteredProcessBuilder(productID).
		WithName(processName).
		WithVersion(version).
		WithType(processType).
		WithOwner(user.Email).
		WithStatus(status).
		Build()
}

func TestProcessTestSuite(t *testing.T) {
	suite.Run(t, new(ProcessServiceTestSuite))
}

func (s *ProcessServiceTestSuite) SetupSuite() {
	logger := zapr.NewLogger(zap.NewNop())
	s.ctrl = gomock.NewController(s.T())
	s.processRepo = mocks.NewMockProcessRepository(s.ctrl)
	s.versionService = mocks.NewMockVersionService(s.ctrl)
	s.objectStorage = mocks.NewMockObjectStorage(s.T())
	s.processService = usecase.NewProcessService(logger, s.versionService, s.processRepo, s.objectStorage)

	s.registryHost = "test.registry"

	viper.Set(config.RegistryHostKey, s.registryHost)

	monkey.Patch(time.Now, func() time.Time {
		return time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	})
}

func (s *ProcessServiceTestSuite) TearDownSuite() {
	monkey.UnpatchAll()
}

func (s *ProcessServiceTestSuite) TestRegisterProcess() {
	ctx := context.Background()

	testFile, err := os.Open(testFileAddr)
	s.Require().NoError(err)
	expectedBytes, err := os.ReadFile(testFileAddr)
	s.Require().NoError(err)

	expectedRegisteredProcess := s.getTestProcess(entity.RegisterProcessStatusCreating)
	customMatcher := newRegisteredProcessMatcher(expectedRegisteredProcess)

	expectedUpdatedProcess := s.getTestProcess(entity.RegisterProcessStatusCreated)
	customMatcherUpdate := newRegisteredProcessMatcher(expectedUpdatedProcess)

	s.processRepo.EXPECT().GetByID(ctx, productID, expectedRegisteredProcess.ID).Return(nil, usecase.ErrRegisteredProcessNotFound)
	s.processRepo.EXPECT().Create(productID, customMatcher).Return(nil, nil)
	s.objectStorage.EXPECT().UploadImageSources(ctx, productID, expectedRegisteredProcess.Image, expectedBytes).Return(nil)

	s.versionService.EXPECT().RegisterProcess(gomock.Any(), productID, expectedRegisteredProcess.ID, expectedRegisteredProcess.Image).Return("", nil)
	s.processRepo.EXPECT().Update(gomock.Any(), productID, customMatcherUpdate).Return(nil)
	s.objectStorage.EXPECT().DeleteImageSources(ctx, productID, expectedRegisteredProcess.Image).Return(nil)

	returnedProcess, notifyCh, err := s.processService.RegisterProcess(
		ctx, user, productID, version, processName, processType, testFile,
	)
	s.Require().NoError(err)

	s.Equal(expectedRegisteredProcess, returnedProcess)

	registeredProcess := <-notifyCh
	s.Equal(entity.RegisterProcessStatusCreated, registeredProcess.Status)
}

func (s *ProcessServiceTestSuite) TestRegisterProcess_OverrideLatest() {
	ctx := context.Background()

	latestVersion := "latest"

	testFile, err := os.Open(testFileAddr)
	s.Require().NoError(err)
	expectedBytes, err := os.ReadFile(testFileAddr)
	s.Require().NoError(err)

	expectedRegisteredProcess := testhelpers.NewRegisteredProcessBuilder(productID).
		WithOwner(user.Email).
		WithVersion(latestVersion).
		WithStatus(entity.RegisterProcessStatusCreating).
		Build()

	customMatcher := newRegisteredProcessMatcher(expectedRegisteredProcess)

	expectedUpdatedProcess := testhelpers.NewRegisteredProcessBuilder(productID).
		WithOwner(user.Email).
		WithVersion(latestVersion).
		WithStatus(entity.RegisterProcessStatusCreated).
		Build()

	customMatcherUpdate := newRegisteredProcessMatcher(expectedUpdatedProcess)

	existingProcess := &entity.RegisteredProcess{
		ID:      "test",
		Version: latestVersion,
	}

	s.processRepo.EXPECT().GetByID(ctx, productID, expectedRegisteredProcess.ID).Return(existingProcess, nil)
	s.processRepo.EXPECT().Update(gomock.Any(), productID, customMatcher).Return(nil)
	s.objectStorage.EXPECT().UploadImageSources(ctx, productID, expectedRegisteredProcess.Image, expectedBytes).Return(nil)
	s.versionService.EXPECT().RegisterProcess(gomock.Any(), productID, expectedRegisteredProcess.ID, expectedRegisteredProcess.Image).Return("", nil)
	s.processRepo.EXPECT().Update(gomock.Any(), productID, customMatcherUpdate).Return(nil)
	s.objectStorage.EXPECT().DeleteImageSources(ctx, productID, expectedRegisteredProcess.Image).Return(nil)

	returnedProcess, notifyCh, err := s.processService.RegisterProcess(
		ctx, user, productID, latestVersion, expectedRegisteredProcess.Name,
		expectedRegisteredProcess.Type, testFile,
	)
	s.Require().NoError(err)

	s.Equal(expectedRegisteredProcess, returnedProcess)

	registeredProcess := <-notifyCh
	s.Equal(entity.RegisterProcessStatusCreated, registeredProcess.Status)
}

func (s *ProcessServiceTestSuite) TestRegisterProcess_ProcessAlreadyExistsWithFailedStatus() {
	ctx := context.Background()

	testFile, err := os.Open(testFileAddr)
	s.Require().NoError(err)
	expectedBytes, err := os.ReadFile(testFileAddr)
	s.Require().NoError(err)

	alreadyRegisteredProcess := s.getTestProcess(entity.RegisterProcessStatusFailed)

	expectedCreatingProcess := s.getTestProcess(entity.RegisterProcessStatusCreating)
	customMatcherCreating := newRegisteredProcessMatcher(expectedCreatingProcess)

	expectedUpdatedProcess := s.getTestProcess(entity.RegisterProcessStatusCreated)
	customMatcherUpdate := newRegisteredProcessMatcher(expectedUpdatedProcess)

	s.processRepo.EXPECT().GetByID(ctx, productID, alreadyRegisteredProcess.ID).Return(alreadyRegisteredProcess, nil)
	s.processRepo.EXPECT().Update(ctx, productID, customMatcherCreating).Return(nil).Times(1)
	s.objectStorage.EXPECT().UploadImageSources(ctx, productID, alreadyRegisteredProcess.Image, expectedBytes).Return(nil)
	s.versionService.EXPECT().RegisterProcess(gomock.Any(), productID, alreadyRegisteredProcess.ID, alreadyRegisteredProcess.Image).Return("", nil)
	s.processRepo.EXPECT().Update(gomock.Any(), productID, customMatcherUpdate).Return(nil).Times(1)
	s.objectStorage.EXPECT().DeleteImageSources(ctx, productID, alreadyRegisteredProcess.Image).Return(nil)

	returnedProcess, notifyCh, err := s.processService.RegisterProcess(
		ctx, user, productID, version, processName, processType, testFile,
	)
	s.Require().NoError(err)

	s.Equal(expectedCreatingProcess, returnedProcess)

	registeredProcess := <-notifyCh
	s.Equal(entity.RegisterProcessStatusCreated, registeredProcess.Status)
}

func (s *ProcessServiceTestSuite) TestRegisterProcess_GetByIDFails() {
	ctx := context.Background()

	testFile, err := os.Open(testFileAddr)
	s.Require().NoError(err)

	alreadyRegisteredProcess := s.getTestProcess(entity.RegisterProcessStatusCreating)

	s.processRepo.EXPECT().GetByID(ctx, productID, alreadyRegisteredProcess.ID).Return(
		nil, fmt.Errorf("all your base are belong to us"),
	)

	_, _, err = s.processService.RegisterProcess(
		ctx, user, productID, version, processName, processType, testFile,
	)
	s.Require().Error(err)
}

func (s *ProcessServiceTestSuite) TestRegisterProcess_ProcessAlreadyExistsAndNotFailed() {
	ctx := context.Background()

	testFile, err := os.Open(testFileAddr)
	s.Require().NoError(err)

	alreadyRegisteredProcess := s.getTestProcess(entity.RegisterProcessStatusCreating)

	s.processRepo.EXPECT().GetByID(ctx, productID, alreadyRegisteredProcess.ID).Return(
		alreadyRegisteredProcess, nil,
	)

	_, _, err = s.processService.RegisterProcess(
		ctx, user, productID, version, processName, processType, testFile,
	)
	s.Require().Error(err)

	s.ErrorIs(err, usecase.ErrProcessAlreadyRegistered)
}

func (s *ProcessServiceTestSuite) TestRegisterProcess_ProcessAlreadyExistsWithFailedStatus_UpdateError() {
	ctx := context.Background()

	testFile, err := os.Open(testFileAddr)
	s.Require().NoError(err)

	alreadyRegisteredProcess := s.getTestProcess(entity.RegisterProcessStatusFailed)

	expectedCreatingProcess := s.getTestProcess(entity.RegisterProcessStatusCreating)
	customMatcherCreating := newRegisteredProcessMatcher(expectedCreatingProcess)

	s.processRepo.EXPECT().GetByID(ctx, productID, alreadyRegisteredProcess.ID).Return(alreadyRegisteredProcess, nil)
	s.processRepo.EXPECT().Update(ctx, productID, customMatcherCreating).Return(fmt.Errorf("doctor maligno"))

	_, _, err = s.processService.RegisterProcess(
		ctx, user, productID, version, processName, processType, testFile,
	)
	s.Require().Error(err)
}

func (s *ProcessServiceTestSuite) TestRegisterProcess_NoFileError() {
	ctx := context.Background()

	testFile, err := os.Open("no-file")
	s.Require().Error(err)

	expectedRegisteredProcess := s.getTestProcess(entity.RegisterProcessStatusCreating)
	customMatcher := newRegisteredProcessMatcher(expectedRegisteredProcess)

	expectedUpdatedProcess := s.getTestProcess(entity.RegisterProcessStatusFailed)
	expectedUpdatedProcess.Logs = "copying temp file for version: invalid argument"
	customMatcherUpdate := newRegisteredProcessMatcher(expectedUpdatedProcess)

	s.processRepo.EXPECT().GetByID(ctx, productID, expectedRegisteredProcess.ID).Return(nil, usecase.ErrRegisteredProcessNotFound)
	s.processRepo.EXPECT().Create(productID, customMatcher).Return(nil, nil)
	s.processRepo.EXPECT().Update(gomock.Any(), productID, customMatcherUpdate).Return(nil)

	returnedRef, notifyCh, err := s.processService.RegisterProcess(
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

	expectedRegisteredProcess := s.getTestProcess(entity.RegisterProcessStatusCreating)
	customMatcher := newRegisteredProcessMatcher(expectedRegisteredProcess)

	expectedUpdatedProcess := s.getTestProcess(entity.RegisterProcessStatusFailed)
	expectedUpdatedProcess.Logs = "registering process: mocked error"
	customMatcherUpdate := newRegisteredProcessMatcher(expectedUpdatedProcess)

	s.processRepo.EXPECT().GetByID(ctx, productID, expectedRegisteredProcess.ID).Return(nil, usecase.ErrRegisteredProcessNotFound)
	s.processRepo.EXPECT().Create(productID, customMatcher).Return(nil, nil)
	s.versionService.EXPECT().
		RegisterProcess(gomock.Any(), productID, expectedRegisteredProcess.ID, expectedRegisteredProcess.Image).
		Return("", fmt.Errorf("mocked error"))
	s.processRepo.EXPECT().Update(gomock.Any(), productID, customMatcherUpdate).Return(nil)

	returnedRef, notifyCh, err := s.processService.RegisterProcess(
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

	expectedRegisteredProcess := s.getTestProcess(entity.RegisterProcessStatusCreating)
	customMatcher := newRegisteredProcessMatcher(expectedRegisteredProcess)

	s.processRepo.EXPECT().GetByID(ctx, productID, expectedRegisteredProcess.ID).Return(nil, usecase.ErrRegisteredProcessNotFound)
	s.processRepo.EXPECT().Create(productID, customMatcher).Return(nil, fmt.Errorf("mocked error"))

	_, _, err = s.processService.RegisterProcess(
		ctx, user, productID, version, processName, processType, testFile,
	)
	s.Require().Error(err)
}

func (s *ProcessServiceTestSuite) TestRegisterProcess_UpdateError() {
	ctx := context.Background()

	testFile, err := os.Open(testFileAddr)
	s.Require().NoError(err)

	expectedRegisteredProcess := s.getTestProcess(entity.RegisterProcessStatusCreating)
	customMatcher := newRegisteredProcessMatcher(expectedRegisteredProcess)

	expectedUpdatedProcess := s.getTestProcess(entity.RegisterProcessStatusCreated)
	customMatcherUpdate := newRegisteredProcessMatcher(expectedUpdatedProcess)

	s.processRepo.EXPECT().GetByID(ctx, productID, expectedRegisteredProcess.ID).Return(nil, usecase.ErrRegisteredProcessNotFound)
	s.processRepo.EXPECT().Create(productID, customMatcher).Return(nil, nil)
	s.versionService.EXPECT().RegisterProcess(gomock.Any(), productID, expectedRegisteredProcess.ID, expectedRegisteredProcess.Image).Return("", nil)
	s.processRepo.EXPECT().Update(gomock.Any(), productID, customMatcherUpdate).Return(fmt.Errorf("listen to Death Grips"))

	returnedProcess, notifyCh, err := s.processService.RegisterProcess(
		ctx, user, productID, version, processName, processType, testFile,
	)
	s.Require().NoError(err)

	s.Equal(expectedRegisteredProcess, returnedProcess)

	registeredProcess := <-notifyCh
	s.Equal(entity.RegisterProcessStatusFailed, registeredProcess.Status)
}

func (s *ProcessServiceTestSuite) TestListByProduct_WithTypeFilter() {
	var (
		ctx               = context.Background()
		filter            = repository.SearchFilter{ProcessType: entity.ProcessTypeTrigger}
		productProcesses  = []*entity.RegisteredProcess{testhelpers.NewRegisteredProcessBuilder(productID).Build()}
		kaiProcesses      = []*entity.RegisteredProcess{testhelpers.NewRegisteredProcessBuilder("kai").Build()}
		expectedProcesses = append(productProcesses, kaiProcesses...)
	)

	s.processRepo.EXPECT().SearchByProduct(ctx, productID, filter).Return(productProcesses, nil)
	s.processRepo.EXPECT().GlobalSearch(ctx, filter).Return(kaiProcesses, nil)

	returnedRegisteredProcess, err := s.processService.ListByProductAndType(ctx, user, productID, filter.ProcessType.String())
	s.Require().NoError(err)

	s.Equal(expectedProcesses, returnedRegisteredProcess)
}

func (s *ProcessServiceTestSuite) TestListByProduct_NoTypeFilter() {
	ctx := context.Background()

	filter := repository.SearchFilter{}
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

	s.processRepo.EXPECT().SearchByProduct(ctx, productID, filter).Return(expectedRegisteredProcess, nil)
	s.processRepo.EXPECT().GlobalSearch(ctx, filter).Return(nil, nil)

	returnedRegisteredProcess, err := s.processService.ListByProductAndType(ctx, user, productID, "")
	s.Require().NoError(err)

	s.Equal(expectedRegisteredProcess, returnedRegisteredProcess)
}

func (s *ProcessServiceTestSuite) TestListByProduct_InvalidTypeFilterFilter() {
	ctx := context.Background()

	typeFilter := "invalid type"

	_, err := s.processService.ListByProductAndType(ctx, user, productID, typeFilter)
	s.Require().Error(err)
}
