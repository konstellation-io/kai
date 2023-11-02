//go:build unit

package version_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/version"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
)

const (
	_globalKeyValueStore = "test-global-kv-store"
)

var (
	prod = &entity.Product{
		KeyValueStore: _globalKeyValueStore,
	}
)

func (s *versionSuite) TestStart_OK() {
	// GIVEN a valid user and version
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionWithConfigsBuilder().
		WithStatus(entity.VersionStatusCreated).
		Build()

	versionStreamResources := s.getVersionStreamingResources(vers)
	keyValueStoreResources := versionStreamResources.KeyValueStores

	workflow := vers.Workflows[0]
	process := workflow.Processes[0]

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

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStartVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)
	s.productRepo.EXPECT().GetByID(ctx, productID).Times(1).Return(prod, nil)

	s.natsManagerService.EXPECT().CreateStreams(ctx, productID, vers).Return(versionStreamResources.Streams, nil)
	s.natsManagerService.EXPECT().CreateObjectStores(ctx, productID, vers).Return(versionStreamResources.ObjectStores, nil)
	s.natsManagerService.EXPECT().CreateVersionKeyValueStores(ctx, productID, vers).Return(keyValueStoreResources, nil)
	s.natsManagerService.EXPECT().UpdateKeyValueConfiguration(ctx, configurationsToUpdate).Return(nil)
	s.versionRepo.EXPECT().SetStatus(ctx, productID, vers.Tag, entity.VersionStatusStarting).Return(nil)

	// go rutine expected calls
	s.versionService.EXPECT().Start(gomock.Any(), productID, vers, versionStreamResources).Return(nil)
	s.versionRepo.EXPECT().SetStatus(gomock.Any(), productID, vers.Tag, entity.VersionStatusStarted).Return(nil)
	s.userActivityInteractor.EXPECT().RegisterStartAction(user.Email, productID, vers, "testing").Return(nil)

	// WHEN starting the version
	startingVer, notifyChn, err := s.handler.Start(ctx, user, productID, versionTag, "testing")
	s.NoError(err)

	// THEN the version status first is starting
	vers.Status = entity.VersionStatusStarting
	s.Equal(vers, startingVer)

	// THEN the version status when the go rutine ends is started
	versionStatus := <-notifyChn
	s.Equal(entity.VersionStatusStarted, versionStatus.Status)
}

func (s *versionSuite) TestStart_ErrorUserNotAuthorized() {
	// GIVEN an unauthorized user and a version
	ctx := context.Background()
	badUser := testhelpers.NewUserBuilder().Build()
	expectedVer := &entity.Version{Tag: versionTag}
	versionMatcher := newVersionMatcher(expectedVer)

	customErr := errors.New("git good")

	s.accessControl.EXPECT().CheckProductGrants(badUser, productID, auth.ActStartVersion).Return(customErr)
	s.userActivityInteractor.EXPECT().RegisterStartAction(badUser.Email, productID, versionMatcher, version.ErrUserNotAuthorized.Error()).Return(nil)

	// WHEN starting the version
	_, _, err := s.handler.Start(ctx, badUser, productID, expectedVer.Tag, "testing")

	// THEN an error is returned
	s.Error(err)
	s.ErrorIs(err, customErr)
}

func (s *versionSuite) TestStart_ErrorNonExistingVersion() {
	// GIVEN a valid user and a non existent version
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	expectedVer := &entity.Version{Tag: versionTag}
	versionMatcher := newVersionMatcher(expectedVer)

	customErr := errors.New("24h cinderella")

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStartVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(nil, customErr)
	s.userActivityInteractor.EXPECT().RegisterStartAction(user.Email, productID, versionMatcher, version.ErrVersionNotFound.Error()).Return(nil)

	// WHEN starting the version
	_, _, err := s.handler.Start(ctx, user, productID, expectedVer.Tag, "testing")

	// THEN an error is returned
	s.Error(err)
	s.ErrorIs(err, customErr)
}

func (s *versionSuite) TestStart_ErrorInvalidVersionStatus() {
	// GIVEN a valid user and a non existent version
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithTag(versionTag).
		WithStatus(entity.VersionStatusStarted).
		Build()
	versionMatcher := newVersionMatcher(vers)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStartVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)

	s.userActivityInteractor.EXPECT().RegisterStartAction(user.Email, productID, versionMatcher, version.ErrVersionCannotBeStarted.Error()).Return(nil)

	// WHEN starting the version
	_, _, err := s.handler.Start(ctx, user, productID, versionTag, "testing")

	// THEN an error is returned
	s.Error(err)
	s.ErrorIs(err, version.ErrVersionCannotBeStarted)
}

func (s *versionSuite) TestStart_ErrorGetVersionConfig_CreateStreams() {
	// GIVEN a valid user and a non existent version
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithTag(versionTag).
		WithStatus(entity.VersionStatusCreated).
		Build()
	versionMatcher := newVersionMatcher(vers)

	customErr := errors.New("brother Nishiki")

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStartVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)
	s.productRepo.EXPECT().GetByID(ctx, productID).Times(1).Return(prod, nil)

	s.natsManagerService.EXPECT().CreateStreams(ctx, productID, vers).Return(nil, customErr)

	s.userActivityInteractor.EXPECT().RegisterStartAction(user.Email, productID, versionMatcher, version.ErrCreatingNATSResources.Error()).Return(nil)

	// WHEN starting the version
	_, _, err := s.handler.Start(ctx, user, productID, versionTag, "testing")

	// THEN an error is returned
	s.Error(err)
	s.ErrorIs(err, customErr)
}

func (s *versionSuite) TestStart_ErrorGetVersionConfig_CreateObjectStore() {
	// GIVEN a valid user and a non existent version
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithTag(versionTag).
		WithStatus(entity.VersionStatusCreated).
		Build()
	versionMatcher := newVersionMatcher(vers)

	customErr := errors.New("Majima constructions")

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStartVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)
	s.productRepo.EXPECT().GetByID(ctx, productID).Times(1).Return(prod, nil)

	s.natsManagerService.EXPECT().CreateStreams(ctx, productID, vers).Return(nil, nil)
	s.natsManagerService.EXPECT().CreateObjectStores(ctx, productID, vers).Return(nil, customErr)

	s.userActivityInteractor.EXPECT().RegisterStartAction(user.Email, productID, versionMatcher, version.ErrCreatingNATSResources.Error()).Return(nil)

	// WHEN starting the version
	_, _, err := s.handler.Start(ctx, user, productID, versionTag, "testing")

	// THEN an error is returned
	s.Error(err)
	s.ErrorIs(err, customErr)
}

func (s *versionSuite) TestStart_ErrorGetVersionConfig_CreateKeyValueStore() {
	// GIVEN a valid user and a non existent version
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithTag(versionTag).
		WithStatus(entity.VersionStatusCreated).
		Build()
	versionMatcher := newVersionMatcher(vers)

	customErr := errors.New("dame da ne")

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStartVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)
	s.productRepo.EXPECT().GetByID(ctx, productID).Times(1).Return(prod, nil)

	s.natsManagerService.EXPECT().CreateStreams(ctx, productID, vers).Return(nil, nil)
	s.natsManagerService.EXPECT().CreateObjectStores(ctx, productID, vers).Return(nil, nil)
	s.natsManagerService.EXPECT().CreateVersionKeyValueStores(ctx, productID, vers).Return(nil, customErr)

	s.userActivityInteractor.EXPECT().RegisterStartAction(user.Email, productID, versionMatcher, version.ErrCreatingNATSResources.Error()).Return(nil)

	// WHEN starting the version
	_, _, err := s.handler.Start(ctx, user, productID, versionTag, "testing")

	// THEN an error is returned
	s.Error(err)
	s.ErrorIs(err, customErr)
}

func (s *versionSuite) TestStart_CheckNonBlockingErrorLogging() {
	// GIVEN a valid user and version
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithTag(versionTag).
		WithStatus(entity.VersionStatusCreated).
		Build()

	versionStreamResources := s.getVersionStreamingResources(vers)

	setStatusErrStarting := errors.New("hello this error")
	setStatusErrStarted := errors.New("no, this is patrick")
	registerActionErr := errors.New("this is sparta remix")

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStartVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)
	s.productRepo.EXPECT().GetByID(ctx, productID).Times(1).Return(prod, nil)

	s.natsManagerService.EXPECT().CreateStreams(ctx, productID, vers).Return(versionStreamResources.Streams, nil)
	s.natsManagerService.EXPECT().CreateObjectStores(ctx, productID, vers).Return(versionStreamResources.ObjectStores, nil)
	s.natsManagerService.EXPECT().CreateVersionKeyValueStores(ctx, productID, vers).Return(versionStreamResources.KeyValueStores, nil)

	// GIVEN first set status errors
	s.versionRepo.EXPECT().SetStatus(ctx, productID, vers.Tag, entity.VersionStatusStarting).
		Return(setStatusErrStarting)

	// go rutine expected calls
	s.versionService.EXPECT().Start(gomock.Any(), productID, vers, versionStreamResources).Return(nil)
	// GIVEN second set status errors
	s.versionRepo.EXPECT().SetStatus(gomock.Any(), productID, vers.Tag, entity.VersionStatusStarted).
		Return(setStatusErrStarted)
	// GIVEN register start action errors
	s.userActivityInteractor.EXPECT().RegisterStartAction(user.Email, productID, vers, "testing").
		Return(registerActionErr)

	// WHEN starting the version
	startingVer, notifyChn, err := s.handler.Start(ctx, user, productID, versionTag, "testing")
	s.NoError(err)

	// THEN the version status first is starting
	vers.Status = entity.VersionStatusStarting
	s.Equal(vers, startingVer)

	// THEN the version status when the go rutine ends is started
	versionStatus := <-notifyChn
	s.Equal(entity.VersionStatusStarted, versionStatus.Status)

	// THEN both set status are logged
	s.Require().Len(s.observedLogs.All(), 4)
	log1 := s.observedLogs.All()[1]
	s.Equal(log1.ContextMap()["error"], setStatusErrStarting.Error())
	log2 := s.observedLogs.All()[2]
	s.Equal(log2.ContextMap()["error"], setStatusErrStarted.Error())
	log3 := s.observedLogs.All()[3]
	s.Equal(log3.ContextMap()["error"], registerActionErr.Error())
}

func (s *versionSuite) TestStart_ErrorUserNotAuthorized_ErrorRegisterAction() {
	// GIVEN an unauthorized user and a version
	ctx := context.Background()
	badUser := testhelpers.NewUserBuilder().Build()
	expectedVer := &entity.Version{Tag: versionTag}
	versionMatcher := newVersionMatcher(expectedVer)

	customErr := errors.New("git good")
	regiserActionErr := errors.New("also failed")

	s.accessControl.EXPECT().CheckProductGrants(badUser, productID, auth.ActStartVersion).Return(customErr)
	// GIVEN error registering action
	s.userActivityInteractor.EXPECT().RegisterStartAction(badUser.Email, productID, versionMatcher, version.ErrUserNotAuthorized.Error()).Return(regiserActionErr)

	// WHEN starting the version
	_, _, err := s.handler.Start(ctx, badUser, productID, expectedVer.Tag, "testing")

	// THEN an error is returned
	s.Error(err)
	s.ErrorIs(err, customErr)

	// THEN failed registered action is logged
	s.Require().Len(s.observedLogs.All(), 1)
	log1 := s.observedLogs.All()[0]
	s.Equal(log1.ContextMap()["error"], regiserActionErr.Error())
}

func (s *versionSuite) TestStart_ErrorVersionServiceStart() {
	// GIVEN a valid user and version
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithTag(versionTag).
		WithStatus(entity.VersionStatusCreated).
		Build()
	errStartingVersion := "error starting version"
	setErrorErr := errors.New("bomb rush crew")

	streamResources := s.getVersionStreamingResources(vers)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStartVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)
	s.productRepo.EXPECT().GetByID(ctx, productID).Times(1).Return(prod, nil)

	s.natsManagerService.EXPECT().CreateStreams(ctx, productID, vers).Return(streamResources.Streams, nil)
	s.natsManagerService.EXPECT().CreateObjectStores(ctx, productID, vers).Return(streamResources.ObjectStores, nil)
	s.natsManagerService.EXPECT().CreateVersionKeyValueStores(ctx, productID, vers).Return(streamResources.KeyValueStores, nil)
	s.versionRepo.EXPECT().SetStatus(ctx, productID, vers.Tag, entity.VersionStatusStarting).Return(nil)

	// go rutine expected calls
	s.versionService.EXPECT().Start(gomock.Any(), productID, vers, streamResources).
		Return(fmt.Errorf(errStartingVersion))
	s.userActivityInteractor.EXPECT().RegisterStartAction(user.Email, productID, vers, version.ErrStartingVersion.Error()).Return(nil)

	// GIVEN set status
	s.versionRepo.EXPECT().SetError(gomock.Any(), productID, vers, errStartingVersion).Return(nil, setErrorErr)

	// WHEN starting the version
	startingVer, notifyChn, err := s.handler.Start(ctx, user, productID, versionTag, "testing")
	s.NoError(err)

	// THEN the version status first is starting
	vers.Status = entity.VersionStatusStarting
	s.Equal(vers, startingVer)

	// THEN the version status when the go rutine ends is error
	versionStatus := <-notifyChn
	s.Equal(entity.VersionStatusError, versionStatus.Status)
	s.Equal(errStartingVersion, versionStatus.Error)

	// THEN set error is logged
	s.Require().Len(s.observedLogs.All(), 2)
	log1 := s.observedLogs.All()[1]
	s.Equal(log1.ContextMap()["error"], setErrorErr.Error())
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
