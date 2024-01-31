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
)

func (s *versionSuite) TestUnpublish_OK() {
	// GIVEN a valid user and published version
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithTag(_versionTag).
		WithStatus(entity.VersionStatusPublished).
		Build()
	product := testhelpers.NewProductBuilder().
		WithPublishedVersion(&vers.Tag).
		Build()

	s.accessControl.EXPECT().CheckProductGrants(user, _productID, auth.ActUnpublishVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, _productID, _versionTag).Return(vers, nil)
	s.productRepo.EXPECT().GetByID(ctx, _productID).Return(product, nil)

	s.versionService.EXPECT().Unpublish(ctx, _productID, vers).Return(nil)
	s.versionRepo.EXPECT().Update(_productID, vers).Return(nil)
	s.productRepo.EXPECT().Update(ctx, product).Return(nil)
	s.userActivityInteractor.EXPECT().RegisterUnpublishAction(user.Email, _productID, vers, "unpublishing").Return(nil)

	// WHEN unpublishing the version
	unpublishedVer, err := s.handler.Unpublish(ctx, user, _productID, _versionTag, "unpublishing")

	// THEN the version status is started, publication fields are cleared, and it's not published
	s.NoError(err)
	s.Equal(entity.VersionStatusStarted, unpublishedVer.Status)
	s.Nil(unpublishedVer.PublicationAuthor)
	s.Nil(unpublishedVer.PublicationDate)
	s.False(product.HasVersionPublished()) // product is a pointer, so it's updated
}

func (s *versionSuite) TestUnpublish_ErrorUserNotAuthorized() {
	// GIVEN an unauthorized user and a published version
	ctx := context.Background()
	badUser := testhelpers.NewUserBuilder().Build()
	expectedVer := &entity.Version{Tag: _versionTag}

	s.accessControl.EXPECT().CheckProductGrants(badUser, _productID, auth.ActUnpublishVersion).Return(
		fmt.Errorf("git good"),
	)

	// WHEN unpublishing the version
	_, err := s.handler.Unpublish(ctx, badUser, _productID, expectedVer.Tag, "unpublishing")

	// THEN an error is returned
	s.Error(err)
}

func (s *versionSuite) TestUnpublish_ErrorVersionNotFound() {
	// GIVEN a valid user and a version not found
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	expectedVer := &entity.Version{Tag: _versionTag}

	s.accessControl.EXPECT().CheckProductGrants(user, _productID, auth.ActUnpublishVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, _productID, expectedVer.Tag).Return(nil, fmt.Errorf("no version found"))

	// WHEN unpublishing the version
	_, err := s.handler.Unpublish(ctx, user, _productID, expectedVer.Tag, "unpublishing")

	// THEN an error is returned
	s.Error(err)
}

func (s *versionSuite) TestUnpublish_ErrorVersionCannotBeUnpublished() {
	// GIVEN a valid user and a version that cannot be unpublished
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithTag(_versionTag).
		WithStatus(entity.VersionStatusStarted).
		Build()

	s.accessControl.EXPECT().CheckProductGrants(user, _productID, auth.ActUnpublishVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, _productID, _versionTag).Return(vers, nil)

	// WHEN unpublishing the version
	_, err := s.handler.Unpublish(ctx, user, _productID, _versionTag, "unpublishing")

	// THEN an error is returned
	s.ErrorIs(err, version.ErrVersionCannotBeUnpublished)
}

func (s *versionSuite) TestUnpublish_ErrorProductNotFound() {
	// GIVEN a valid user and a published version, but product not found
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithTag(_versionTag).
		WithStatus(entity.VersionStatusPublished).
		Build()

	s.accessControl.EXPECT().CheckProductGrants(user, _productID, auth.ActUnpublishVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, _productID, _versionTag).Return(vers, nil)
	s.productRepo.EXPECT().GetByID(ctx, _productID).Return(nil, fmt.Errorf("no product found"))

	s.versionService.EXPECT().Unpublish(ctx, _productID, vers).Return(nil)
	s.versionRepo.EXPECT().Update(_productID, vers).Return(nil)

	// WHEN unpublishing the version
	_, err := s.handler.Unpublish(ctx, user, _productID, _versionTag, "unpublishing")

	// THEN an error is returned
	s.Error(err)
}

func (s *versionSuite) TestUnpublish_ErrorUnpublishingVersion() {
	// GIVEN a valid user and a published version, but error during unpublishing
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithTag(_versionTag).
		WithStatus(entity.VersionStatusPublished).
		Build()
	unpubErr := errors.New("unpublish error in k8s service")

	s.accessControl.EXPECT().CheckProductGrants(user, _productID, auth.ActUnpublishVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, _productID, _versionTag).Return(vers, nil)

	s.versionService.EXPECT().Unpublish(ctx, _productID, vers).Return(unpubErr)

	// WHEN unpublishing the version
	_, err := s.handler.Unpublish(ctx, user, _productID, _versionTag, "unpublishing")

	// THEN an error is returned
	s.ErrorIs(err, version.ErrUnpublishingVersion)
}

func (s *versionSuite) TestUnpublish_CheckNonBlockingErrorLogging() {
	// GIVEN a valid user and a published version, but error during unpublishing
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithTag(_versionTag).
		WithStatus(entity.VersionStatusPublished).
		Build()
	product := testhelpers.NewProductBuilder().
		WithPublishedVersion(&vers.Tag).
		Build()

	setStatusErr := errors.New("error updating version status")
	updateProductErr := errors.New("error updating product published version")
	registerActionErr := errors.New("error registering unpublish action")

	s.accessControl.EXPECT().CheckProductGrants(user, _productID, auth.ActUnpublishVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, _productID, _versionTag).Return(vers, nil)
	s.productRepo.EXPECT().GetByID(ctx, _productID).Return(product, nil)

	s.versionService.EXPECT().Unpublish(ctx, _productID, vers).Return(nil)
	s.versionRepo.EXPECT().Update(_productID, vers).Return(setStatusErr)
	s.productRepo.EXPECT().Update(ctx, product).Return(updateProductErr)
	s.userActivityInteractor.EXPECT().RegisterUnpublishAction(user.Email, _productID, vers, "unpublishing").Return(registerActionErr)

	// WHEN unpublishing the version
	_, err := s.handler.Unpublish(ctx, user, _productID, _versionTag, "unpublishing")
	s.NoError(err)

	s.Require().Len(s.observedLogs.All(), 4)
	print(s.observedLogs.All())
	log1 := s.observedLogs.All()[1]
	s.Equal(log1.ContextMap()["error"], setStatusErr.Error())
	log2 := s.observedLogs.All()[2]
	s.Equal(log2.ContextMap()["error"], updateProductErr.Error())
	log3 := s.observedLogs.All()[3]
	s.Equal(log3.ContextMap()["error"], registerActionErr.Error())

}
