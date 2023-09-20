//go:build unit

package version_test

import (
	"context"

	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
)

func (s *versionSuite) TestGetByVersion() {
	// GIVEN a productID and a version tag
	productID := "product-1"
	ctx := context.Background()
	testVersion := testhelpers.NewVersionBuilder().WithTag("test-tag").Build()

	s.versionRepo.EXPECT().GetByVersion(ctx, productID, testVersion.Tag).Return(testVersion, nil)

	actual, err := s.versionRepo.GetByVersion(ctx, productID, testVersion.Tag)
	s.Require().NoError(err)

	s.Equal(testVersion, actual)
}
