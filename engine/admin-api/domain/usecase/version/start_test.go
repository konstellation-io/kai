//go:build unit

package version_test

import (
	"context"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/version/utils"
)

// TODO refactorizar set errors para solo guardar un mensaje de error
// Poner los 3 user activities para el start
// Quitar el handler_test y mover los tests al test que toque

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

	s.userActivityInteractor.EXPECT().RegisterStartAction(user.ID, productID, vers, "testing").Return(nil)

	s.natsManagerService.EXPECT().CreateStreams(ctx, productID, vers).Return(nil, nil)
	s.natsManagerService.EXPECT().CreateObjectStores(ctx, productID, vers).Return(nil, nil)
	s.natsManagerService.EXPECT().CreateKeyValueStores(ctx, productID, vers).Return(nil, nil)
	s.versionRepo.EXPECT().SetStatus(ctx, productID, vers.ID, entity.VersionStatusStarting).Return(nil)

	expectedVersionConfig := &entity.VersionConfig{}

	// go rutine expecected calls
	s.versionService.EXPECT().Start(gomock.Any(), productID, vers, expectedVersionConfig).Return(nil)
	s.versionRepo.EXPECT().SetStatus(gomock.Any(), productID, vers.ID, entity.VersionStatusStarted).Return(nil)

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
