//go:build unit

package version_test

import (
	"context"

	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
)

func (s *versionSuite) TestGetLatest() {
	// GIVEN a productID and a version tag
	productID := "product-1"
	ctx := context.Background()
	testVersion := testhelpers.NewVersionBuilder().WithTag("test-tag").Build()

	s.versionRepo.EXPECT().GetLatest(ctx, productID).Return(testVersion, nil)

	actual, err := s.versionRepo.GetLatest(ctx, productID)
	s.Require().NoError(err)

	s.Equal(testVersion, actual)
}
