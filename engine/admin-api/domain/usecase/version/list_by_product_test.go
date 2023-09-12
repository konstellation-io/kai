//go:build unit

package version_test

import (
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/version/utils"
)

func (s *VersionUsecaseTestSuite) ListVersionsByProduct() {
	// GIVEN a productID and an ID
	productID := "product-1"
	testVersion := utils.InitTestVersion().WithVersionID("test-ID").GetVersion()

	s.versionRepo.EXPECT().GetByID(productID, testVersion.ID).Return(testVersion, nil)

	actual, err := s.versionRepo.GetByID(productID, testVersion.ID)
	s.Require().NoError(err)

	s.Equal(testVersion, actual)
}
