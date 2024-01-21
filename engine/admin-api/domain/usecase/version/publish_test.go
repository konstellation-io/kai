//go:build unit

package version_test

import (
	"context"
	"errors"

	"github.com/golang/mock/gomock"
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
	product := testhelpers.NewProductBuilder().Build()

	expectedURLs := map[string]string{
		"test-trigger": "test-url",
	}

	s.accessControl.EXPECT().CheckProductGrants(user, product.ID, auth.ActPublishVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, product.ID).Return(product, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, product.ID, versionTag).Return(vers, nil)

	s.versionService.EXPECT().Publish(gomock.Any(), product.ID, vers.Tag).Return(expectedURLs, nil)
	s.versionRepo.EXPECT().Update(product.ID, vers).Return(nil)
	s.productRepo.EXPECT().Update(gomock.Any(), product).Return(nil)
	s.userActivityInteractor.EXPECT().RegisterPublishAction(user.Email, product.ID, vers, "publishing").Return(nil)

	// WHEN publishing the version
	_, notifyCh, err := s.handler.Publish(ctx, user, version.PublishParams{
		ProductID:  product.ID,
		VersionTag: versionTag,
		Comment:    "publishing",
	})

	// THEN the version status is publishing
	s.Require().NoError(err)
	//s.Assert().Equal(entity.VersionStatusPublishing, publishingVersion.Status)

	publishedVersion, ok := <-notifyCh
	s.Require().True(ok)

	s.Assert().Equal(user.Email, *publishedVersion.PublicationAuthor)
	s.Assert().Equal(entity.VersionStatusPublished, publishedVersion.Status)
}

func (s *versionSuite) TestPublishing_ErrorUserNotAuthorized() {
	// GIVEN an unauthorized user and a started version
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	expectedVer := &entity.Version{Tag: versionTag}
	product := testhelpers.NewProductBuilder().Build()

	expectedError := errors.New("unauthorized")

	s.accessControl.EXPECT().CheckProductGrants(user, product.ID, auth.ActPublishVersion).Return(expectedError)

	// WHEN publishing the version
	_, _, err := s.handler.Publish(ctx, user, version.PublishParams{
		ProductID:  product.ID,
		VersionTag: expectedVer.Tag,
		Comment:    "publishing",
	})

	// THEN an error is returned
	s.Assert().ErrorIs(err, expectedError)
}

func (s *versionSuite) TestPublish_ErrorVersionNotFound() {
	// GIVEN a valid user and a version not found
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	expectedVer := testhelpers.NewVersionBuilder().Build()
	product := testhelpers.NewProductBuilder().Build()

	expectedError := errors.New("version not found")

	s.accessControl.EXPECT().CheckProductGrants(user, product.ID, auth.ActPublishVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, product.ID).Return(product, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, product.ID, expectedVer.Tag).Return(nil, expectedError)

	// WHEN unpublishing the version
	_, _, err := s.handler.Publish(ctx, user, version.PublishParams{
		ProductID:  product.ID,
		VersionTag: expectedVer.Tag,
		Comment:    "publishing",
	})

	// THEN an error is returned
	s.Assert().ErrorIs(err, expectedError)
}

func (s *versionSuite) TestPublish_ErrorVersionCannotBePublished_NoForce() {
	// GIVEN a valid user and a version that cannot be unpublished
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithTag(versionTag).
		WithStatus(entity.VersionStatusCreated).
		Build()
	product := testhelpers.NewProductBuilder().Build()

	s.accessControl.EXPECT().CheckProductGrants(user, product.ID, auth.ActPublishVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, product.ID).Return(product, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, product.ID, versionTag).Return(vers, nil)

	// WHEN unpublishing the version
	_, _, err := s.handler.Publish(ctx, user, version.PublishParams{
		ProductID:  product.ID,
		VersionTag: versionTag,
		Comment:    "publishing",
	})

	// THEN an error is returned
	s.Assert().ErrorIs(err, version.ErrVersionIsNotStarted)
}

func (s *versionSuite) TestPublish_ProductWithVersionAlreadyPublished_NoForce() {
	// GIVEN a valid user and a product with a published version
	var (
		ctx     = context.Background()
		user    = testhelpers.NewUserBuilder().Build()
		product = testhelpers.NewProductBuilder().
			WithPublishedVersion(testhelpers.StrPointer("another-version")).
			Build()
		vers = testhelpers.NewVersionBuilder().
			WithTag(versionTag).
			WithStatus(entity.VersionStatusStarted).
			Build()
	)

	s.accessControl.EXPECT().CheckProductGrants(user, product.ID, auth.ActPublishVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, product.ID).Return(product, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, product.ID, vers.Tag).Return(vers, nil)

	// WHEN publishing the version
	_, _, err := s.handler.Publish(ctx, user, version.PublishParams{
		ProductID:  product.ID,
		VersionTag: vers.Tag,
		Comment:    "publishing",
	})

	// THEN an error is returned
	s.Assert().ErrorIs(err, version.ErrProductAlreadyPublished)
}

func (s *versionSuite) TestPublish_ErrorPublishingVersion() {
	// GIVEN a valid user and a published version, but error during publishing
	var (
		ctx  = context.Background()
		user = testhelpers.NewUserBuilder().Build()
		vers = testhelpers.NewVersionBuilder().
			WithStatus(entity.VersionStatusStarted).
			Build()

		product = testhelpers.NewProductBuilder().Build()

		expectedError = errors.New("publish error in k8s service")
		errStrMatcher = newStringContainsMatcher(expectedError.Error())
	)

	s.accessControl.EXPECT().CheckProductGrants(user, product.ID, auth.ActPublishVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, product.ID, vers.Tag).Return(vers, nil)
	s.productRepo.EXPECT().GetByID(ctx, product.ID).Return(product, nil)

	s.versionService.EXPECT().Publish(gomock.Any(), product.ID, vers.Tag).Return(nil, expectedError)

	s.versionRepo.EXPECT().SetErrorStatusWithError(gomock.Any(), product.ID, vers.Tag, errStrMatcher).Return(nil)

	// WHEN no error in the initial return (the versionService publish is executed if a goroutine)
	_, notifyCh, _ := s.handler.Publish(ctx, user, version.PublishParams{
		ProductID:  product.ID,
		VersionTag: vers.Tag,
		Comment:    "publishing",
	})

	failedVersion, ok := <-notifyCh
	s.Require().True(ok)

	s.Equal(entity.VersionStatusError, failedVersion.Status)
}

func (s *versionSuite) TestPublish_ErrorInStatusAndRegisteringAction() {
	// GIVEN a valid user and a published version, but error during publishing
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithStatus(entity.VersionStatusStarted).
		Build()

	product := testhelpers.NewProductBuilder().
		WithPublishedVersion(nil).
		Build()

	expectedURLs := map[string]string{
		"test-trigger": "test-url",
	}

	expectedError := errors.New("error registering user activity")
	errStrMatcher := newStringContainsMatcher(expectedError.Error())

	s.accessControl.EXPECT().CheckProductGrants(user, product.ID, auth.ActPublishVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, product.ID).Return(product, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, product.ID, vers.Tag).Return(vers, nil)

	s.versionService.EXPECT().Publish(gomock.Any(), product.ID, vers.Tag).Return(expectedURLs, nil)

	versionMatcher := newVersionMatcher(vers)

	s.productRepo.EXPECT().Update(gomock.Any(), product).Times(2).Return(nil)
	s.versionRepo.EXPECT().Update(product.ID, versionMatcher).Times(2).Return(nil)

	s.userActivityInteractor.EXPECT().RegisterPublishAction(user.Email, product.ID, gomock.Any(), "publishing").
		Return(expectedError)

	s.versionService.EXPECT().Unpublish(gomock.Any(), product.ID, versionMatcher).Return(nil)
	// WHEN publishing the version
	_, notifyCh, err := s.handler.Publish(ctx, user, version.PublishParams{
		ProductID:  product.ID,
		VersionTag: vers.Tag,
		Comment:    "publishing",
	})

	s.versionRepo.EXPECT().SetErrorStatusWithError(gomock.Any(), product.ID, vers.Tag, errStrMatcher).Return(nil)

	// THEN no error is returned (error happens in goroutine)
	s.Require().NoError(err)
	//s.Equal(entity.VersionStatusPublishing, v.Status)

	failedVersion, ok := <-notifyCh
	s.Require().True(ok)

	s.Assert().Equal(entity.VersionStatusError, failedVersion.Status)
}

func (s *versionSuite) TestPublish_VersionIsNotStarted_Forced() {
	// GIVEN a valid user and a non started version
	var (
		ctx     = context.Background()
		user    = testhelpers.NewUserBuilder().Build()
		product = testhelpers.NewProductBuilder().Build()
		vers    = testhelpers.NewVersionBuilder().
			WithTag(versionTag).
			WithStatus(entity.VersionStatusCreated).
			Build()
	)

	s.accessControl.EXPECT().CheckProductGrants(user, product.ID, auth.ActPublishVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, product.ID).Return(product, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, product.ID, versionTag).Return(vers, nil)

	s.mockStartVersion(user, product, vers)

	s.versionService.EXPECT().Publish(gomock.Any(), product.ID, vers.Tag).Return(nil, nil)
	s.versionRepo.EXPECT().Update(product.ID, vers).Return(nil)
	s.productRepo.EXPECT().Update(gomock.Any(), product).Return(nil)
	s.userActivityInteractor.EXPECT().RegisterPublishAction(user.Email, product.ID, vers, "publishing").Return(nil)

	// WHEN publish the version with the param Force set to true
	_, notifyCh, err := s.handler.Publish(ctx, user, version.PublishParams{
		ProductID:  product.ID,
		VersionTag: versionTag,
		Comment:    "publishing",
		Force:      true,
	})

	// THEN an error is returned
	s.Require().NoError(err)

	startedVersion, ok := <-notifyCh
	s.Require().True(ok)

	s.Equal(entity.VersionStatusPublished, startedVersion.Status)
}

func (s *versionSuite) TestPublish_AnotherVersionPublished_Forced() {
	// GIVEN a valid user and a non started version
	var (
		ctx        = context.Background()
		user       = testhelpers.NewUserBuilder().Build()
		oldVersion = testhelpers.NewVersionBuilder().
				WithTag("old-version").
				WithStatus(entity.VersionStatusPublished).
				Build()
		product = testhelpers.NewProductBuilder().
			WithPublishedVersion(&oldVersion.Tag).
			Build()
		vers = testhelpers.NewVersionBuilder().
			WithTag(versionTag).
			WithStatus(entity.VersionStatusStarted).
			Build()
	)

	s.accessControl.EXPECT().CheckProductGrants(user, product.ID, auth.ActPublishVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, product.ID).Return(product, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, product.ID, versionTag).Return(vers, nil)

	s.mockUnpublishVersion(user, product, oldVersion)

	s.versionService.EXPECT().Publish(gomock.Any(), product.ID, vers.Tag).Return(nil, nil)
	s.versionRepo.EXPECT().Update(product.ID, vers).Return(nil)
	s.productRepo.EXPECT().Update(gomock.Any(), product).Return(nil)
	s.userActivityInteractor.EXPECT().RegisterPublishAction(user.Email, product.ID, vers, "publishing").Return(nil)

	// WHEN publish the version with the param Force set to true
	_, notifyCh, err := s.handler.Publish(ctx, user, version.PublishParams{
		ProductID:  product.ID,
		VersionTag: versionTag,
		Comment:    "publishing",
		Force:      true,
	})

	// THEN an error is returned
	s.Require().NoError(err)

	startedVersion, ok := <-notifyCh
	s.Require().True(ok)

	s.Equal(entity.VersionStatusPublished, startedVersion.Status)
}

func (s *versionSuite) TestPublish_FailsIfTheVersionIsAlreadyPublished() {
	// GIVEN a valid user and a non started version
	var (
		ctx     = context.Background()
		user    = testhelpers.NewUserBuilder().Build()
		product = testhelpers.NewProductBuilder().
			Build()
		vers = testhelpers.NewVersionBuilder().
			WithTag(versionTag).
			WithStatus(entity.VersionStatusPublished).
			Build()
	)

	s.accessControl.EXPECT().CheckProductGrants(user, product.ID, auth.ActPublishVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, product.ID).Return(product, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, product.ID, versionTag).Return(vers, nil)

	// WHEN publish the version with the param Force set to true
	_, _, err := s.handler.Publish(ctx, user, version.PublishParams{
		ProductID:  product.ID,
		VersionTag: versionTag,
		Comment:    "publishing",
	})

	// THEN an error is returned
	s.Require().ErrorIs(err, version.ErrVersionAlreadyPublished)
}

func (s *versionSuite) TestPublish_FailsIfTheVersionIsBeingPublished() {
	// GIVEN a valid user and a non started version
	var (
		ctx     = context.Background()
		user    = testhelpers.NewUserBuilder().Build()
		product = testhelpers.NewProductBuilder().
			Build()
		vers = testhelpers.NewVersionBuilder().
			WithTag(versionTag).
			WithStatus(entity.VersionStatusPublishing).
			Build()
	)

	s.accessControl.EXPECT().CheckProductGrants(user, product.ID, auth.ActPublishVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, product.ID).Return(product, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, product.ID, versionTag).Return(vers, nil)

	// WHEN publish the version with the param Force set to true
	_, _, err := s.handler.Publish(ctx, user, version.PublishParams{
		ProductID:  product.ID,
		VersionTag: versionTag,
		Comment:    "publishing",
	})

	// THEN an error is returned
	s.Require().ErrorIs(err, version.ErrVersionBeingPublished)
}

func (s *versionSuite) mockStartVersion(user *entity.User, product *entity.Product, vers *entity.Version) {
	versionStreamResources := s.getVersionStreamingResources(vers)
	keyValueStoreResources := versionStreamResources.KeyValueStores

	s.accessControl.EXPECT().CheckProductGrants(user, product.ID, auth.ActStartVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(gomock.Any(), product.ID, versionTag).Return(vers, nil)
	s.productRepo.EXPECT().GetByID(gomock.Any(), product.ID).Times(1).Return(product, nil)

	s.natsManagerService.EXPECT().CreateStreams(gomock.Any(), product.ID, vers).Return(versionStreamResources.Streams, nil)
	s.natsManagerService.EXPECT().CreateObjectStores(gomock.Any(), product.ID, vers).Return(versionStreamResources.ObjectStores, nil)
	s.natsManagerService.EXPECT().CreateVersionKeyValueStores(gomock.Any(), product.ID, vers).Return(keyValueStoreResources, nil)
	s.versionRepo.EXPECT().SetStatus(gomock.Any(), product.ID, vers.Tag, entity.VersionStatusStarting).Return(nil)

	// goroutine calls
	s.versionService.EXPECT().Start(gomock.Any(), product, vers, versionStreamResources).Return(nil)
	s.versionRepo.EXPECT().SetStatus(gomock.Any(), product.ID, vers.Tag, entity.VersionStatusStarted).Return(nil)
	s.userActivityInteractor.EXPECT().RegisterStartAction(user.Email, product.ID, vers, gomock.Any()).Return(nil)
}

func (s *versionSuite) mockUnpublishVersion(user *entity.User, product *entity.Product, vers *entity.Version) {
	s.accessControl.EXPECT().CheckProductGrants(user, product.ID, auth.ActUnpublishVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(gomock.Any(), product.ID, vers.Tag).Return(vers, nil)
	s.productRepo.EXPECT().GetByID(gomock.Any(), product.ID).Return(product, nil)

	s.versionService.EXPECT().Unpublish(gomock.Any(), product.ID, vers).Return(nil)
	s.versionRepo.EXPECT().Update(product.ID, vers).Return(nil)
	s.productRepo.EXPECT().Update(gomock.Any(), product)
	s.userActivityInteractor.EXPECT().RegisterUnpublishAction(user.Email, product.ID, vers, gomock.Any()).Return(nil)
}
