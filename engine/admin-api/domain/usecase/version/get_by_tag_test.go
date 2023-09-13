//go:build unit

package version_test

import (
	"context"

	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/version/utils"
)

func (s *VersionUsecaseTestSuite) TestGetByTag() {
	// GIVEN a productID and a version tag
	productID := "product-1"
	ctx := context.Background()
	testVersion := utils.InitTestVersion().WithTag("test-tag").GetVersion()

	s.versionRepo.EXPECT().GetByTag(ctx, productID, testVersion.Tag).Return(testVersion, nil)

	actual, err := s.versionRepo.GetByTag(ctx, productID, testVersion.Tag)
	s.Require().NoError(err)

	s.Equal(testVersion, actual)
}
