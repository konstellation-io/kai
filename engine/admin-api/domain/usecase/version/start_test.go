//go:build unit

package version_test

import (
	"context"
	"errors"
	"sync"
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
		ID:            productID,
		KeyValueStore: _globalKeyValueStore,
	}
)

func (s *versionSuite) TestStart_OK() {
	// GIVEN a valid user and version
	var (
		ctx  = context.Background()
		wg   = sync.WaitGroup{}
		user = testhelpers.NewUserBuilder().Build()
		vers = testhelpers.NewVersionWithConfigsBuilder().
			WithStatus(entity.VersionStatusCreated).
			Build()

		versionStreamResources = s.getVersionStreamingResources(vers)
		keyValueStoreResources = versionStreamResources.KeyValueStores

		workflow = vers.Workflows[0]

		configurationsToUpdate = s.getTestKeyValueConfigurations(keyValueStoreResources, vers, workflow, workflow.Processes[0])
	)

	wg.Add(1)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStartVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)
	s.productRepo.EXPECT().GetByID(ctx, productID).Times(1).Return(prod, nil)

	s.natsManagerService.EXPECT().CreateStreams(gomock.Any(), productID, vers).Return(versionStreamResources.Streams, nil)
	s.natsManagerService.EXPECT().CreateObjectStores(gomock.Any(), productID, vers).Return(versionStreamResources.ObjectStores, nil)
	s.natsManagerService.EXPECT().CreateVersionKeyValueStores(gomock.Any(), productID, vers).Return(keyValueStoreResources, nil)
	s.natsManagerService.EXPECT().UpdateKeyValueConfiguration(gomock.Any(), configurationsToUpdate).Return(nil)
	s.versionRepo.EXPECT().SetStatus(ctx, productID, vers.Tag, entity.VersionStatusStarting).Return(nil)

	// goroutine calls
	s.versionService.EXPECT().Start(gomock.Any(), prod, vers, versionStreamResources).Return(nil)
	s.versionRepo.EXPECT().SetStatus(gomock.Any(), productID, vers.Tag, entity.VersionStatusStarted).
		DoAndReturn(func(a1, a2, a3, a4 interface{}) error {
			wg.Done()
			return nil
		})
	s.userActivityInteractor.EXPECT().RegisterStartAction(user.Email, productID, vers, "testing").Return(nil)

	// WHEN starting the version
	startingVer, err := s.handler.Start(ctx, user, productID, versionTag, "testing")
	s.Require().NoError(err)

	// THEN
	s.Equal(entity.VersionStatusStarting, startingVer.Status)

	s.Require().NoError(testhelpers.WaitOrTimeout(&wg, _waitGroupTimeout))
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

func (s *versionSuite) TestStart_ErrorUserNotAuthorized() {
	// GIVEN an unauthorized user and a version
	ctx := context.Background()
	badUser := testhelpers.NewUserBuilder().Build()
	expectedVer := &entity.Version{Tag: versionTag}

	expectedError := errors.New("unauthorized")

	s.accessControl.EXPECT().CheckProductGrants(badUser, productID, auth.ActStartVersion).Return(expectedError)

	// WHEN starting the version
	_, err := s.handler.Start(ctx, badUser, productID, expectedVer.Tag, "testing")

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

	s.accessControl.EXPECT().CheckProductGrants(user, prod.ID, auth.ActStartVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, prod.ID).Return(prod, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, prod.ID, vers.Tag).Return(vers, nil)
	s.accessControl.EXPECT().CheckProductGrants(user, prod.ID, auth.ActStartCriticalVersion).Return(expectedError)

	// WHEN starting the version
	_, err := s.handler.Start(ctx, user, productID, vers.Tag, "testing")

	// THEN an error is returned
	s.ErrorIs(err, expectedError)
}

func (s *versionSuite) TestStart_ErrorNonExistingVersion() {
	// GIVEN a valid user and a non-existent version
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	expectedVer := &entity.Version{Tag: versionTag}

	expectedError := errors.New("version repo error")

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStartVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, productID).Return(prod, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(nil, expectedError)

	// WHEN starting the version
	_, err := s.handler.Start(ctx, user, productID, expectedVer.Tag, "testing")

	// THEN an error is returned
	s.ErrorIs(err, expectedError)
}

func (s *versionSuite) TestStart_ErrorInvalidVersionStatus() {
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithTag(versionTag).
		WithStatus(entity.VersionStatusStarted).
		Build()

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStartVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, productID).Return(prod, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)

	// WHEN starting the version
	_, err := s.handler.Start(ctx, user, productID, versionTag, "testing")

	// THEN an error is returned
	s.ErrorIs(err, version.ErrVersionCannotBeStarted)
}

func (s *versionSuite) TestStart_ErrorCreatingStreams() {
	// GIVEN a valid user and a non-existent version
	var (
		ctx  = context.Background()
		wg   = sync.WaitGroup{}
		user = testhelpers.NewUserBuilder().Build()
		vers = testhelpers.NewVersionBuilder().
			WithTag(versionTag).
			WithStatus(entity.VersionStatusCreated).
			Build()

		expectedError = errors.New("stream creation error")
		errStrMatcher = newStringContainsMatcher(expectedError.Error())
	)

	wg.Add(1)
	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStartVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, productID).Times(1).Return(prod, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)
	s.versionRepo.EXPECT().SetStatus(ctx, productID, versionTag, entity.VersionStatusStarting).Return(nil)

	s.natsManagerService.EXPECT().CreateStreams(gomock.Any(), productID, vers).Return(nil, expectedError)
	s.versionRepo.EXPECT().SetErrorStatusWithError(ctx, productID, vers.Tag, errStrMatcher).
		DoAndReturn(func(a1, a2, a3, a4 interface{}) error {
			wg.Done()
			return nil
		})

	// WHEN starting the version
	_, err := s.handler.Start(ctx, user, productID, versionTag, "testing")
	s.Require().NoError(err)

	s.Require().NoError(testhelpers.WaitOrTimeout(&wg, _waitGroupTimeout))
}

func (s *versionSuite) TestStart_ErrorCreatingObjectStore() {
	// GIVEN a valid user and a version
	ctx := context.Background()
	wg := sync.WaitGroup{}
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithTag(versionTag).
		WithStatus(entity.VersionStatusCreated).
		Build()

	expectedError := errors.New("error creating object-stores")

	errStrMatcher := newStringContainsMatcher(expectedError.Error())

	wg.Add(1)
	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStartVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, productID).Times(1).Return(prod, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)
	s.versionRepo.EXPECT().SetStatus(ctx, productID, versionTag, entity.VersionStatusStarting).Return(nil)

	s.natsManagerService.EXPECT().CreateStreams(gomock.Any(), productID, vers).Return(nil, nil)
	s.natsManagerService.EXPECT().CreateObjectStores(gomock.Any(), productID, vers).Return(nil, expectedError)
	// Compensation calls
	s.natsManagerService.EXPECT().DeleteStreams(gomock.Any(), productID, vers.Tag).Return(nil)
	s.versionRepo.EXPECT().SetErrorStatusWithError(gomock.Any(), productID, vers.Tag, errStrMatcher).
		DoAndReturn(func(a1, a2, a3, a4 interface{}) error {
			wg.Done()
			return nil
		})

	// WHEN starting the version
	_, err := s.handler.Start(ctx, user, productID, versionTag, "testing")
	s.Require().NoError(err)

	s.Require().NoError(testhelpers.WaitOrTimeout(&wg, _waitGroupTimeout))

}

func (s *versionSuite) TestStart_ErrorCreatingKeyValueStores() {
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithTag(versionTag).
		WithStatus(entity.VersionStatusCreated).
		Build()

	expectedError := errors.New("error creating key-value store")

	errStrMatcher := newStringContainsMatcher(expectedError.Error())

	wg := sync.WaitGroup{}
	wg.Add(1)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStartVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)
	s.productRepo.EXPECT().GetByID(ctx, productID).Times(1).Return(prod, nil)
	s.versionRepo.EXPECT().SetStatus(ctx, productID, versionTag, entity.VersionStatusStarting).Return(nil)

	s.natsManagerService.EXPECT().CreateStreams(gomock.Any(), productID, vers).Return(nil, nil)
	s.natsManagerService.EXPECT().CreateObjectStores(gomock.Any(), productID, vers).Return(nil, nil)
	s.natsManagerService.EXPECT().CreateVersionKeyValueStores(gomock.Any(), productID, vers).Return(nil, expectedError)
	s.natsManagerService.EXPECT().DeleteObjectStores(gomock.Any(), productID, vers.Tag).Return(nil)
	s.natsManagerService.EXPECT().DeleteStreams(gomock.Any(), productID, vers.Tag).Return(nil)
	s.versionRepo.EXPECT().SetErrorStatusWithError(gomock.Any(), productID, vers.Tag, errStrMatcher).
		DoAndReturn(func(arg1, arg2, arg3, arg4 interface{}) error {
			wg.Done()
			return nil
		})

	_, err := s.handler.Start(ctx, user, productID, versionTag, "testing")
	s.Require().NoError(err)

	s.Require().NoError(testhelpers.WaitOrTimeout(&wg, _waitGroupTimeout))
}

func (s *versionSuite) TestStart_ErrorVersionServiceStart() {
	// GIVEN a valid user and version
	var (
		ctx  = context.Background()
		wg   = sync.WaitGroup{}
		user = testhelpers.NewUserBuilder().Build()
		vers = testhelpers.NewVersionBuilder().
			WithTag(versionTag).
			WithStatus(entity.VersionStatusCreated).
			Build()

		expectedError   = errors.New("error starting version")
		errStrMatcher   = newStringContainsMatcher(expectedError.Error())
		streamResources = s.getVersionStreamingResources(vers)
	)

	wg.Add(1)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStartVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)
	s.productRepo.EXPECT().GetByID(ctx, productID).Times(1).Return(prod, nil)
	s.versionRepo.EXPECT().SetStatus(ctx, productID, vers.Tag, entity.VersionStatusStarting).Return(nil)

	s.natsManagerService.EXPECT().CreateStreams(gomock.Any(), productID, vers).Return(streamResources.Streams, nil)
	s.natsManagerService.EXPECT().CreateObjectStores(gomock.Any(), productID, vers).Return(streamResources.ObjectStores, nil)
	s.natsManagerService.EXPECT().CreateVersionKeyValueStores(gomock.Any(), productID, vers).Return(streamResources.KeyValueStores, nil)

	s.versionService.EXPECT().Start(gomock.Any(), prod, vers, streamResources).
		Return(expectedError)

	s.natsManagerService.EXPECT().DeleteObjectStores(gomock.Any(), productID, vers.Tag).Return(nil)
	s.natsManagerService.EXPECT().DeleteStreams(gomock.Any(), productID, vers.Tag).Return(nil)

	s.versionRepo.EXPECT().SetErrorStatusWithError(gomock.Any(), productID, vers.Tag, errStrMatcher).
		DoAndReturn(func(a1, a2, a3, a4 interface{}) error {
			wg.Done()
			return nil
		})

	// WHEN starting the version
	startingVer, err := s.handler.Start(ctx, user, productID, versionTag, "testing")
	s.NoError(err)

	// THEN the version status first is starting
	vers.Status = entity.VersionStatusStarting
	s.Equal(vers, startingVer)

	s.Require().NoError(testhelpers.WaitOrTimeout(&wg, _waitGroupTimeout))
}

func (s *versionSuite) TestStart_ErrorRegisteringUserActivity() {
	// GIVEN a valid user and version
	ctx := context.Background()
	wg := sync.WaitGroup{}
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithTag(versionTag).
		WithStatus(entity.VersionStatusCreated).
		Build()

	comment := "testing"

	expectedError := errors.New("error registering user activity")

	errStrMatcher := newStringContainsMatcher(expectedError.Error())

	streamResources := s.getVersionStreamingResources(vers)

	wg.Add(1)
	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStartVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)
	s.productRepo.EXPECT().GetByID(ctx, productID).Times(1).Return(prod, nil)
	s.versionRepo.EXPECT().SetStatus(ctx, productID, vers.Tag, entity.VersionStatusStarting).Return(nil)

	s.natsManagerService.EXPECT().CreateStreams(gomock.Any(), productID, vers).Return(streamResources.Streams, nil)
	s.natsManagerService.EXPECT().CreateObjectStores(gomock.Any(), productID, vers).Return(streamResources.ObjectStores, nil)
	s.natsManagerService.EXPECT().CreateVersionKeyValueStores(gomock.Any(), productID, vers).Return(streamResources.KeyValueStores, nil)
	s.versionRepo.EXPECT().SetStatus(gomock.Any(), productID, vers.Tag, entity.VersionStatusStarted).Return(nil)
	s.versionService.EXPECT().Start(gomock.Any(), prod, vers, streamResources).
		Return(nil)

	s.userActivityInteractor.EXPECT().RegisterStartAction(user.Email, productID, vers, comment).Return(expectedError)

	s.natsManagerService.EXPECT().DeleteObjectStores(gomock.Any(), productID, vers.Tag).Return(nil)
	s.natsManagerService.EXPECT().DeleteStreams(gomock.Any(), productID, vers.Tag).Return(nil)
	s.versionService.EXPECT().Stop(gomock.Any(), productID, vers).Return(nil)
	s.versionRepo.EXPECT().SetErrorStatusWithError(gomock.Any(), productID, vers.Tag, errStrMatcher).
		DoAndReturn(func(a1, a2, a3, a4 interface{}) error {
			wg.Done()
			return nil
		})

	// WHEN starting the version
	startingVer, err := s.handler.Start(ctx, user, productID, versionTag, comment)
	s.NoError(err)

	// THEN the version status first is starting
	vers.Status = entity.VersionStatusStarting
	s.Equal(vers, startingVer)

	s.Require().NoError(testhelpers.WaitOrTimeout(&wg, _waitGroupTimeout))
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
