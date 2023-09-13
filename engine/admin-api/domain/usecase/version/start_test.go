//go:build unit

package version_test

import (
	"context"
	"fmt"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/version"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/version/utils"
	internalerrors "github.com/konstellation-io/kai/engine/admin-api/internal/errors"
)

func (s *VersionUsecaseTestSuite) TestStart_OK() {
	// GIVEN a valid user and version
	ctx := context.Background()
	user := s.getTestUser()
	vers := utils.InitTestVersion().
		WithVersionID(versionID).
		WithTag(versionTag).
		WithStatus(entity.VersionStatusCreated).
		GetVersion()

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStartVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)

	s.natsManagerService.EXPECT().CreateStreams(ctx, productID, vers).Return(nil, nil)
	s.natsManagerService.EXPECT().CreateObjectStores(ctx, productID, vers).Return(nil, nil)
	s.natsManagerService.EXPECT().CreateKeyValueStores(ctx, productID, vers).Return(nil, nil)
	s.versionRepo.EXPECT().SetStatus(ctx, productID, vers.ID, entity.VersionStatusStarting).Return(nil)

	expectedVersionConfig := &entity.VersionConfig{}

	// go rutine expected to be called
	s.versionService.EXPECT().Start(gomock.Any(), productID, vers, expectedVersionConfig).Return(nil)
	s.versionRepo.EXPECT().SetStatus(gomock.Any(), productID, vers.ID, entity.VersionStatusStarted).Return(nil)
	s.userActivityInteractor.EXPECT().RegisterStartAction(user.ID, productID, vers, "testing").Return(nil)

	// WHEN starting the version
	startingVer, notifyChn, err := s.handler.Start(ctx, user, productID, vers.Tag, "testing")
	s.NoError(err)

	// THEN the version status first is starting
	vers.Status = entity.VersionStatusStarting
	s.Equal(vers, startingVer)

	// THEN the version status when the go rutine ends is started
	versionStatus := <-notifyChn
	s.Equal(entity.VersionStatusStarted, versionStatus.Status)
}

func (s *VersionUsecaseTestSuite) TestStart_ErrorUserNotAuthorized() {
	// GIVEN an unauthorized user and a version
	ctx := context.Background()
	badUser := s.getTestUser()
	expectedVer := &entity.Version{Tag: versionTag}
	versionMatcher := newVersionMatcher(expectedVer)

	s.accessControl.EXPECT().CheckProductGrants(badUser, productID, auth.ActStartVersion).Return(
		fmt.Errorf("git good"),
	)
	s.userActivityInteractor.EXPECT().RegisterStartAction(badUser.ID, productID, versionMatcher, version.CommentUserNotAuthorized).Return(nil)

	// WHEN starting the version
	_, _, err := s.handler.Start(ctx, badUser, productID, expectedVer.Tag, "testing")

	// THEN an error is returned
	s.Error(err)
}

func (s *VersionUsecaseTestSuite) TestStart_ErrorNonExistingVersion() {
	// GIVEN a valid user and a non existent version
	ctx := context.Background()
	user := s.getTestUser()
	expectedVer := &entity.Version{Tag: versionTag}
	versionMatcher := newVersionMatcher(expectedVer)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStartVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(nil, fmt.Errorf("no version"))
	s.userActivityInteractor.EXPECT().RegisterStartAction(user.ID, productID, versionMatcher, version.CommentVersionNotFound).Return(nil)

	// WHEN starting the version
	_, _, err := s.handler.Start(ctx, user, productID, expectedVer.Tag, "testing")
	s.Error(err)
}

func (s *VersionUsecaseTestSuite) TestStart_ErrorInvalidVersionStatus() {
	// GIVEN a valid user and a non existent version
	ctx := context.Background()
	user := s.getTestUser()
	vers := utils.InitTestVersion().
		WithVersionID(versionID).
		WithTag(versionTag).
		WithStatus(entity.VersionStatusStarted).
		GetVersion()
	versionMatcher := newVersionMatcher(vers)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStartVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)

	s.userActivityInteractor.EXPECT().RegisterStartAction(user.ID, productID, versionMatcher, version.CommentInvalidVersionStatus).Return(nil)

	// WHEN starting the version
	_, _, err := s.handler.Start(ctx, user, productID, vers.Tag, "testing")
	s.Error(err)
	s.ErrorIs(err, internalerrors.ErrInvalidVersionStatusBeforeStarting)
}

func (s *VersionUsecaseTestSuite) TestStart_ErrorGetVersionConfig_CreateStreams() {
	// GIVEN a valid user and a non existent version
	ctx := context.Background()
	user := s.getTestUser()
	vers := utils.InitTestVersion().
		WithVersionID(versionID).
		WithTag(versionTag).
		WithStatus(entity.VersionStatusCreated).
		GetVersion()
	versionMatcher := newVersionMatcher(vers)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStartVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)

	s.natsManagerService.EXPECT().CreateStreams(ctx, productID, vers).Return(nil, fmt.Errorf("error creating streams"))

	s.userActivityInteractor.EXPECT().RegisterStartAction(user.ID, productID, versionMatcher, version.CommentErrorCreatingNATSResources).Return(nil)

	// WHEN starting the version
	_, _, err := s.handler.Start(ctx, user, productID, vers.Tag, "testing")
	s.Error(err)
}

func (s *VersionUsecaseTestSuite) TestStart_ErrorGetVersionConfig_CreateObjectStore() {
	// GIVEN a valid user and a non existent version
	ctx := context.Background()
	user := s.getTestUser()
	vers := utils.InitTestVersion().
		WithVersionID(versionID).
		WithTag(versionTag).
		WithStatus(entity.VersionStatusCreated).
		GetVersion()
	versionMatcher := newVersionMatcher(vers)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStartVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)

	s.natsManagerService.EXPECT().CreateStreams(ctx, productID, vers).Return(nil, nil)
	s.natsManagerService.EXPECT().CreateObjectStores(ctx, productID, vers).Return(nil, fmt.Errorf("error creating object stores"))

	s.userActivityInteractor.EXPECT().RegisterStartAction(user.ID, productID, versionMatcher, version.CommentErrorCreatingNATSResources).Return(nil)

	// WHEN starting the version
	_, _, err := s.handler.Start(ctx, user, productID, vers.Tag, "testing")
	s.Error(err)
}

func (s *VersionUsecaseTestSuite) TestStart_ErrorGetVersionConfig_CreateKeyValueStore() {
	// GIVEN a valid user and a non existent version
	ctx := context.Background()
	user := s.getTestUser()
	vers := utils.InitTestVersion().
		WithVersionID(versionID).
		WithTag(versionTag).
		WithStatus(entity.VersionStatusCreated).
		GetVersion()
	versionMatcher := newVersionMatcher(vers)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStartVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)

	s.natsManagerService.EXPECT().CreateStreams(ctx, productID, vers).Return(nil, nil)
	s.natsManagerService.EXPECT().CreateObjectStores(ctx, productID, vers).Return(nil, nil)
	s.natsManagerService.EXPECT().CreateKeyValueStores(ctx, productID, vers).Return(nil, fmt.Errorf("error creating key value stores"))

	s.userActivityInteractor.EXPECT().RegisterStartAction(user.ID, productID, versionMatcher, version.CommentErrorCreatingNATSResources).Return(nil)

	// WHEN starting the version
	_, _, err := s.handler.Start(ctx, user, productID, vers.Tag, "testing")
	s.Error(err)
}

func (s *VersionUsecaseTestSuite) TestStart_CheckNonBlockingErrorLogging() {
	// GIVEN a valid user and version
	ctx := context.Background()
	user := s.getTestUser()
	vers := utils.InitTestVersion().
		WithVersionID(versionID).
		WithTag(versionTag).
		WithStatus(entity.VersionStatusCreated).
		GetVersion()

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStartVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)

	s.natsManagerService.EXPECT().CreateStreams(ctx, productID, vers).Return(nil, nil)
	s.natsManagerService.EXPECT().CreateObjectStores(ctx, productID, vers).Return(nil, nil)
	s.natsManagerService.EXPECT().CreateKeyValueStores(ctx, productID, vers).Return(nil, nil)

	// GIVEN first set status errors
	s.versionRepo.EXPECT().SetStatus(ctx, productID, vers.ID, entity.VersionStatusStarting).
		Return(fmt.Errorf("hello this error"))

	expectedVersionConfig := &entity.VersionConfig{}

	// go rutine expecected calls
	s.versionService.EXPECT().Start(gomock.Any(), productID, vers, expectedVersionConfig).Return(nil)
	// GIVEN second set status errors
	s.versionRepo.EXPECT().SetStatus(gomock.Any(), productID, vers.ID, entity.VersionStatusStarted).
		Return(fmt.Errorf("no, this is patrick"))
	// GIVEN register start action errors
	s.userActivityInteractor.EXPECT().RegisterStartAction(user.ID, productID, vers, "testing").
		Return(fmt.Errorf("this is sparta remix"))

	// WHEN starting the version
	startingVer, notifyChn, err := s.handler.Start(ctx, user, productID, vers.Tag, "testing")
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
	s.Equal(log1.ContextMap()["error"], version.ErrUpdatingVersionStatus.Error())
	log2 := s.observedLogs.All()[2]
	s.Equal(log2.ContextMap()["error"], version.ErrUpdatingVersionStatus.Error())
	log3 := s.observedLogs.All()[3]
	s.Equal(log3.ContextMap()["error"], version.ErrRegisteringUserActivity.Error())
}

func (s *VersionUsecaseTestSuite) TestStart_ErrorUserNotAuthorized_ErrorRegisterAction() {
	// GIVEN an unauthorized user and a version
	ctx := context.Background()
	badUser := s.getTestUser()
	expectedVer := &entity.Version{Tag: versionTag}
	versionMatcher := newVersionMatcher(expectedVer)

	s.accessControl.EXPECT().CheckProductGrants(badUser, productID, auth.ActStartVersion).Return(
		fmt.Errorf("git good"),
	)
	// GIVEN error registering action
	s.userActivityInteractor.EXPECT().RegisterStartAction(badUser.ID, productID, versionMatcher, version.CommentUserNotAuthorized).
		Return(fmt.Errorf("also failed"))

	// WHEN starting the version
	_, _, err := s.handler.Start(ctx, badUser, productID, expectedVer.Tag, "testing")

	// THEN an error is returned
	s.Error(err)

	// THEN failed registered action is logged
	s.Require().Len(s.observedLogs.All(), 2)
	log1 := s.observedLogs.All()[1]
	s.Equal(log1.ContextMap()["error"], version.ErrRegisteringUserActivity.Error())
}

func (s *VersionUsecaseTestSuite) TestStart_ErrorVersionServiceStart() {
	// GIVEN a valid user and version
	ctx := context.Background()
	user := s.getTestUser()
	vers := utils.InitTestVersion().
		WithVersionID(versionID).
		WithTag(versionTag).
		WithStatus(entity.VersionStatusCreated).
		GetVersion()
	errStartingVersion := "error starting version"

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStartVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)

	s.natsManagerService.EXPECT().CreateStreams(ctx, productID, vers).Return(nil, nil)
	s.natsManagerService.EXPECT().CreateObjectStores(ctx, productID, vers).Return(nil, nil)
	s.natsManagerService.EXPECT().CreateKeyValueStores(ctx, productID, vers).Return(nil, nil)
	s.versionRepo.EXPECT().SetStatus(ctx, productID, vers.ID, entity.VersionStatusStarting).Return(nil)

	expectedVersionConfig := &entity.VersionConfig{}

	// go rutine expecected calls
	s.versionService.EXPECT().Start(gomock.Any(), productID, vers, expectedVersionConfig).
		Return(fmt.Errorf(errStartingVersion))
	s.userActivityInteractor.EXPECT().RegisterStartAction(user.ID, productID, vers, version.CommentErrorStartingVersion).Return(nil)

	// GIVEN set status

	s.versionRepo.EXPECT().SetError(gomock.Any(), productID, vers, errStartingVersion).
		Return(nil, fmt.Errorf("bomb rush crew"))

	// WHEN starting the version
	startingVer, notifyChn, err := s.handler.Start(ctx, user, productID, vers.Tag, "testing")
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
	s.Equal(log1.ContextMap()["error"], version.ErrUpdatingVersionError.Error())
}
