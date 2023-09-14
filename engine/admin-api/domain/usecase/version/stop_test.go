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

const TEST_COMMENT = "test comment"

func (s *VersionUsecaseTestSuite) TestStop_OK() {
	// GIVEN a valid user and version
	ctx := context.Background()
	user := s.getTestUser()
	vers := utils.InitTestVersion().
		WithVersionID(versionID).
		WithTag(versionTag).
		WithStatus(entity.VersionStatusStarted).
		GetVersion()

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStopVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)

	s.natsManagerService.EXPECT().DeleteStreams(ctx, productID, vers.Tag).Return(nil)
	s.natsManagerService.EXPECT().DeleteObjectStores(ctx, productID, vers.Tag).Return(nil)
	s.versionRepo.EXPECT().SetStatus(ctx, productID, vers.ID, entity.VersionStatusStopping).Return(nil)

	// go rutine expected to be called
	s.versionService.EXPECT().Stop(gomock.Any(), productID, vers).Return(nil)
	s.versionRepo.EXPECT().SetStatus(gomock.Any(), productID, vers.ID, entity.VersionStatusStopped).Return(nil)
	s.userActivityInteractor.EXPECT().RegisterStopAction(user.ID, productID, vers, TEST_COMMENT).Return(nil)

	// WHEN stopping the version
	stoppingVer, notifyChn, err := s.handler.Stop(ctx, user, productID, vers.Tag, TEST_COMMENT)
	s.NoError(err)

	// THEN the version status is stopping
	vers.Status = entity.VersionStatusStopping
	s.Equal(vers, stoppingVer)

	// THEN the version status when the go rutine ends is stopped
	versionStatus := <-notifyChn
	s.Equal(entity.VersionStatusStopped, versionStatus.Status)
}

func (s *VersionUsecaseTestSuite) TestStop_ErrorUserNotAuthorized() {
	// GIVEN an unauthorized user and a version
	ctx := context.Background()
	badUser := s.getTestUser()
	expectedVer := &entity.Version{Tag: versionTag}
	versionMatcher := newVersionMatcher(expectedVer)

	s.accessControl.EXPECT().CheckProductGrants(badUser, productID, auth.ActStopVersion).Return(
		fmt.Errorf("git good"),
	)
	s.userActivityInteractor.EXPECT().RegisterStopAction(badUser.ID, productID, versionMatcher, version.CommentUserNotAuthorized).Return(nil)

	// WHEN stopping the version
	_, _, err := s.handler.Stop(ctx, badUser, productID, expectedVer.Tag, TEST_COMMENT)

	// THEN an error is returned
	s.Error(err)
}

func (s *VersionUsecaseTestSuite) TestStop_ErrorVersionNotFound() {
	// GIVEN a valid user and a version not found
	ctx := context.Background()
	user := s.getTestUser()
	expectedVer := &entity.Version{Tag: versionTag}
	versionMatcher := newVersionMatcher(expectedVer)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStopVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, expectedVer.Tag).Return(nil, fmt.Errorf("no version found"))
	s.userActivityInteractor.EXPECT().RegisterStopAction(user.ID, productID, versionMatcher, version.CommentVersionNotFound).Return(nil)

	// WHEN stopping the version
	_, _, err := s.handler.Stop(ctx, user, productID, expectedVer.Tag, TEST_COMMENT)

	// THEN an error is returned
	s.Error(err)
}

func (s *VersionUsecaseTestSuite) TestStop_ErrorInvalidVersionStatus() {
	// GIVEN a valid user and an invalid version
	ctx := context.Background()
	user := s.getTestUser()
	vers := utils.InitTestVersion().
		WithVersionID(versionID).
		WithTag(versionTag).
		WithStatus(entity.VersionStatusStopped).
		GetVersion()
	versionMatcher := newVersionMatcher(vers)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStopVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)

	s.userActivityInteractor.EXPECT().RegisterStopAction(user.ID, productID, versionMatcher, version.CommentInvalidVersionStatusBeforeStopping).Return(nil)

	// WHEN stopping the version
	_, _, err := s.handler.Stop(ctx, user, productID, vers.Tag, TEST_COMMENT)

	// THEN an error is returned
	s.Error(err)
	s.ErrorIs(err, internalerrors.ErrInvalidVersionStatusBeforeStopping)
}

func (s *VersionUsecaseTestSuite) TestDeleteNatsResources_ErrorDeletingStreams() {
	// GIVEN a valid user and a version
	ctx := context.Background()
	user := s.getTestUser()
	vers := utils.InitTestVersion().
		WithVersionID(versionID).
		WithTag(versionTag).
		WithStatus(entity.VersionStatusStarted).
		GetVersion()
	versionMatcher := newVersionMatcher(vers)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStopVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)

	s.natsManagerService.EXPECT().DeleteStreams(ctx, productID, vers.Tag).Return(fmt.Errorf("error deleting streams"))
	s.userActivityInteractor.EXPECT().RegisterStopAction(user.ID, productID, versionMatcher, version.CommentErrorDeletingNATSResources).Return(nil)

	// WHEN stopping the version
	_, _, err := s.handler.Stop(ctx, user, productID, vers.Tag, TEST_COMMENT)

	// THEN an error is returned
	s.Error(err)
}

func (s *VersionUsecaseTestSuite) TestDeleteNatsResources_ErrorDeletingObjectStores() {
	// GIVEN a valid user and a version
	ctx := context.Background()
	user := s.getTestUser()
	vers := utils.InitTestVersion().
		WithVersionID(versionID).
		WithTag(versionTag).
		WithStatus(entity.VersionStatusStarted).
		GetVersion()
	versionMatcher := newVersionMatcher(vers)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStopVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)

	s.natsManagerService.EXPECT().DeleteStreams(ctx, productID, vers.Tag).Return(nil)
	s.natsManagerService.EXPECT().DeleteObjectStores(ctx, productID, vers.Tag).Return(fmt.Errorf("error deleting object stores"))
	s.userActivityInteractor.EXPECT().RegisterStopAction(user.ID, productID, versionMatcher, version.CommentErrorDeletingNATSResources).Return(nil)

	// WHEN stopping the version
	_, _, err := s.handler.Stop(ctx, user, productID, vers.Tag, TEST_COMMENT)

	// THEN an error is returned
	s.Error(err)
}

func (s *VersionUsecaseTestSuite) TestStop_CheckNonBlockingErrorLogging() {
	// GIVEN a valid user and version
	ctx := context.Background()
	user := s.getTestUser()
	vers := utils.InitTestVersion().
		WithVersionID(versionID).
		WithTag(versionTag).
		WithStatus(entity.VersionStatusStarted).
		GetVersion()

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStopVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)

	s.natsManagerService.EXPECT().DeleteStreams(ctx, productID, vers.Tag).Return(nil)
	s.natsManagerService.EXPECT().DeleteObjectStores(ctx, productID, vers.Tag).Return(nil)

	// GIVEN first set status errors
	s.versionRepo.EXPECT().SetStatus(ctx, productID, vers.ID, entity.VersionStatusStopping).
		Return(fmt.Errorf("set status error"))

	// go rutine expected calls
	s.versionService.EXPECT().Stop(gomock.Any(), productID, vers).Return(nil)
	// GIVEN second set status errors
	s.versionRepo.EXPECT().SetStatus(gomock.Any(), productID, vers.ID, entity.VersionStatusStopped).
		Return(fmt.Errorf("not again"))
	// GIVEN register stop action errors
	s.userActivityInteractor.EXPECT().RegisterStopAction(user.ID, productID, vers, TEST_COMMENT).
		Return(fmt.Errorf("this is the end"))

	// WHEN stopping the version
	stoppingVer, notifyChn, err := s.handler.Stop(ctx, user, productID, vers.Tag, TEST_COMMENT)
	s.NoError(err)

	// THEN the version status first is stopping
	vers.Status = entity.VersionStatusStopping
	s.Equal(vers, stoppingVer)

	// THEN the version status when the go rutine ends is stopped
	versionStatus := <-notifyChn
	s.Equal(entity.VersionStatusStopped, versionStatus.Status)

	// THEN both set status are logged
	s.Require().Len(s.observedLogs.All(), 4)
	print(s.observedLogs.All())
	log1 := s.observedLogs.All()[1]
	s.Equal(log1.ContextMap()["error"], version.ErrUpdatingVersionStatus.Error())
	log2 := s.observedLogs.All()[2]
	s.Equal(log2.ContextMap()["error"], version.ErrUpdatingVersionStatus.Error())
	log3 := s.observedLogs.All()[3]
	s.Equal(log3.ContextMap()["error"], version.ErrRegisteringUserActivity.Error())
}

func (s *VersionUsecaseTestSuite) TestStop_ErrorUserNotAuthorized_ErrorRegisterAction() {
	// GIVEN an unauthorized user and a version
	ctx := context.Background()
	badUser := s.getTestUser()
	expectedVer := &entity.Version{Tag: versionTag}
	versionMatcher := newVersionMatcher(expectedVer)

	s.accessControl.EXPECT().CheckProductGrants(badUser, productID, auth.ActStopVersion).Return(
		fmt.Errorf("oh no"),
	)
	// Given error registering action
	s.userActivityInteractor.EXPECT().RegisterStopAction(badUser.ID, productID, versionMatcher, version.CommentUserNotAuthorized).Return(
		fmt.Errorf("a bad day"),
	)

	// WHEN stopping the version
	_, _, err := s.handler.Stop(ctx, badUser, productID, expectedVer.Tag, TEST_COMMENT)

	// THEN an error is returned
	s.Error(err)

	// THEN failed registered action is logged
	s.Require().Len(s.observedLogs.All(), 2)
	log1 := s.observedLogs.All()[1]
	s.Equal(log1.ContextMap()["error"], version.ErrRegisteringUserActivity.Error())
}

func (s *VersionUsecaseTestSuite) TestStopAndNotify_ErrorVersionServiceStop() {
	// GIVEN a valid user and version
	ctx := context.Background()
	user := s.getTestUser()
	vers := utils.InitTestVersion().
		WithVersionID(versionID).
		WithTag(versionTag).
		WithStatus(entity.VersionStatusStarted).
		GetVersion()
	errStoppingVersion := "error stopping version"

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStopVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)

	s.natsManagerService.EXPECT().DeleteStreams(ctx, productID, vers.Tag).Return(nil)
	s.natsManagerService.EXPECT().DeleteObjectStores(ctx, productID, vers.Tag).Return(nil)
	s.versionRepo.EXPECT().SetStatus(ctx, productID, vers.ID, entity.VersionStatusStopping).Return(nil)

	// go rutine expected to be called
	s.versionService.EXPECT().Stop(gomock.Any(), productID, vers).Return(fmt.Errorf(errStoppingVersion))
	s.userActivityInteractor.EXPECT().RegisterStopAction(user.ID, productID, vers, version.CommentErrorStoppingVersion).Return(nil)

	// Given set status
	s.versionRepo.EXPECT().SetError(gomock.Any(), productID, vers, errStoppingVersion).Return(
		nil, fmt.Errorf("error setting error"),
	)

	// WHEN stopping the version
	stoppingVer, notifyChn, err := s.handler.Stop(ctx, user, productID, vers.Tag, TEST_COMMENT)
	s.NoError(err)

	// THEN the version status is stopping
	vers.Status = entity.VersionStatusStopping
	s.Equal(vers, stoppingVer)

	// THEN the version status when the go rutine ends is error
	versionStatus := <-notifyChn
	s.Equal(entity.VersionStatusError, versionStatus.Status)
	s.Equal(errStoppingVersion, versionStatus.Error)

	// THEN set error is logged
	s.Require().Len(s.observedLogs.All(), 2)
	log1 := s.observedLogs.All()[1]
	s.Equal(log1.ContextMap()["error"], version.ErrUpdatingVersionError.Error())
}
