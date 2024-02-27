//go:build unit

package version_test

import (
	"context"
	"errors"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/version"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
)

const (
	_globalKeyValueStore = "test-global-kv-store"
	_waitGroupTimeout    = 500 * time.Millisecond
)

var (
	prod = &entity.Product{
		ID:            _productID,
		KeyValueStore: _globalKeyValueStore,
	}
)

func (s *versionSuite) TestStart_OK() {
	// GIVEN a valid user and version
	var (
		ctx  = context.Background()
		user = testhelpers.NewUserBuilder().Build()
		vers = testhelpers.NewVersionWithConfigsBuilder().
			WithStatus(entity.VersionStatusCreated).
			Build()

		versionStreamResources = s.getVersionStreamingResources(vers)
		keyValueStoreResources = versionStreamResources.KeyValueStores

		workflow = vers.Workflows[0]

		configurationsToUpdate = s.getTestKeyValueConfigurations(keyValueStoreResources, vers, workflow, workflow.Processes[0])
	)

	s.accessControl.EXPECT().CheckProductGrants(user, _productID, auth.ActManageVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, _productID, _versionTag).Return(vers, nil)
	s.productRepo.EXPECT().GetByID(ctx, _productID).Times(1).Return(prod, nil)

	s.natsManagerService.EXPECT().CreateStreams(gomock.Any(), _productID, vers).Return(versionStreamResources.Streams, nil)
	s.natsManagerService.EXPECT().CreateObjectStores(gomock.Any(), _productID, vers).Return(versionStreamResources.ObjectStores, nil)
	s.natsManagerService.EXPECT().CreateVersionKeyValueStores(gomock.Any(), _productID, vers).Return(keyValueStoreResources, nil)
	s.natsManagerService.EXPECT().UpdateKeyValueConfiguration(gomock.Any(), configurationsToUpdate).Return(nil)
	s.versionRepo.EXPECT().SetStatus(ctx, _productID, vers.Tag, entity.VersionStatusStarting).Return(nil)

	// goroutine calls
	s.versionService.EXPECT().Start(gomock.Any(), prod, vers, versionStreamResources).Return(nil)
	s.versionRepo.EXPECT().SetStatus(gomock.Any(), _productID, vers.Tag, entity.VersionStatusStarted).Return(nil)
	s.userActivityInteractor.EXPECT().RegisterStartAction(user.Email, _productID, vers, "testing").Return(nil)

	// WHEN starting the version
	startingVer, notifyCh, err := s.handler.Start(ctx, user, _productID, _versionTag, "testing")
	s.Require().NoError(err)

	// THEN
	s.Equal(entity.VersionStatusStarting, startingVer.Status)

	startedVersion, ok := <-notifyCh
	s.Require().True(ok)

	s.Equal(entity.VersionStatusStarted, startedVersion.Status)
}

func (s *versionSuite) TestStart_ErrorUserNotAuthorized() {
	// GIVEN an unauthorized user and a version
	ctx := context.Background()
	badUser := testhelpers.NewUserBuilder().Build()
	expectedVer := &entity.Version{Tag: _versionTag}

	expectedError := errors.New("unauthorized")

	s.accessControl.EXPECT().CheckProductGrants(badUser, _productID, auth.ActManageVersion).Return(expectedError)

	// WHEN starting the version
	_, _, err := s.handler.Start(ctx, badUser, _productID, expectedVer.Tag, "testing")

	// THEN an error is returned
	s.ErrorIs(err, expectedError)
}

func (s *versionSuite) TestStart_ErrorCriticalVersionUnauthorized() {
	// GIVEN an unauthorized user and a version
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithStatus(entity.VersionStatusCritical).
		Build()

	expectedError := errors.New("unauthorized")

	s.accessControl.EXPECT().CheckProductGrants(user, prod.ID, auth.ActManageVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, prod.ID).Return(prod, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, prod.ID, vers.Tag).Return(vers, nil)
	s.accessControl.EXPECT().CheckProductGrants(user, prod.ID, auth.ActManageCriticalVersion).Return(expectedError)

	// WHEN starting the version
	_, _, err := s.handler.Start(ctx, user, _productID, vers.Tag, "testing")

	// THEN an error is returned
	s.ErrorIs(err, expectedError)
}

func (s *versionSuite) TestStart_ErrorNonExistingVersion() {
	// GIVEN a valid user and a non-existent version
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	expectedVer := &entity.Version{Tag: _versionTag}

	expectedError := errors.New("version repo error")

	s.accessControl.EXPECT().CheckProductGrants(user, _productID, auth.ActManageVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, _productID).Return(prod, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, _productID, _versionTag).Return(nil, expectedError)

	// WHEN starting the version
	_, _, err := s.handler.Start(ctx, user, _productID, expectedVer.Tag, "testing")

	// THEN an error is returned
	s.ErrorIs(err, expectedError)
}

func (s *versionSuite) TestStart_ErrorInvalidVersionStatus() {
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithTag(_versionTag).
		WithStatus(entity.VersionStatusStarted).
		Build()

	s.accessControl.EXPECT().CheckProductGrants(user, _productID, auth.ActManageVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, _productID).Return(prod, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, _productID, _versionTag).Return(vers, nil)

	// WHEN starting the version
	_, _, err := s.handler.Start(ctx, user, _productID, _versionTag, "testing")

	// THEN an error is returned
	s.ErrorIs(err, version.ErrVersionCannotBeStarted)
}

func (s *versionSuite) TestStart_ErrorCreatingStreams() {
	// GIVEN a valid user and a non-existent version
	var (
		ctx  = context.Background()
		user = testhelpers.NewUserBuilder().Build()
		vers = testhelpers.NewVersionBuilder().
			WithTag(_versionTag).
			WithStatus(entity.VersionStatusCreated).
			Build()

		expectedError = errors.New("stream creation error")
		errStrMatcher = newStringContainsMatcher(expectedError.Error())
	)

	s.accessControl.EXPECT().CheckProductGrants(user, _productID, auth.ActManageVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, _productID).Times(1).Return(prod, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, _productID, _versionTag).Return(vers, nil)
	s.versionRepo.EXPECT().SetStatus(ctx, _productID, _versionTag, entity.VersionStatusStarting).Return(nil)

	s.natsManagerService.EXPECT().CreateStreams(gomock.Any(), _productID, vers).Return(nil, expectedError)
	s.versionRepo.EXPECT().SetErrorStatusWithError(ctx, _productID, vers.Tag, errStrMatcher).Return(nil)

	// WHEN starting the version
	_, notifyCh, err := s.handler.Start(ctx, user, _productID, _versionTag, "testing")
	s.Require().NoError(err)

	failedVersion, ok := <-notifyCh
	s.Require().True(ok)

	s.Equal(entity.VersionStatusError, failedVersion.Status)
	s.Contains(failedVersion.Error, expectedError.Error())
}

func (s *versionSuite) TestStart_ErrorCreatingObjectStore() {
	// GIVEN a valid user and a version
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithTag(_versionTag).
		WithStatus(entity.VersionStatusCreated).
		Build()

	expectedError := errors.New("error creating object-stores")

	errStrMatcher := newStringContainsMatcher(expectedError.Error())

	s.accessControl.EXPECT().CheckProductGrants(user, _productID, auth.ActManageVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, _productID).Times(1).Return(prod, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, _productID, _versionTag).Return(vers, nil)
	s.versionRepo.EXPECT().SetStatus(ctx, _productID, _versionTag, entity.VersionStatusStarting).Return(nil)

	s.natsManagerService.EXPECT().CreateStreams(gomock.Any(), _productID, vers).Return(nil, nil)
	s.natsManagerService.EXPECT().CreateObjectStores(gomock.Any(), _productID, vers).Return(nil, expectedError)
	// Compensation calls
	s.natsManagerService.EXPECT().DeleteStreams(gomock.Any(), _productID, vers.Tag).Return(nil)
	s.versionRepo.EXPECT().SetErrorStatusWithError(gomock.Any(), _productID, vers.Tag, errStrMatcher).Return(nil)

	// WHEN starting the version
	_, notifyCh, err := s.handler.Start(ctx, user, _productID, _versionTag, "testing")
	s.Require().NoError(err)

	failedVersion, ok := <-notifyCh
	s.Require().True(ok)

	s.Equal(entity.VersionStatusError, failedVersion.Status)
	s.Contains(failedVersion.Error, expectedError.Error())
}

func (s *versionSuite) TestStart_ErrorCreatingKeyValueStores() {
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithTag(_versionTag).
		WithStatus(entity.VersionStatusCreated).
		Build()

	expectedError := errors.New("error creating key-value store")

	errStrMatcher := newStringContainsMatcher(expectedError.Error())

	s.accessControl.EXPECT().CheckProductGrants(user, _productID, auth.ActManageVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, _productID, _versionTag).Return(vers, nil)
	s.productRepo.EXPECT().GetByID(ctx, _productID).Times(1).Return(prod, nil)
	s.versionRepo.EXPECT().SetStatus(ctx, _productID, _versionTag, entity.VersionStatusStarting).Return(nil)

	s.natsManagerService.EXPECT().CreateStreams(gomock.Any(), _productID, vers).Return(nil, nil)
	s.natsManagerService.EXPECT().CreateObjectStores(gomock.Any(), _productID, vers).Return(nil, nil)
	s.natsManagerService.EXPECT().CreateVersionKeyValueStores(gomock.Any(), _productID, vers).Return(nil, expectedError)
	s.natsManagerService.EXPECT().DeleteObjectStores(gomock.Any(), _productID, vers.Tag).Return(nil)
	s.natsManagerService.EXPECT().DeleteStreams(gomock.Any(), _productID, vers.Tag).Return(nil)
	s.versionRepo.EXPECT().SetErrorStatusWithError(gomock.Any(), _productID, vers.Tag, errStrMatcher).Return(nil)

	_, notifyCh, err := s.handler.Start(ctx, user, _productID, _versionTag, "testing")
	s.Require().NoError(err)

	failedVersion, ok := <-notifyCh
	s.Require().True(ok)

	s.Equal(entity.VersionStatusError, failedVersion.Status)
	s.Contains(failedVersion.Error, expectedError.Error())
}

func (s *versionSuite) TestStart_ErrorVersionServiceStart() {
	// GIVEN a valid user and version
	var (
		ctx  = context.Background()
		user = testhelpers.NewUserBuilder().Build()
		vers = testhelpers.NewVersionBuilder().
			WithTag(_versionTag).
			WithStatus(entity.VersionStatusCreated).
			Build()

		expectedError   = errors.New("error starting version")
		errStrMatcher   = newStringContainsMatcher(expectedError.Error())
		streamResources = s.getVersionStreamingResources(vers)
	)

	s.accessControl.EXPECT().CheckProductGrants(user, _productID, auth.ActManageVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, _productID, _versionTag).Return(vers, nil)
	s.productRepo.EXPECT().GetByID(ctx, _productID).Times(1).Return(prod, nil)
	s.versionRepo.EXPECT().SetStatus(ctx, _productID, vers.Tag, entity.VersionStatusStarting).Return(nil)

	s.natsManagerService.EXPECT().CreateStreams(gomock.Any(), _productID, vers).Return(streamResources.Streams, nil)
	s.natsManagerService.EXPECT().CreateObjectStores(gomock.Any(), _productID, vers).Return(streamResources.ObjectStores, nil)
	s.natsManagerService.EXPECT().CreateVersionKeyValueStores(gomock.Any(), _productID, vers).Return(streamResources.KeyValueStores, nil)

	s.versionService.EXPECT().Start(gomock.Any(), prod, vers, streamResources).
		Return(expectedError)

	s.natsManagerService.EXPECT().DeleteObjectStores(gomock.Any(), _productID, vers.Tag).Return(nil)
	s.natsManagerService.EXPECT().DeleteStreams(gomock.Any(), _productID, vers.Tag).Return(nil)
	s.natsManagerService.EXPECT().DeleteVersionKeyValueStores(gomock.Any(), _productID, vers).Return(nil)

	s.versionRepo.EXPECT().SetErrorStatusWithError(gomock.Any(), _productID, vers.Tag, errStrMatcher).Return(nil)

	// WHEN starting the version
	startingVer, notifyCh, err := s.handler.Start(ctx, user, _productID, _versionTag, "testing")
	s.NoError(err)

	// THEN the version status first is starting
	vers.Status = entity.VersionStatusStarting
	s.Equal(vers, startingVer)

	failedVersion, ok := <-notifyCh
	s.Require().True(ok)

	s.Equal(entity.VersionStatusError, failedVersion.Status)
	s.Contains(failedVersion.Error, expectedError.Error())
}

func (s *versionSuite) TestStart_ErrorRegisteringUserActivity() {
	// GIVEN a valid user and version
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithTag(_versionTag).
		WithStatus(entity.VersionStatusCreated).
		Build()

	comment := "testing"

	expectedError := errors.New("error registering user activity")

	errStrMatcher := newStringContainsMatcher(expectedError.Error())

	streamResources := s.getVersionStreamingResources(vers)

	s.accessControl.EXPECT().CheckProductGrants(user, _productID, auth.ActManageVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, _productID, _versionTag).Return(vers, nil)
	s.productRepo.EXPECT().GetByID(ctx, _productID).Times(1).Return(prod, nil)
	s.versionRepo.EXPECT().SetStatus(ctx, _productID, vers.Tag, entity.VersionStatusStarting).Return(nil)

	s.natsManagerService.EXPECT().CreateStreams(gomock.Any(), _productID, vers).Return(streamResources.Streams, nil)
	s.natsManagerService.EXPECT().CreateObjectStores(gomock.Any(), _productID, vers).Return(streamResources.ObjectStores, nil)
	s.natsManagerService.EXPECT().CreateVersionKeyValueStores(gomock.Any(), _productID, vers).Return(streamResources.KeyValueStores, nil)
	s.versionRepo.EXPECT().SetStatus(gomock.Any(), _productID, vers.Tag, entity.VersionStatusStarted).Return(nil)
	s.versionService.EXPECT().Start(gomock.Any(), prod, vers, streamResources).
		Return(nil)

	s.userActivityInteractor.EXPECT().RegisterStartAction(user.Email, _productID, vers, comment).Return(expectedError)

	s.natsManagerService.EXPECT().DeleteObjectStores(gomock.Any(), _productID, vers.Tag).Return(nil)
	s.natsManagerService.EXPECT().DeleteStreams(gomock.Any(), _productID, vers.Tag).Return(nil)
	s.natsManagerService.EXPECT().DeleteVersionKeyValueStores(gomock.Any(), _productID, vers).Return(nil)
	s.versionService.EXPECT().Stop(gomock.Any(), _productID, vers).Return(nil)
	s.versionRepo.EXPECT().SetErrorStatusWithError(gomock.Any(), _productID, vers.Tag, errStrMatcher).Return(nil)

	// WHEN starting the version
	startingVer, notifyCh, err := s.handler.Start(ctx, user, _productID, _versionTag, comment)
	s.NoError(err)

	// THEN the version status first is starting
	vers.Status = entity.VersionStatusStarting
	s.Equal(vers, startingVer)

	failedVersion, ok := <-notifyCh
	s.Require().True(ok)

	s.Equal(entity.VersionStatusError, failedVersion.Status)
	s.Contains(failedVersion.Error, expectedError.Error())
}

func (s *versionSuite) getVersionStreamingResources(vers *entity.Version) *entity.VersionStreamingResources {
	s.Require().Greater(len(vers.Workflows), 0)
	s.Require().Greater(len(vers.Workflows[0].Processes), 0)

	workflow := vers.Workflows[0]
	process := workflow.Processes[0]

	versionKeyValueStores := &entity.KeyValueStores{
		GlobalKeyValueStore:  _globalKeyValueStore,
		VersionKeyValueStore: "version-kv-store",
		Workflows: map[string]*entity.WorkflowKeyValueStores{
			workflow.Name: {
				KeyValueStore: "workflow-kv-store",
				Processes: map[string]string{
					process.Name: "process-kv-store",
				},
			},
		},
	}

	return &entity.VersionStreamingResources{
		KeyValueStores: versionKeyValueStores,
		Streams: &entity.VersionStreams{
			Workflows: map[string]entity.WorkflowStreamResources{
				workflow.Name: {
					Stream: "workflow-stream",
					Processes: map[string]entity.ProcessStreamConfig{
						process.Name: {
							Subject:       "process-subject",
							Subscriptions: nil,
						},
					},
				},
			},
		},
	}
}

func (s *versionSuite) getTestKeyValueConfigurations(
	keyValueStoreResources *entity.KeyValueStores,
	vers *entity.Version,
	workflow entity.Workflow,
	process entity.Process,
) []entity.KeyValueConfiguration {
	configurationsToUpdate := []entity.KeyValueConfiguration{
		{
			Store:         keyValueStoreResources.VersionKeyValueStore,
			Configuration: vers.Config,
		},
		{
			Store:         keyValueStoreResources.Workflows[workflow.Name].KeyValueStore,
			Configuration: workflow.Config,
		},
		{
			Store:         keyValueStoreResources.Workflows[workflow.Name].Processes[process.Name],
			Configuration: process.Config,
		},
	}

	return configurationsToUpdate
}
