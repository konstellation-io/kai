//go:build unit

package process_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/process"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
)

func (s *ProcessHandlerTestSuite) TestRegisterProcess() {
	ctx := context.Background()
	wg := sync.WaitGroup{}

	wg.Add(1)

	testFile, err := os.Open(_testFilePath)
	s.Require().NoError(err)
	expectedBytes, err := os.ReadFile(_testFilePath)
	s.Require().NoError(err)

	expectedRegisteredProcess := s.getTestProcess(_productID, entity.RegisterProcessStatusCreating)
	expectedCreatedProcess := s.getTestProcess(_productID, entity.RegisterProcessStatusCreated)

	s.accessControl.EXPECT().CheckProductGrants(user, _productID, auth.ActRegisterProcess).Return(nil)
	s.processRepo.EXPECT().GetByID(ctx, _productID, expectedRegisteredProcess.ID).Return(nil, process.ErrRegisteredProcessNotFound)
	s.processRepo.EXPECT().Create(ctx, _productID, expectedRegisteredProcess).Return(nil)
	s.objectStorage.EXPECT().UploadImageSources(ctx, _productID, expectedRegisteredProcess.Image, expectedBytes).Return(nil)

	s.versionService.EXPECT().RegisterProcess(gomock.Any(), _productID, expectedRegisteredProcess.ID, expectedRegisteredProcess.Image).
		Return("", nil)
	s.processRepo.EXPECT().Update(gomock.Any(), _productID, expectedCreatedProcess).Return(nil)
	s.objectStorage.EXPECT().DeleteImageSources(ctx, _productID, expectedRegisteredProcess.Image).
		RunAndReturn(func(_ context.Context, _, _ string) error {
			wg.Done()
			return nil
		}).Once()

	returnedProcess, err := s.processHandler.RegisterProcess(
		ctx, user,
		process.RegisterProcessOpts{
			Product:     _productID,
			Version:     _version,
			Process:     _processName,
			ProcessType: _processType,
			Sources:     testFile,
		},
	)
	s.Require().NoError(err)

	s.Equal(expectedRegisteredProcess, returnedProcess)
	s.Require().NoError(testhelpers.WaitOrTimeout(&wg, 1*time.Second))
}

func (s *ProcessHandlerTestSuite) TestRegisterProcess_OverrideLatest() {
	ctx := context.Background()
	wg := sync.WaitGroup{}
	wg.Add(1)

	latestVersion := "latest"

	testFile, err := os.Open(_testFilePath)
	s.Require().NoError(err)
	expectedBytes, err := os.ReadFile(_testFilePath)
	s.Require().NoError(err)

	expectedRegisteredProcess := testhelpers.NewRegisteredProcessBuilder(_productID).
		WithOwner(user.Email).
		WithVersion(latestVersion).
		WithStatus(entity.RegisterProcessStatusCreating).
		Build()

	customMatcher := newRegisteredProcessMatcher(expectedRegisteredProcess)

	expectedUpdatedProcess := testhelpers.NewRegisteredProcessBuilder(_productID).
		WithOwner(user.Email).
		WithVersion(latestVersion).
		WithStatus(entity.RegisterProcessStatusCreated).
		Build()

	customMatcherUpdate := newRegisteredProcessMatcher(expectedUpdatedProcess)

	existingProcess := &entity.RegisteredProcess{
		ID:      "test",
		Version: latestVersion,
	}

	s.accessControl.EXPECT().CheckProductGrants(user, _productID, auth.ActRegisterProcess).Return(nil)
	s.processRepo.EXPECT().GetByID(ctx, _productID, expectedRegisteredProcess.ID).Return(existingProcess, nil)
	s.processRepo.EXPECT().Update(gomock.Any(), _productID, customMatcher).Return(nil)
	s.objectStorage.EXPECT().UploadImageSources(ctx, _productID, expectedRegisteredProcess.Image, expectedBytes).Return(nil)
	s.versionService.EXPECT().RegisterProcess(gomock.Any(), _productID, expectedRegisteredProcess.ID, expectedRegisteredProcess.Image).
		Return("", nil)
	s.processRepo.EXPECT().Update(gomock.Any(), _productID, customMatcherUpdate).Return(nil)
	s.objectStorage.EXPECT().DeleteImageSources(ctx, _productID, expectedRegisteredProcess.Image).
		RunAndReturn(func(_ context.Context, _, _ string) error {
			wg.Done()
			return nil
		}).Once()

	returnedProcess, err := s.processHandler.RegisterProcess(
		ctx, user,
		process.RegisterProcessOpts{
			Product:     _productID,
			Version:     latestVersion,
			Process:     expectedRegisteredProcess.Name,
			ProcessType: expectedRegisteredProcess.Type,
			Sources:     testFile,
		},
	)
	s.Require().NoError(err)
	s.Equal(expectedRegisteredProcess, returnedProcess)

	s.Require().NoError(testhelpers.WaitOrTimeout(&wg, 1*time.Second))
}

func (s *ProcessHandlerTestSuite) TestRegisterProcess_ProcessAlreadyExistsWithFailedStatus() {
	ctx := context.Background()
	wg := sync.WaitGroup{}
	wg.Add(1)

	testFile, err := os.Open(_testFilePath)
	s.Require().NoError(err)
	expectedBytes, err := os.ReadFile(_testFilePath)
	s.Require().NoError(err)

	alreadyRegisteredProcess := s.getTestProcess(_productID, entity.RegisterProcessStatusFailed)

	expectedCreatingProcess := s.getTestProcess(_productID, entity.RegisterProcessStatusCreating)
	customMatcherCreating := newRegisteredProcessMatcher(expectedCreatingProcess)

	expectedUpdatedProcess := s.getTestProcess(_productID, entity.RegisterProcessStatusCreated)
	customMatcherUpdate := newRegisteredProcessMatcher(expectedUpdatedProcess)

	s.accessControl.EXPECT().CheckProductGrants(user, _productID, auth.ActRegisterProcess).Return(nil)
	s.processRepo.EXPECT().GetByID(ctx, _productID, alreadyRegisteredProcess.ID).Return(alreadyRegisteredProcess, nil)
	s.processRepo.EXPECT().Update(ctx, _productID, customMatcherCreating).Return(nil).Times(1)
	s.objectStorage.EXPECT().UploadImageSources(ctx, _productID, alreadyRegisteredProcess.Image, expectedBytes).Return(nil)
	s.versionService.EXPECT().RegisterProcess(gomock.Any(), _productID, alreadyRegisteredProcess.ID, alreadyRegisteredProcess.Image).
		Return("", nil)
	s.processRepo.EXPECT().Update(gomock.Any(), _productID, customMatcherUpdate).Return(nil).Times(1)
	s.objectStorage.EXPECT().DeleteImageSources(ctx, _productID, alreadyRegisteredProcess.Image).
		RunAndReturn(func(_ context.Context, _, _ string) error {
			wg.Done()
			return nil
		}).Once()

	returnedProcess, err := s.processHandler.RegisterProcess(
		ctx, user,
		process.RegisterProcessOpts{
			Product:     _productID,
			Version:     _version,
			Process:     _processName,
			ProcessType: _processType,
			Sources:     testFile,
		},
	)
	s.Require().NoError(err)

	s.Equal(expectedCreatingProcess, returnedProcess)

	s.Require().NoError(testhelpers.WaitOrTimeout(&wg, 1*time.Second))
}

func (s *ProcessHandlerTestSuite) TestRegisterProcess_MissingProductInRegisterOptions() {
	ctx := context.Background()

	_, err := s.processHandler.RegisterProcess(
		ctx, user,
		process.RegisterProcessOpts{
			Version:     _version,
			Process:     _processName,
			ProcessType: _processType,
			Sources:     nil,
		},
	)
	s.ErrorIs(err, process.ErrMissingProductInParams)
}

func (s *ProcessHandlerTestSuite) TestRegisterProcess_IsPublicAndHasProduct() {
	ctx := context.Background()

	_, err := s.processHandler.RegisterProcess(
		ctx, user,
		process.RegisterProcessOpts{
			Product:  _productID,
			Version:  _version,
			Process:  _processName,
			IsPublic: true,
			Sources:  nil,
		},
	)
	s.ErrorIs(err, process.ErrIsPublicAndHasProduct)
}

func (s *ProcessHandlerTestSuite) TestRegisterProcess_MissingVersionInRegisterOptions() {
	ctx := context.Background()

	_, err := s.processHandler.RegisterProcess(
		ctx, user,
		process.RegisterProcessOpts{
			Product:     _productID,
			Process:     _processName,
			ProcessType: _processType,
			Sources:     nil,
		},
	)
	s.ErrorIs(err, process.ErrMissingVersionInParams)
}

func (s *ProcessHandlerTestSuite) TestRegisterProcess_MissingProcessInRegisterOptions() {
	ctx := context.Background()

	_, err := s.processHandler.RegisterProcess(
		ctx, user,
		process.RegisterProcessOpts{
			Product:     _productID,
			Version:     _version,
			ProcessType: _processType,
			Sources:     nil,
		},
	)
	s.ErrorIs(err, process.ErrMissingProcessInParams)
}

func (s *ProcessHandlerTestSuite) TestRegisterProcess_InvalidProcessTypeInRegisterOptions() {
	ctx := context.Background()

	_, err := s.processHandler.RegisterProcess(
		ctx, user,
		process.RegisterProcessOpts{
			Product:     _productID,
			Version:     _version,
			Process:     _processName,
			ProcessType: "invalid",
			Sources:     nil,
		},
	)
	s.ErrorIs(err, entity.ErrInvalidProcessType)
}

func (s *ProcessHandlerTestSuite) TestRegisterProcess_MissingSourcesInRegisterOptions() {
	ctx := context.Background()

	_, err := s.processHandler.RegisterProcess(
		ctx, user,
		process.RegisterProcessOpts{
			Product:     _productID,
			Version:     _version,
			Process:     _processName,
			ProcessType: _processType,
			Sources:     nil,
		},
	)
	s.ErrorIs(err, process.ErrMissingSourcesInParams)
}

func (s *ProcessHandlerTestSuite) TestRegisterProcess_NoProductGrants() {
	ctx := context.Background()

	expectedErr := errors.New("auth error")

	testFile, err := os.Open(_testFilePath)
	s.Require().NoError(err)

	s.accessControl.EXPECT().CheckProductGrants(user, _productID, auth.ActRegisterProcess).Return(expectedErr)

	_, err = s.processHandler.RegisterProcess(
		ctx, user,
		process.RegisterProcessOpts{
			Product:     _productID,
			Version:     _version,
			Process:     _processName,
			ProcessType: _processType,
			Sources:     testFile,
		},
	)
	s.Require().ErrorIs(err, expectedErr)
}

func (s *ProcessHandlerTestSuite) TestRegisterProcess_GetByIDFails() {
	ctx := context.Background()

	testFile, err := os.Open(_testFilePath)
	s.Require().NoError(err)

	alreadyRegisteredProcess := s.getTestProcess(_productID, entity.RegisterProcessStatusCreating)

	s.accessControl.EXPECT().CheckProductGrants(user, _productID, auth.ActRegisterProcess).Return(nil)
	s.processRepo.EXPECT().GetByID(ctx, _productID, alreadyRegisteredProcess.ID).Return(
		nil, fmt.Errorf("all your base are belong to us"),
	)

	_, err = s.processHandler.RegisterProcess(
		ctx, user,
		process.RegisterProcessOpts{
			Product:     _productID,
			Version:     _version,
			Process:     _processName,
			ProcessType: _processType,
			Sources:     testFile,
		},
	)
	s.Require().Error(err)
}

func (s *ProcessHandlerTestSuite) TestRegisterProcess_ProcessAlreadyExistsAndNotFailed() {
	ctx := context.Background()

	testFile, err := os.Open(_testFilePath)
	s.Require().NoError(err)

	alreadyRegisteredProcess := s.getTestProcess(_productID, entity.RegisterProcessStatusCreating)

	s.accessControl.EXPECT().CheckProductGrants(user, _productID, auth.ActRegisterProcess).Return(nil)
	s.processRepo.EXPECT().GetByID(ctx, _productID, alreadyRegisteredProcess.ID).Return(
		alreadyRegisteredProcess, nil,
	)

	_, err = s.processHandler.RegisterProcess(
		ctx, user,
		process.RegisterProcessOpts{
			Product:     _productID,
			Version:     _version,
			Process:     _processName,
			ProcessType: _processType,
			Sources:     testFile,
		},
	)
	s.Require().Error(err)

	s.ErrorIs(err, process.ErrProcessAlreadyRegistered)
}

func (s *ProcessHandlerTestSuite) TestRegisterProcess_ProcessAlreadyExistsWithFailedStatus_UpdateError() {
	ctx := context.Background()

	testFile, err := os.Open(_testFilePath)
	s.Require().NoError(err)

	alreadyRegisteredProcess := s.getTestProcess(_productID, entity.RegisterProcessStatusFailed)

	expectedCreatingProcess := s.getTestProcess(_productID, entity.RegisterProcessStatusCreating)
	customMatcherCreating := newRegisteredProcessMatcher(expectedCreatingProcess)

	s.accessControl.EXPECT().CheckProductGrants(user, _productID, auth.ActRegisterProcess).Return(nil)
	s.processRepo.EXPECT().GetByID(ctx, _productID, alreadyRegisteredProcess.ID).Return(alreadyRegisteredProcess, nil)
	s.processRepo.EXPECT().Update(ctx, _productID, customMatcherCreating).Return(fmt.Errorf("doctor maligno"))

	_, err = s.processHandler.RegisterProcess(
		ctx, user,
		process.RegisterProcessOpts{
			Product:     _productID,
			Version:     _version,
			Process:     _processName,
			ProcessType: _processType,
			Sources:     testFile,
		},
	)
	s.Require().Error(err)
}

func (s *ProcessHandlerTestSuite) TestRegisterProcess_NoFileError() {
	ctx := context.Background()

	_, err := s.processHandler.RegisterProcess(
		ctx, user,
		process.RegisterProcessOpts{
			Product:     _productID,
			Version:     _version,
			Process:     _processName,
			ProcessType: _processType,
			Sources:     nil,
		},
	)
	s.Require().ErrorIs(err, process.ErrMissingSourcesInParams)
}

func (s *ProcessHandlerTestSuite) TestRegisterProcess_K8sServiceError() {
	ctx := context.Background()
	wg := sync.WaitGroup{}
	wg.Add(1)

	testFile, err := os.Open(_testFilePath)
	s.Require().NoError(err)
	expectedBytes, err := os.ReadFile(_testFilePath)
	s.Require().NoError(err)

	expectedError := errors.New("registering process: mocked error")
	processToRegister := s.getTestProcess(_productID, entity.RegisterProcessStatusCreating)

	customMatcher := newRegisteredProcessMatcher(processToRegister)

	expectedFailedProcess := s.getTestProcess(_productID, entity.RegisterProcessStatusFailed)
	expectedFailedProcess.Logs = expectedError.Error()
	customMatcherUpdate := newRegisteredProcessMatcher(expectedFailedProcess)

	s.accessControl.EXPECT().CheckProductGrants(user, _productID, auth.ActRegisterProcess).Return(nil)
	s.processRepo.EXPECT().GetByID(ctx, _productID, processToRegister.ID).Return(nil, process.ErrRegisteredProcessNotFound)
	s.processRepo.EXPECT().Create(ctx, _productID, customMatcher).Return(nil)
	s.objectStorage.EXPECT().UploadImageSources(ctx, _productID, processToRegister.Image, expectedBytes).Return(nil)
	s.versionService.EXPECT().
		RegisterProcess(gomock.Any(), _productID, processToRegister.ID, processToRegister.Image).
		Return("", errors.New("mocked error"))
	s.processRepo.EXPECT().Update(gomock.Any(), _productID, customMatcherUpdate).Return(nil)
	s.objectStorage.EXPECT().DeleteImageSources(ctx, _productID, processToRegister.Image).
		RunAndReturn(func(_ context.Context, _, _ string) error {
			wg.Done()
			return nil
		}).Once()

	returnedRef, err := s.processHandler.RegisterProcess(
		ctx, user,
		process.RegisterProcessOpts{
			Product:     _productID,
			Version:     _version,
			Process:     _processName,
			ProcessType: _processType,
			Sources:     testFile,
		},
	)
	s.Require().NoError(err)

	s.Require().NoError(testhelpers.WaitOrTimeout(&wg, 1*time.Second))
	s.Equal(expectedFailedProcess, returnedRef)
}

func (s *ProcessHandlerTestSuite) TestRegisterProcess_RepositoryError() {
	ctx := context.Background()

	testFile, err := os.Open(_testFilePath)
	s.Require().NoError(err)

	expectedError := errors.New("process repo error")
	expectedRegisteredProcess := s.getTestProcess(_productID, entity.RegisterProcessStatusCreating)
	customMatcher := newRegisteredProcessMatcher(expectedRegisteredProcess)

	s.accessControl.EXPECT().CheckProductGrants(user, _productID, auth.ActRegisterProcess).Return(nil)
	s.processRepo.EXPECT().GetByID(ctx, _productID, expectedRegisteredProcess.ID).Return(nil, process.ErrRegisteredProcessNotFound)
	s.processRepo.EXPECT().Create(ctx, _productID, customMatcher).Return(expectedError)

	_, err = s.processHandler.RegisterProcess(
		ctx, user,
		process.RegisterProcessOpts{
			Product:     _productID,
			Version:     _version,
			Process:     _processName,
			ProcessType: _processType,
			Sources:     testFile,
		},
	)
	s.Require().ErrorIs(err, expectedError)
}

func (s *ProcessHandlerTestSuite) TestRegisterProcess_UpdateError() {
	ctx := context.Background()
	wg := sync.WaitGroup{}
	wg.Add(1)

	testFile, err := os.Open(_testFilePath)
	s.Require().NoError(err)
	expectedBytes, err := os.ReadFile(_testFilePath)
	s.Require().NoError(err)

	expectedRegisteredProcess := s.getTestProcess(_productID, entity.RegisterProcessStatusCreating)
	customMatcher := newRegisteredProcessMatcher(expectedRegisteredProcess)
	expectedUpdatedProcess := s.getTestProcess(_productID, entity.RegisterProcessStatusCreated)

	s.accessControl.EXPECT().CheckProductGrants(user, _productID, auth.ActRegisterProcess).Return(nil)
	s.processRepo.EXPECT().GetByID(ctx, _productID, expectedRegisteredProcess.ID).Return(nil, process.ErrRegisteredProcessNotFound)
	s.processRepo.EXPECT().Create(ctx, _productID, customMatcher).Return(nil)
	s.objectStorage.EXPECT().UploadImageSources(ctx, _productID, expectedRegisteredProcess.Image, expectedBytes).Return(nil)
	s.processRepo.EXPECT().Update(gomock.Any(), _productID, expectedUpdatedProcess).Return(errors.New("update error"))
	s.versionService.EXPECT().RegisterProcess(gomock.Any(), _productID, expectedRegisteredProcess.ID, expectedRegisteredProcess.Image).
		Return("", nil)
	s.objectStorage.EXPECT().DeleteImageSources(ctx, _productID, expectedRegisteredProcess.Image).
		RunAndReturn(func(_ context.Context, _, _ string) error {
			wg.Done()
			return nil
		}).Once()

	returnedProcess, err := s.processHandler.RegisterProcess(
		ctx, user,
		process.RegisterProcessOpts{
			Product:     _productID,
			Version:     _version,
			Process:     _processName,
			ProcessType: _processType,
			Sources:     testFile,
		},
	)
	s.Require().NoError(err)
	s.Equal(expectedRegisteredProcess, returnedProcess)

	s.Require().NoError(testhelpers.WaitOrTimeout(&wg, 1*time.Second))
}

func (s *ProcessHandlerTestSuite) TestRegisterProcess_Public() {
	var (
		ctx                       = context.Background()
		wg                        = sync.WaitGroup{}
		processWithCreatingStatus = s.getTestProcess(_publicRegistry, entity.RegisterProcessStatusCreating, true)
		processWithCreatedStatus  = s.getTestProcess(_publicRegistry, entity.RegisterProcessStatusCreated, true)
	)

	wg.Add(1)

	testFile, err := os.Open(_testFilePath)
	s.Require().NoError(err)
	expectedBytes, err := os.ReadFile(_testFilePath)
	s.Require().NoError(err)

	s.accessControl.EXPECT().CheckRoleGrants(user, auth.ActRegisterPublicProcess).Return(nil)
	s.processRepo.EXPECT().GetByID(ctx, _publicRegistry, processWithCreatingStatus.ID).Return(nil, process.ErrRegisteredProcessNotFound)
	s.processRepo.EXPECT().Create(ctx, _publicRegistry, processWithCreatingStatus).Return(nil)
	s.objectStorage.EXPECT().UploadImageSources(ctx, _publicRegistry, processWithCreatingStatus.Image, expectedBytes).Return(nil)

	s.versionService.EXPECT().RegisterProcess(gomock.Any(), _publicRegistry, processWithCreatingStatus.ID, processWithCreatingStatus.Image).
		Return("", nil)
	s.processRepo.EXPECT().Update(gomock.Any(), _publicRegistry, processWithCreatedStatus).Return(nil)
	s.objectStorage.EXPECT().DeleteImageSources(ctx, _publicRegistry, processWithCreatedStatus.Image).
		RunAndReturn(func(_ context.Context, _, _ string) error {
			wg.Done()
			return nil
		}).Once()

	returnedProcess, err := s.processHandler.RegisterProcess(
		ctx, user,
		process.RegisterProcessOpts{
			Version:     _version,
			Process:     _processName,
			ProcessType: _processType,
			Sources:     testFile,
			IsPublic:    true,
		},
	)
	s.Require().NoError(err)

	s.Equal(processWithCreatingStatus, returnedProcess)
	s.Require().NoError(testhelpers.WaitOrTimeout(&wg, 1*time.Second))
}

func (s *ProcessHandlerTestSuite) TestRegisterProcess_ErrorImageIDTooLong() {
	ctx := context.Background()
	longProcessName := faker.UUIDHyphenated(options.WithRandomStringLength(64))

	s.accessControl.EXPECT().CheckProductGrants(user, _productID, auth.ActRegisterProcess).Return(nil)

	_, err := s.processHandler.RegisterProcess(
		ctx, user,
		process.RegisterProcessOpts{
			Product:     _productID,
			Version:     _version,
			Process:     longProcessName,
			ProcessType: _processType,
			Sources:     bytes.NewReader([]byte("sources")),
		},
	)
	s.Require().ErrorIs(err, process.ErrProcessNameTooLong)
}
