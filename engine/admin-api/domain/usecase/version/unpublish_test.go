//go:build unit

package version_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/version"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
	"github.com/stretchr/testify/assert"
)

func (s *versionSuite) TestUnpublish_OK() {
	// GIVEN a valid user and published version
	ctx := context.Background()
	user := s.getTestUser()
	vers := testhelpers.NewVersionBuilder().
		WithTag(versionTag).
		WithStatus(entity.VersionStatusPublished).
		Build()

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActUnpublishVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)

	s.versionService.EXPECT().Unpublish(ctx, productID, vers).Return(nil)
	s.versionRepo.EXPECT().Update(productID, vers).Return(nil)
	s.userActivityInteractor.EXPECT().RegisterUnpublishAction(user.ID, productID, vers, "unpublishing").Return(nil)

	// WHEN unpublishing the version
	unpublishedVer, err := s.handler.Unpublish(ctx, user, productID, versionTag, "unpublishing")

	// THEN the version status is started, publication fields are cleared, and it's not published
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), entity.VersionStatusStarted, unpublishedVer.Status)
	assert.Nil(s.T(), unpublishedVer.PublicationAuthor)
	assert.Nil(s.T(), unpublishedVer.PublicationDate)
}

func (s *versionSuite) TestUnpublish_ErrorUserNotAuthorized() {
	// GIVEN an unauthorized user and a published version
	ctx := context.Background()
	badUser := s.getTestUser()
	expectedVer := &entity.Version{Tag: versionTag}
	versionMatcher := newVersionMatcher(expectedVer)

	s.accessControl.EXPECT().CheckProductGrants(badUser, productID, auth.ActUnpublishVersion).Return(
		fmt.Errorf("git good"),
	)
	s.userActivityInteractor.EXPECT().RegisterUnpublishAction(badUser.ID, productID, versionMatcher, version.ErrUserNotAuthorized.Error()).Return(nil)

	// WHEN unpublishing the version
	_, err := s.handler.Unpublish(ctx, badUser, productID, expectedVer.Tag, "unpublishing")

	// THEN an error is returned
	assert.Error(s.T(), err)
}

func (s *versionSuite) TestUnpublish_ErrorVersionNotFound() {
	// GIVEN a valid user and a version not found
	ctx := context.Background()
	user := s.getTestUser()
	expectedVer := &entity.Version{Tag: versionTag}
	versionMatcher := newVersionMatcher(expectedVer)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActUnpublishVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, expectedVer.Tag).Return(nil, fmt.Errorf("no version found"))
	s.userActivityInteractor.EXPECT().RegisterUnpublishAction(user.ID, productID, versionMatcher, version.ErrVersionNotFound.Error()).Return(nil)

	// WHEN unpublishing the version
	_, err := s.handler.Unpublish(ctx, user, productID, expectedVer.Tag, "unpublishing")

	// THEN an error is returned
	assert.Error(s.T(), err)
}

func (s *versionSuite) TestUnpublish_ErrorVersionCannotBeUnpublished() {
	// GIVEN a valid user and a version that cannot be unpublished
	ctx := context.Background()
	user := s.getTestUser()
	vers := testhelpers.NewVersionBuilder().
		WithTag(versionTag).
		WithStatus(entity.VersionStatusStarted).
		Build()
	versionMatcher := newVersionMatcher(vers)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActUnpublishVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)

	s.userActivityInteractor.EXPECT().RegisterUnpublishAction(user.ID, productID, versionMatcher, version.ErrVersionCannotBeUnpublished.Error()).Return(nil)

	// WHEN unpublishing the version
	_, err := s.handler.Unpublish(ctx, user, productID, versionTag, "unpublishing")

	// THEN an error is returned
	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, version.ErrVersionCannotBeUnpublished)
}

func (s *versionSuite) TestUnpublish_ErrorUnpublishingVersion() {
	// GIVEN a valid user and a published version, but error during unpublishing
	ctx := context.Background()
	user := s.getTestUser()
	vers := testhelpers.NewVersionBuilder().
		WithTag(versionTag).
		WithStatus(entity.VersionStatusPublished).
		Build()
	versionMatcher := newVersionMatcher(vers)
	unpubErr := errors.New("unpublish error in k8s service")

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActUnpublishVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)

	s.versionService.EXPECT().Unpublish(ctx, productID, vers).Return(unpubErr)

	s.userActivityInteractor.EXPECT().RegisterUnpublishAction(user.ID, productID, versionMatcher, version.ErrUnpublishingVersion.Error()).Return(nil)

	// WHEN unpublishing the version
	_, err := s.handler.Unpublish(ctx, user, productID, versionTag, "unpublishing")

	// THEN an error is returned
	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, version.ErrUnpublishingVersion)
}

func (s *versionSuite) TestUnpublish_CheckNonBlockingErrorLogging() {
	// GIVEN a valid user and a published version, but error during unpublishing
	ctx := context.Background()
	user := s.getTestUser()
	vers := testhelpers.NewVersionBuilder().
		WithTag(versionTag).
		WithStatus(entity.VersionStatusPublished).
		Build()

	setStatusErr := errors.New("error updating version status")
	registerActionErr := errors.New("error registering unpublish action")

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActUnpublishVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)

	s.versionService.EXPECT().Unpublish(ctx, productID, vers).Return(nil)
	s.versionRepo.EXPECT().Update(productID, vers).Return(setStatusErr)
	s.userActivityInteractor.EXPECT().RegisterUnpublishAction(user.ID, productID, vers, "unpublishing").Return(registerActionErr)

	// WHEN unpublishing the version
	_, err := s.handler.Unpublish(ctx, user, productID, versionTag, "unpublishing")
	s.NoError(err)

	s.Require().Len(s.observedLogs.All(), 3)
	print(s.observedLogs.All())
	log1 := s.observedLogs.All()[1]
	s.Equal(log1.ContextMap()["error"], setStatusErr.Error())
	log2 := s.observedLogs.All()[2]
	s.Equal(log2.ContextMap()["error"], registerActionErr.Error())
}
