//go:build unit

package version_test

import (
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
)

func (s *versionSuite) ListVersionsByProduct() {
	// GIVEN a productID and an ID
	productID := "product-1"
	testVersion := testhelpers.NewVersionBuilder().WithID("test-ID").Build()

	s.versionRepo.EXPECT().GetByID(productID, testVersion.ID).Return(testVersion, nil)

	actual, err := s.versionRepo.GetByID(productID, testVersion.ID)
	s.Require().NoError(err)

	s.Equal(testVersion, actual)
}
