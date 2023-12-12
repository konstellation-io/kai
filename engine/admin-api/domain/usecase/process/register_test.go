//go:build unit

package process_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/process"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
)

func (s *ProcessServiceTestSuite) TestRegisterProcess() {
	ctx := context.Background()
	wg := sync.WaitGroup{}

	wg.Add(1)

	testFile, err := os.Open(testFileAddr)
	s.Require().NoError(err)
	expectedBytes, err := os.ReadFile(testFileAddr)
	s.Require().NoError(err)

	expectedRegisteredProcess := s.getTestProcess(productID, entity.RegisterProcessStatusCreating)
	expectedCreatedProcess := s.getTestProcess(productID, entity.RegisterProcessStatusCreated)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActRegisterProcess).Return(nil)
	s.processRepo.EXPECT().GetByID(ctx, productID, expectedRegisteredProcess.ID).Return(nil, process.ErrRegisteredProcessNotFound)
	s.processRepo.EXPECT().Create(ctx, productID, expectedRegisteredProcess).Return(nil)
	s.objectStorage.EXPECT().UploadImageSources(ctx, productID, expectedRegisteredProcess.Image, expectedBytes).Return(nil)

	s.versionService.EXPECT().RegisterProcess(gomock.Any(), productID, expectedRegisteredProcess.ID, expectedRegisteredProcess.Image).
		Return("", nil)
	s.processRepo.EXPECT().Update(gomock.Any(), productID, expectedCreatedProcess).Return(nil)
	s.objectStorage.EXPECT().DeleteImageSources(ctx, productID, expectedRegisteredProcess.Image).
		RunAndReturn(func(_ context.Context, _, _ string) error {
			wg.Done()
			return nil
		}).Once()

	returnedProcess, err := s.processService.RegisterProcess(
		ctx, user,
		process.RegisterProcessOpts{
			Product:     productID,
			Version:     version,
			Process:     processName,
			ProcessType: processType,
			Sources:     testFile,
		},
	)
	s.Require().NoError(err)

	s.Equal(expectedRegisteredProcess, returnedProcess)
	s.Require().NoError(testhelpers.WaitOrTimeout(&wg, 1*time.Second))
}

func (s *ProcessServiceTestSuite) TestRegisterProcess_OverrideLatest() {
	ctx := context.Background()
	wg := sync.WaitGroup{}
	wg.Add(1)

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

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActRegisterProcess).Return(nil)
	s.processRepo.EXPECT().GetByID(ctx, productID, expectedRegisteredProcess.ID).Return(existingProcess, nil)
	s.processRepo.EXPECT().Update(gomock.Any(), productID, customMatcher).Return(nil)
	s.objectStorage.EXPECT().UploadImageSources(ctx, productID, expectedRegisteredProcess.Image, expectedBytes).Return(nil)
	s.versionService.EXPECT().RegisterProcess(gomock.Any(), productID, expectedRegisteredProcess.ID, expectedRegisteredProcess.Image).
		Return("", nil)
	s.processRepo.EXPECT().Update(gomock.Any(), productID, customMatcherUpdate).Return(nil)
	s.objectStorage.EXPECT().DeleteImageSources(ctx, productID, expectedRegisteredProcess.Image).
		RunAndReturn(func(_ context.Context, _, _ string) error {
			wg.Done()
			return nil
		}).Once()

	returnedProcess, err := s.processService.RegisterProcess(
		ctx, user,
		process.RegisterProcessOpts{
			Product:     productID,
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

func (s *ProcessServiceTestSuite) TestRegisterProcess_ProcessAlreadyExistsWithFailedStatus() {
	ctx := context.Background()
	wg := sync.WaitGroup{}
	wg.Add(1)

	testFile, err := os.Open(testFileAddr)
	s.Require().NoError(err)
	expectedBytes, err := os.ReadFile(testFileAddr)
	s.Require().NoError(err)

	alreadyRegisteredProcess := s.getTestProcess(productID, entity.RegisterProcessStatusFailed)

	expectedCreatingProcess := s.getTestProcess(productID, entity.RegisterProcessStatusCreating)
	customMatcherCreating := newRegisteredProcessMatcher(expectedCreatingProcess)

	expectedUpdatedProcess := s.getTestProcess(productID, entity.RegisterProcessStatusCreated)
	customMatcherUpdate := newRegisteredProcessMatcher(expectedUpdatedProcess)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActRegisterProcess).Return(nil)
	s.processRepo.EXPECT().GetByID(ctx, productID, alreadyRegisteredProcess.ID).Return(alreadyRegisteredProcess, nil)
	s.processRepo.EXPECT().Update(ctx, productID, customMatcherCreating).Return(nil).Times(1)
	s.objectStorage.EXPECT().UploadImageSources(ctx, productID, alreadyRegisteredProcess.Image, expectedBytes).Return(nil)
	s.versionService.EXPECT().RegisterProcess(gomock.Any(), productID, alreadyRegisteredProcess.ID, alreadyRegisteredProcess.Image).
		Return("", nil)
	s.processRepo.EXPECT().Update(gomock.Any(), productID, customMatcherUpdate).Return(nil).Times(1)
	s.objectStorage.EXPECT().DeleteImageSources(ctx, productID, alreadyRegisteredProcess.Image).
		RunAndReturn(func(_ context.Context, _, _ string) error {
			wg.Done()
			return nil
		}).Once()

	returnedProcess, err := s.processService.RegisterProcess(
		ctx, user,
		process.RegisterProcessOpts{
			Product:     productID,
			Version:     version,
			Process:     processName,
			ProcessType: processType,
			Sources:     testFile,
		},
	)
	s.Require().NoError(err)

	s.Equal(expectedCreatingProcess, returnedProcess)

	s.Require().NoError(testhelpers.WaitOrTimeout(&wg, 1*time.Second))
}

func (s *ProcessServiceTestSuite) TestRegisterProcess_GetByIDFails() {
	ctx := context.Background()

	testFile, err := os.Open(testFileAddr)
	s.Require().NoError(err)

	alreadyRegisteredProcess := s.getTestProcess(productID, entity.RegisterProcessStatusCreating)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActRegisterProcess).Return(nil)
	s.processRepo.EXPECT().GetByID(ctx, productID, alreadyRegisteredProcess.ID).Return(
		nil, fmt.Errorf("all your base are belong to us"),
	)

	_, err = s.processService.RegisterProcess(
		ctx, user,
		process.RegisterProcessOpts{
			Product:     productID,
			Version:     version,
			Process:     processName,
			ProcessType: processType,
			Sources:     testFile,
		},
	)
	s.Require().Error(err)
}

func (s *ProcessServiceTestSuite) TestRegisterProcess_ProcessAlreadyExistsAndNotFailed() {
	ctx := context.Background()

	testFile, err := os.Open(testFileAddr)
	s.Require().NoError(err)

	alreadyRegisteredProcess := s.getTestProcess(productID, entity.RegisterProcessStatusCreating)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActRegisterProcess).Return(nil)
	s.processRepo.EXPECT().GetByID(ctx, productID, alreadyRegisteredProcess.ID).Return(
		alreadyRegisteredProcess, nil,
	)

	_, err = s.processService.RegisterProcess(
		ctx, user,
		process.RegisterProcessOpts{
			Product:     productID,
			Version:     version,
			Process:     processName,
			ProcessType: processType,
			Sources:     testFile,
		},
	)
	s.Require().Error(err)

	s.ErrorIs(err, process.ErrProcessAlreadyRegistered)
}

func (s *ProcessServiceTestSuite) TestRegisterProcess_ProcessAlreadyExistsWithFailedStatus_UpdateError() {
	ctx := context.Background()

	testFile, err := os.Open(testFileAddr)
	s.Require().NoError(err)

	alreadyRegisteredProcess := s.getTestProcess(productID, entity.RegisterProcessStatusFailed)

	expectedCreatingProcess := s.getTestProcess(productID, entity.RegisterProcessStatusCreating)
	customMatcherCreating := newRegisteredProcessMatcher(expectedCreatingProcess)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActRegisterProcess).Return(nil)
	s.processRepo.EXPECT().GetByID(ctx, productID, alreadyRegisteredProcess.ID).Return(alreadyRegisteredProcess, nil)
	s.processRepo.EXPECT().Update(ctx, productID, customMatcherCreating).Return(fmt.Errorf("doctor maligno"))

	_, err = s.processService.RegisterProcess(
		ctx, user,
		process.RegisterProcessOpts{
			Product:     productID,
			Version:     version,
			Process:     processName,
			ProcessType: processType,
			Sources:     testFile,
		},
	)
	s.Require().Error(err)
}

func (s *ProcessServiceTestSuite) TestRegisterProcess_NoFileError() {
	ctx := context.Background()

	testFile, err := os.Open("no-file")
	s.Require().Error(err)

	expectedRegisteredProcess := s.getTestProcess(productID, entity.RegisterProcessStatusCreating)
	customMatcher := newRegisteredProcessMatcher(expectedRegisteredProcess)

	expectedUpdatedProcess := s.getTestProcess(productID, entity.RegisterProcessStatusFailed)
	expectedUpdatedProcess.Logs = "copying temp file for version: invalid argument"
	customMatcherUpdate := newRegisteredProcessMatcher(expectedUpdatedProcess)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActRegisterProcess).Return(nil)
	s.processRepo.EXPECT().GetByID(ctx, productID, expectedRegisteredProcess.ID).Return(nil, process.ErrRegisteredProcessNotFound)
	s.processRepo.EXPECT().Create(ctx, productID, customMatcher).Return(nil)
	s.processRepo.EXPECT().Update(gomock.Any(), productID, customMatcherUpdate).Return(nil)

	returnedRef, err := s.processService.RegisterProcess(
		ctx, user,
		process.RegisterProcessOpts{
			Product:     productID,
			Version:     version,
			Process:     processName,
			ProcessType: processType,
			Sources:     testFile,
		},
	)
	s.Require().NoError(err)

	s.Equal(expectedRegisteredProcess, returnedRef)

	s.Equal(entity.RegisterProcessStatusFailed, expectedUpdatedProcess.Status)
}

func (s *ProcessServiceTestSuite) TestRegisterProcess_K8sServiceError() {
	ctx := context.Background()
	wg := sync.WaitGroup{}
	wg.Add(1)

	testFile, err := os.Open(testFileAddr)
	s.Require().NoError(err)
	expectedBytes, err := os.ReadFile(testFileAddr)
	s.Require().NoError(err)

	expectedError := errors.New("registering process: mocked error")
	processToRegister := s.getTestProcess(productID, entity.RegisterProcessStatusCreating)

	customMatcher := newRegisteredProcessMatcher(processToRegister)

	expectedFailedProcess := s.getTestProcess(productID, entity.RegisterProcessStatusFailed)
	expectedFailedProcess.Logs = expectedError.Error()
	customMatcherUpdate := newRegisteredProcessMatcher(expectedFailedProcess)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActRegisterProcess).Return(nil)
	s.processRepo.EXPECT().GetByID(ctx, productID, processToRegister.ID).Return(nil, process.ErrRegisteredProcessNotFound)
	s.processRepo.EXPECT().Create(ctx, productID, customMatcher).Return(nil)
	s.objectStorage.EXPECT().UploadImageSources(ctx, productID, processToRegister.Image, expectedBytes).Return(nil)
	s.versionService.EXPECT().
		RegisterProcess(gomock.Any(), productID, processToRegister.ID, processToRegister.Image).
		Return("", errors.New("mocked error"))
	s.processRepo.EXPECT().Update(gomock.Any(), productID, customMatcherUpdate).Return(nil)
	s.objectStorage.EXPECT().DeleteImageSources(ctx, productID, processToRegister.Image).
		RunAndReturn(func(_ context.Context, _, _ string) error {
			wg.Done()
			return nil
		}).Once()

	returnedRef, err := s.processService.RegisterProcess(
		ctx, user,
		process.RegisterProcessOpts{
			Product:     productID,
			Version:     version,
			Process:     processName,
			ProcessType: processType,
			Sources:     testFile,
		},
	)
	s.Require().NoError(err)

	s.Require().NoError(testhelpers.WaitOrTimeout(&wg, 1*time.Second))
	s.Equal(expectedFailedProcess, returnedRef)
}

func (s *ProcessServiceTestSuite) TestRegisterProcess_RepositoryError() {
	ctx := context.Background()

	testFile, err := os.Open(testFileAddr)
	s.Require().NoError(err)

	expectedError := errors.New("process repo error")
	expectedRegisteredProcess := s.getTestProcess(productID, entity.RegisterProcessStatusCreating)
	customMatcher := newRegisteredProcessMatcher(expectedRegisteredProcess)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActRegisterProcess).Return(nil)
	s.processRepo.EXPECT().GetByID(ctx, productID, expectedRegisteredProcess.ID).Return(nil, process.ErrRegisteredProcessNotFound)
	s.processRepo.EXPECT().Create(ctx, productID, customMatcher).Return(expectedError)

	_, err = s.processService.RegisterProcess(
		ctx, user,
		process.RegisterProcessOpts{
			Product:     productID,
			Version:     version,
			Process:     processName,
			ProcessType: processType,
			Sources:     testFile,
		},
	)
	s.Require().ErrorIs(err, expectedError)
}

func (s *ProcessServiceTestSuite) TestRegisterProcess_UpdateError() {
	ctx := context.Background()
	wg := sync.WaitGroup{}
	wg.Add(1)

	testFile, err := os.Open(testFileAddr)
	s.Require().NoError(err)
	expectedBytes, err := os.ReadFile(testFileAddr)
	s.Require().NoError(err)

	expectedRegisteredProcess := s.getTestProcess(productID, entity.RegisterProcessStatusCreating)
	customMatcher := newRegisteredProcessMatcher(expectedRegisteredProcess)
	expectedUpdatedProcess := s.getTestProcess(productID, entity.RegisterProcessStatusCreated)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActRegisterProcess).Return(nil)
	s.processRepo.EXPECT().GetByID(ctx, productID, expectedRegisteredProcess.ID).Return(nil, process.ErrRegisteredProcessNotFound)
	s.processRepo.EXPECT().Create(ctx, productID, customMatcher).Return(nil)
	s.objectStorage.EXPECT().UploadImageSources(ctx, productID, expectedRegisteredProcess.Image, expectedBytes).Return(nil)
	s.processRepo.EXPECT().Update(gomock.Any(), productID, expectedUpdatedProcess).Return(errors.New("update error"))
	s.versionService.EXPECT().RegisterProcess(gomock.Any(), productID, expectedRegisteredProcess.ID, expectedRegisteredProcess.Image).
		Return("", nil)
	s.objectStorage.EXPECT().DeleteImageSources(ctx, productID, expectedRegisteredProcess.Image).
		RunAndReturn(func(_ context.Context, _, _ string) error {
			wg.Done()
			return nil
		}).Once()

	returnedProcess, err := s.processService.RegisterProcess(
		ctx, user,
		process.RegisterProcessOpts{
			Product:     productID,
			Version:     version,
			Process:     processName,
			ProcessType: processType,
			Sources:     testFile,
		},
	)
	s.Require().NoError(err)
	s.Equal(expectedRegisteredProcess, returnedProcess)

	s.Require().NoError(testhelpers.WaitOrTimeout(&wg, 1*time.Second))
}

func (s *ProcessServiceTestSuite) TestRegisterProcess_Public() {
	var (
		ctx                       = context.Background()
		wg                        = sync.WaitGroup{}
		processWithCreatingStatus = s.getTestProcess(_publicRegistry, entity.RegisterProcessStatusCreating)
		processWithCreatedStatus  = s.getTestProcess(_publicRegistry, entity.RegisterProcessStatusCreated)
	)

	wg.Add(1)

	testFile, err := os.Open(testFileAddr)
	s.Require().NoError(err)
	expectedBytes, err := os.ReadFile(testFileAddr)
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

	returnedProcess, err := s.processService.RegisterProcess(
		ctx, user,
		process.RegisterProcessOpts{
			Version:     version,
			Process:     processName,
			ProcessType: processType,
			Sources:     testFile,
			IsPublic:    true,
		},
	)
	s.Require().NoError(err)

	s.Equal(processWithCreatingStatus, returnedProcess)
	s.Require().NoError(testhelpers.WaitOrTimeout(&wg, 1*time.Second))
}
