//go:build unit

package version_test

import (
	"context"
	"errors"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/version"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
)

func (s *versionSuite) TestPublish_OK() {
	// GIVEN a valid user and published version
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithTag(versionTag).
		WithStatus(entity.VersionStatusStarted).
		Build()

	expectedURLs := map[string]string{
		"test-trigger": "test-url",
	}

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActPublishVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)

	s.versionService.EXPECT().Publish(ctx, productID, vers.Tag).Return(expectedURLs, nil)
	s.versionRepo.EXPECT().Update(productID, vers).Return(nil)
	s.userActivityInteractor.EXPECT().RegisterPublishAction(user.Email, productID, vers, "publishing").Return(nil)

	// WHEN publishing the version
	urls, err := s.handler.Publish(ctx, user, productID, versionTag, "publishing")

	// THEN the version status is published and the published triggers urls are returned
	s.Require().NoError(err)
	s.Assert().Equal(expectedURLs, urls)

	s.Assert().Equal(user.Email, *vers.PublicationAuthor)
	s.Assert().Equal(entity.VersionStatusPublished, vers.Status)
}

func (s *versionSuite) TestPublishing_ErrorUserNotAuthorized() {
	// GIVEN an unauthorized user and a started version
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	expectedVer := &entity.Version{Tag: versionTag}
	expectedError := errors.New("unauthorized")

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActPublishVersion).Return(expectedError)

	// WHEN publishing the version
	_, err := s.handler.Publish(ctx, user, productID, expectedVer.Tag, "publishing")

	// THEN an error is returned
	s.Assert().ErrorIs(err, expectedError)
}

func (s *versionSuite) TestPublish_ErrorVersionNotFound() {
	// GIVEN a valid user and a version not found
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	expectedVer := testhelpers.NewVersionBuilder().Build()

	expectedError := errors.New("version not found")

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActPublishVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, expectedVer.Tag).Return(nil, expectedError)

	// WHEN unpublishing the version
	_, err := s.handler.Publish(ctx, user, productID, expectedVer.Tag, "publishing")

	// THEN an error is returned
	s.Assert().ErrorIs(err, expectedError)
}

func (s *versionSuite) TestPublish_ErrorVersionCannotBePublished() {
	// GIVEN a valid user and a version that cannot be unpublished
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithTag(versionTag).
		WithStatus(entity.VersionStatusCreated).
		Build()

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActPublishVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)

	// WHEN unpublishing the version
	_, err := s.handler.Publish(ctx, user, productID, versionTag, "publishing")

	// THEN an error is returned
	s.Assert().ErrorIs(err, version.ErrVersionCannotBePublished)
}

func (s *versionSuite) TestPublish_ErrorPublishingVersion() {
	// GIVEN a valid user and a published version, but error during publishing
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithStatus(entity.VersionStatusStarted).
		Build()

	expectedError := errors.New("publish error in k8s service")

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActPublishVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, vers.Tag).Return(vers, nil)

	s.versionService.EXPECT().Publish(ctx, productID, vers.Tag).Return(nil, expectedError)

	// WHEN publishing the version
	_, err := s.handler.Publish(ctx, user, productID, vers.Tag, "publishing")

	// THEN an error is returned
	s.Assert().ErrorIs(err, expectedError)
}

func (s *versionSuite) TestPublish_ErrorInStatusAndRegisteringAction() {
	// GIVEN a valid user and a published version, but error during publishing
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithStatus(entity.VersionStatusStarted).
		Build()

	expectedURLs := map[string]string{
		"test-trigger": "test-url",
	}

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActPublishVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, vers.Tag).Return(vers, nil)

	s.versionService.EXPECT().Publish(ctx, productID, vers.Tag).Return(expectedURLs, nil)

	versionMatcher := newVersionMatcher(vers)

	s.versionRepo.EXPECT().Update(productID, versionMatcher).Return(errors.New("updating version"))
	s.userActivityInteractor.EXPECT().RegisterPublishAction(user.Email, productID, versionMatcher, "publishing").Return(errors.New("registering action"))

	// WHEN publishing the version
	urls, err := s.handler.Publish(ctx, user, productID, vers.Tag, "publishing")

	// THEN an error is returned
	s.Assert().NoError(err)
	s.Assert().Equal(expectedURLs, urls)

	s.Assert().Equal(user.Email, *vers.PublicationAuthor)
	s.Assert().Equal(entity.VersionStatusPublished, vers.Status)
}
