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

func (s *versionSuite) TestStop_OK() {
	// GIVEN a valid user and version
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithTag(versionTag).
		WithStatus(entity.VersionStatusStarted).
		Build()

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStopVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)

	s.natsManagerService.EXPECT().DeleteStreams(ctx, productID, versionTag).Return(nil)
	s.natsManagerService.EXPECT().DeleteObjectStores(ctx, productID, versionTag).Return(nil)
	s.natsManagerService.EXPECT().DeleteVersionKeyValueStores(ctx, productID, vers).Return(nil)
	s.versionRepo.EXPECT().SetStatus(ctx, productID, vers.Tag, entity.VersionStatusStopping).Return(nil)

	// go rutine expected to be called
	s.versionService.EXPECT().Stop(gomock.Any(), productID, vers).Return(nil)
	s.versionRepo.EXPECT().SetStatus(gomock.Any(), productID, vers.Tag, entity.VersionStatusStopped).Return(nil)
	s.userActivityInteractor.EXPECT().RegisterStopAction(user.Email, productID, vers, "testing").Return(nil)

	// WHEN stopping the version
	stoppingVer, notifyChn, err := s.handler.Stop(ctx, user, productID, versionTag, "testing")
	s.NoError(err)

	// THEN the version status is stopping
	vers.Status = entity.VersionStatusStopping
	s.Equal(vers, stoppingVer)

	// THEN the version status when the go rutine ends is stopped
	versionStatus := <-notifyChn
	s.Equal(entity.VersionStatusStopped, versionStatus.Status)
}

func (s *versionSuite) TestStop_ErrorUserNotAuthorized() {
	// GIVEN an unauthorized user and a version
	ctx := context.Background()
	badUser := testhelpers.NewUserBuilder().Build()
	expectedVer := &entity.Version{Tag: versionTag}
	versionMatcher := newVersionMatcher(expectedVer)

	s.accessControl.EXPECT().CheckProductGrants(badUser, productID, auth.ActStopVersion).Return(
		fmt.Errorf("git good"),
	)
	s.userActivityInteractor.EXPECT().RegisterStopAction(badUser.Email, productID, versionMatcher, version.ErrUserNotAuthorized.Error()).Return(nil)

	// WHEN stopping the version
	_, _, err := s.handler.Stop(ctx, badUser, productID, expectedVer.Tag, "testing")

	// THEN an error is returned
	s.Error(err)
}

func (s *versionSuite) TestStop_ErrorVersionNotFound() {
	// GIVEN a valid user and a version not found
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	expectedVer := &entity.Version{Tag: versionTag}
	versionMatcher := newVersionMatcher(expectedVer)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStopVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, expectedVer.Tag).Return(nil, fmt.Errorf("no version found"))
	s.userActivityInteractor.EXPECT().RegisterStopAction(user.Email, productID, versionMatcher, version.ErrVersionNotFound.Error()).Return(nil)

	// WHEN stopping the version
	_, _, err := s.handler.Stop(ctx, user, productID, expectedVer.Tag, "testing")

	// THEN an error is returned
	s.Error(err)
}

func (s *versionSuite) TestStop_ErrorInvalidVersionStatus() {
	// GIVEN a valid user and an invalid version
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithTag(versionTag).
		WithStatus(entity.VersionStatusStopped).
		Build()
	versionMatcher := newVersionMatcher(vers)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStopVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)

	s.userActivityInteractor.EXPECT().RegisterStopAction(user.Email, productID, versionMatcher, version.ErrVersionCannotBeStopped.Error()).Return(nil)

	// WHEN stopping the version
	_, _, err := s.handler.Stop(ctx, user, productID, versionTag, "testing")

	// THEN an error is returned
	s.Error(err)
	s.ErrorIs(err, version.ErrVersionCannotBeStopped)
}

func (s *versionSuite) TestDeleteNatsResources_ErrorDeletingStreams() {
	// GIVEN a valid user and a version
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithTag(versionTag).
		WithStatus(entity.VersionStatusStarted).
		Build()
	versionMatcher := newVersionMatcher(vers)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStopVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)

	s.natsManagerService.EXPECT().DeleteStreams(ctx, productID, versionTag).Return(fmt.Errorf("error deleting streams"))
	s.userActivityInteractor.EXPECT().RegisterStopAction(user.Email, productID, versionMatcher, version.ErrDeletingNATSResources.Error()).Return(nil)

	// WHEN stopping the version
	_, _, err := s.handler.Stop(ctx, user, productID, versionTag, "testing")

	// THEN an error is returned
	s.Error(err)
}

func (s *versionSuite) TestDeleteNatsResources_ErrorDeletingObjectStores() {
	// GIVEN a valid user and a version
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithTag(versionTag).
		WithStatus(entity.VersionStatusStarted).
		Build()
	versionMatcher := newVersionMatcher(vers)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStopVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)

	s.natsManagerService.EXPECT().DeleteStreams(ctx, productID, versionTag).Return(nil)
	s.natsManagerService.EXPECT().DeleteObjectStores(ctx, productID, versionTag).Return(fmt.Errorf("error deleting object stores"))
	s.userActivityInteractor.EXPECT().RegisterStopAction(user.Email, productID, versionMatcher, version.ErrDeletingNATSResources.Error()).Return(nil)

	// WHEN stopping the version
	_, _, err := s.handler.Stop(ctx, user, productID, versionTag, "testing")

	// THEN an error is returned
	s.Error(err)
}

func (s *versionSuite) TestStop_CheckNonBlockingErrorLogging() {
	// GIVEN a valid user and version
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithTag(versionTag).
		WithStatus(entity.VersionStatusStarted).
		Build()

	setStatusErrStarting := errors.New("set status error")
	setStatusErrStarted := errors.New("not again")
	registerActionErr := errors.New("this is the end")

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStopVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)

	s.natsManagerService.EXPECT().DeleteStreams(ctx, productID, versionTag).Return(nil)
	s.natsManagerService.EXPECT().DeleteObjectStores(ctx, productID, versionTag).Return(nil)
	s.natsManagerService.EXPECT().DeleteVersionKeyValueStores(ctx, productID, vers).Return(nil)

	// GIVEN first set status errors
	s.versionRepo.EXPECT().SetStatus(ctx, productID, vers.Tag, entity.VersionStatusStopping).
		Return(setStatusErrStarting)

	// go rutine expected calls
	s.versionService.EXPECT().Stop(gomock.Any(), productID, vers).Return(nil)
	// GIVEN second set status errors
	s.versionRepo.EXPECT().SetStatus(gomock.Any(), productID, vers.Tag, entity.VersionStatusStopped).
		Return(setStatusErrStarted)
	// GIVEN register stop action errors
	s.userActivityInteractor.EXPECT().RegisterStopAction(user.Email, productID, vers, "testing").
		Return(registerActionErr)

	// WHEN stopping the version
	stoppingVer, notifyChn, err := s.handler.Stop(ctx, user, productID, versionTag, "testing")
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
	s.Equal(log1.ContextMap()["error"], setStatusErrStarting.Error())
	log2 := s.observedLogs.All()[2]
	s.Equal(log2.ContextMap()["error"], setStatusErrStarted.Error())
	log3 := s.observedLogs.All()[3]
	s.Equal(log3.ContextMap()["error"], registerActionErr.Error())
}

func (s *versionSuite) TestStop_ErrorUserNotAuthorized_ErrorRegisterAction() {
	// GIVEN an unauthorized user and a version
	ctx := context.Background()
	badUser := testhelpers.NewUserBuilder().Build()
	expectedVer := &entity.Version{Tag: versionTag}
	versionMatcher := newVersionMatcher(expectedVer)

	customErr := errors.New("oh no")
	registerActionErr := errors.New("a bad day")

	s.accessControl.EXPECT().CheckProductGrants(badUser, productID, auth.ActStopVersion).Return(
		customErr,
	)
	// Given error registering action
	s.userActivityInteractor.EXPECT().RegisterStopAction(badUser.Email, productID, versionMatcher, version.ErrUserNotAuthorized.Error()).Return(
		registerActionErr,
	)

	// WHEN stopping the version
	_, _, err := s.handler.Stop(ctx, badUser, productID, expectedVer.Tag, "testing")

	// THEN an error is returned
	s.Error(err)

	// THEN failed registered action is logged
	s.Require().Len(s.observedLogs.All(), 1)
	log1 := s.observedLogs.All()[0]
	s.Equal(log1.ContextMap()["error"], registerActionErr.Error())
}

func (s *versionSuite) TestStopAndNotify_ErrorVersionServiceStop() {
	// GIVEN a valid user and version
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithTag(versionTag).
		WithStatus(entity.VersionStatusStarted).
		Build()
	errStoppingVersion := "error stopping version"
	setErrorErr := errors.New("error setting error")

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStopVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)

	s.natsManagerService.EXPECT().DeleteStreams(ctx, productID, versionTag).Return(nil)
	s.natsManagerService.EXPECT().DeleteObjectStores(ctx, productID, versionTag).Return(nil)
	s.natsManagerService.EXPECT().DeleteVersionKeyValueStores(ctx, productID, vers).Return(nil)
	s.versionRepo.EXPECT().SetStatus(ctx, productID, vers.Tag, entity.VersionStatusStopping).Return(nil)

	// go rutine expected to be called
	s.versionService.EXPECT().Stop(gomock.Any(), productID, vers).Return(fmt.Errorf(errStoppingVersion))
	s.userActivityInteractor.EXPECT().RegisterStopAction(user.Email, productID, vers, version.ErrStoppingVersion.Error()).Return(nil)

	// Given set status
	s.versionRepo.EXPECT().SetErrorStatusWithError(gomock.Any(), productID, vers.Tag, errStoppingVersion).Return(setErrorErr)

	// WHEN stopping the version
	stoppingVer, notifyChn, err := s.handler.Stop(ctx, user, productID, versionTag, "testing")
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
	s.Equal(log1.ContextMap()["error"], setErrorErr.Error())
}
