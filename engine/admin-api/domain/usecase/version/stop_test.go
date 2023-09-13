//go:build unit

package version_test

import (
	"context"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/version/utils"
)

func (s *VersionUsecaseTestSuite) TestStop_OK() {
	// GIVEN a valid user and version
	ctx := context.Background()
	user := s.getTestUser()
	vers := utils.InitTestVersion().
		WithVersionID(versionID).
		WithTag(versionTag).
		WithStatus(entity.VersionStatusStopped).
		GetVersion()

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActStopVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, versionTag).Return(vers, nil)

	s.natsManagerService.EXPECT().DeleteStreams(ctx, productID, vers.Tag).Return(nil)
	s.natsManagerService.EXPECT().DeleteObjectStores(ctx, productID, vers.Tag).Return(nil)
	s.versionRepo.EXPECT().SetStatus(ctx, productID, vers.ID, entity.VersionStatusStopping).Return(nil)

	// go rutine expected to be called
	s.versionService.EXPECT().Stop(gomock.Any(), productID, vers).Return(nil)
	s.versionRepo.EXPECT().SetStatus(gomock.Any(), productID, vers.ID, entity.VersionStatusStopped).Return(nil)
	s.userActivityInteractor.EXPECT().RegisterStopAction(user.ID, productID, vers, "testing").Return(nil)

	// WHEN stopping the version
	stoppingVer, notifyChn, err := s.handler.Stop(ctx, user, productID, vers.Tag, "testing")
	s.NoError(err)

	// THEN the version status is stopping
	vers.Status = entity.VersionStatusStopping
	s.Equal(vers, stoppingVer)

	// THEN the version status when the go rutine ends is stopped
	versionStatus := <-notifyChn
	s.Equal(entity.VersionStatusStopped, versionStatus.Status)
}
