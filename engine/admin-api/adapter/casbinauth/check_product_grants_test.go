//go:build unit

package casbinauth_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/casbinauth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/auth"
	"github.com/konstellation-io/kai/engine/admin-api/mocks"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckProductGrants(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := mocks.NewMockLogger(ctrl)
	mocks.AddLoggerExpects(logger)

	authorizer, err := casbinauth.NewCasbinAccessControl(logger, casbinModel, casbinPolicy)
	require.NoError(t, err)

	product01 := "product-01"
	product02 := "product-02"

	testCases := []struct {
		name     string
		product  string
		user     *entity.User
		act      auth.Action
		hasError bool
	}{
		{
			name:    "user with grants to view product-01 can view product-01",
			product: product01,
			user: testhelpers.NewUserBuilder().
				WithProductGrants(
					entity.ProductGrants{
						product01: []string{auth.ActViewProduct.String()},
					},
				).
				Build(),
			act:      auth.ActViewProduct,
			hasError: false,
		},
		{
			name:     "user without grants to view product-01 cannot view product-01",
			product:  product01,
			user:     testhelpers.NewUserBuilder().Build(),
			act:      auth.ActViewProduct,
			hasError: true,
		},
		{
			name:    "user with grant to view product-01 but no product-02 cannot view product-02",
			product: product02,
			user: testhelpers.NewUserBuilder().WithProductGrants(
				entity.ProductGrants{
					product01: []string{auth.ActViewProduct.String()},
				},
			).Build(),
			act:      auth.ActViewProduct,
			hasError: true,
		},
		{
			name:    "user with grant to view product-01 cannot create product-01",
			product: product01,
			user: testhelpers.NewUserBuilder().WithProductGrants(
				entity.ProductGrants{
					product01: []string{auth.ActViewProduct.String()},
				},
			).Build(),
			act:      auth.Action(auth.ActCreateProduct.String()),
			hasError: true,
		},
		{
			name:     "admin user can do anything without specifying product grants",
			product:  product01,
			user:     testhelpers.NewUserBuilder().WithRoles([]string{auth.DefaultAdminRole}).Build(),
			act:      auth.Action(auth.ActCreateProduct.String()),
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := authorizer.CheckProductGrants(tc.user, tc.product, tc.act)
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckProductGrants_InvalidAct(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := mocks.NewMockLogger(ctrl)
	mocks.AddLoggerExpects(logger)

	authorizer, err := casbinauth.NewCasbinAccessControl(logger, casbinModel, casbinPolicy)
	require.NoError(t, err)

	user := testhelpers.NewUserBuilder().Build()
	product := "product-01"

	err = authorizer.CheckProductGrants(user, product, "invalid act")
	assert.ErrorIs(t, casbinauth.ErrInvalidAccessControlAction, err)
}
