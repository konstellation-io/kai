//go:build integration

package user

import (
	"context"
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	_adminUser     = "admin"
	_adminPassword = "admin"
)

type KeycloakSuite struct {
	suite.Suite
	keycloakContainer    testcontainers.Container
	keycloakUserRegistry *KeycloakUserRegistry
	keycloakClient       *gocloak.GoCloak
}

func TestKeycloakSuite(t *testing.T) {
	suite.Run(t, new(KeycloakSuite))
}

func (s *KeycloakSuite) SetupSuite() {
	ctx := context.Background()

	absFilePath, err := filepath.Abs("./testdata")
	s.Require().NoError(err)

	req := testcontainers.ContainerRequest{
		Image: "quay.io/keycloak/keycloak:latest",
		Cmd: []string{
			"start-dev",
			"--import-realm",
		},
		ExposedPorts: []string{"8080/tcp"},
		WaitingFor:   wait.ForLog("Listening on:"),
		Env: map[string]string{
			"KEYCLOAK_ADMIN":          _adminUser,
			"KEYCLOAK_ADMIN_PASSWORD": _adminPassword,
		},
		Mounts: []testcontainers.ContainerMount{
			{
				Source: testcontainers.DockerBindMountSource{
					HostPath: absFilePath,
				},
				Target: "/opt/keycloak/data/import",
			},
		},
	}

	keycloakContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	s.Require().NoError(err)

	keycloakEndpoint, err := keycloakContainer.Endpoint(ctx, "http")
	s.Require().NoError(err)

	s.keycloakContainer = keycloakContainer

	s.keycloakClient = WithClient(keycloakEndpoint)

	viper.Set(config.KeycloakURLKey, keycloakEndpoint)
	viper.Set(config.KeycloakRealmKey, "example")
	viper.Set(config.KeycloakMasterRealmKey, "master")
	viper.Set(config.KeycloakAdminUserKey, _adminUser)
	viper.Set(config.KeycloakAdminPasswordKey, _adminPassword)
	viper.Set(config.KeycloakAdminClientIDKey, "admin-cli")
}

func (s *KeycloakSuite) TearDownSuite() {
	err := s.keycloakContainer.Terminate(context.Background())
	s.Require().NoError(err)
}

func (s *KeycloakSuite) SetupTest() {
	keycloakUserRegistry, err := NewKeycloakUserRegistry(s.keycloakClient)
	s.Require().NoError(err)

	s.keycloakUserRegistry = keycloakUserRegistry
}

func (s *KeycloakSuite) TearDownTest() {
	testUser := s.getTestUser()
	testUser.Attributes = &map[string][]string{}
	err := s.keycloakClient.UpdateUser(
		context.Background(),
		s.keycloakUserRegistry.token.AccessToken,
		viper.GetString(config.KeycloakRealmKey),
		*testUser,
	)
	s.Require().NoError(err)
}

func (s *KeycloakSuite) getTestUser() *gocloak.User {
	users, err := s.keycloakClient.GetUsers(
		context.Background(),
		s.keycloakUserRegistry.token.AccessToken,
		viper.GetString(config.KeycloakRealmKey),
		gocloak.GetUsersParams{},
	)
	s.Require().NoError(err)

	return users[0]
}

func (s *KeycloakSuite) TestUpdateUserProductGrantsNoPreviousExisting() {
	// GIVEN a user with no previous existing grants and a product
	ctx := context.Background()
	user := s.getTestUser()
	product := "test-product"

	// WHEN updating grants for a product for the first time
	err := s.keycloakUserRegistry.UpdateUserProductGrants(
		ctx,
		*user.ID,
		product,
		[]string{"grant1", "grant2"},
	)
	s.Require().NoError(err)

	// THEN grants for the product are added
	updatedUser := s.getTestUser()
	marshalledAttributes := (*updatedUser.Attributes)["product_roles"]

	s.Require().NotNil(marshalledAttributes)
	s.Require().Len(marshalledAttributes, 1)

	obtainedResult := make(map[string]interface{})
	err = json.Unmarshal([]byte(marshalledAttributes[0]), &obtainedResult)
	s.Require().NoError(err)

	expectedResult := map[string]interface{}{
		product: []interface{}{"grant1", "grant2"},
	}

	s.Equal(expectedResult, obtainedResult)
}

func (s *KeycloakSuite) TestUpdateUserProductGrantsWithPreviousExisting() {
	// GIVEN a user with no previous existing grants and a product
	ctx := context.Background()
	user := s.getTestUser()
	product := "test-product"

	// GIVEN previous existing grants
	err := s.keycloakUserRegistry.UpdateUserProductGrants(
		ctx,
		*user.ID,
		product,
		[]string{"grant1", "grant2"},
	)
	s.Require().NoError(err)

	// WHEN updating grants for a product already existing
	err = s.keycloakUserRegistry.UpdateUserProductGrants(
		ctx,
		*user.ID,
		product,
		[]string{"grant3"},
	)
	s.Require().NoError(err)

	// THEN grants for the product are updated
	updatedUser := s.getTestUser()
	marshalledAttributes := (*updatedUser.Attributes)["product_roles"]

	s.Require().NotNil(marshalledAttributes)
	s.Require().Len(marshalledAttributes, 1)

	obtainedResult := make(map[string]interface{})
	err = json.Unmarshal([]byte(marshalledAttributes[0]), &obtainedResult)
	s.Require().NoError(err)

	expectedResult := map[string]interface{}{
		product: []interface{}{"grant3"},
	}

	s.Equal(expectedResult, obtainedResult)
}

func (s *KeycloakSuite) TestUpdateUserProductGrantsForOtherProduct() {
	// GIVEN a user with no previous existing grants and two products
	ctx := context.Background()
	user := s.getTestUser()
	product := "test-product"
	product2 := "test-product-2"

	// GIVEN previous existing grants
	err := s.keycloakUserRegistry.UpdateUserProductGrants(
		ctx,
		*user.ID,
		product,
		[]string{"grant1", "grant2"},
	)
	s.Require().NoError(err)

	// WHEN adding grants for other product
	err = s.keycloakUserRegistry.UpdateUserProductGrants(
		ctx,
		*user.ID,
		product2,
		[]string{"grant3", "grant4"},
	)
	s.Require().NoError(err)

	// THEN grants for the other product are added
	updatedUser := s.getTestUser()
	marshalledAttributes := (*updatedUser.Attributes)["product_roles"]

	s.Require().NotNil(marshalledAttributes)
	s.Require().Len(marshalledAttributes, 1)

	obtainedResult := make(map[string]interface{})
	err = json.Unmarshal([]byte(marshalledAttributes[0]), &obtainedResult)
	s.Require().NoError(err)

	expectedResult := map[string]interface{}{
		product:  []interface{}{"grant1", "grant2"},
		product2: []interface{}{"grant3", "grant4"},
	}

	s.Equal(expectedResult, obtainedResult)
}

func (s *KeycloakSuite) TestRevokeUserProductGrants() {
	// GIVEN a user with no previous existing grants and a product
	ctx := context.Background()
	user := s.getTestUser()
	product := "test-product"

	// GIVEN previous existing grants
	err := s.keycloakUserRegistry.UpdateUserProductGrants(
		ctx,
		*user.ID,
		product,
		[]string{"grant1", "grant2"},
	)
	s.Require().NoError(err)

	// WHEN revoking grants
	err = s.keycloakUserRegistry.UpdateUserProductGrants(
		ctx,
		*user.ID,
		product,
		[]string{},
	)
	s.Require().NoError(err)

	// THEN grants are revoked
	updatedUser := s.getTestUser()
	marshalledAttributes := (*updatedUser.Attributes)["product_roles"]

	s.Require().NotNil(marshalledAttributes)
	s.Require().Len(marshalledAttributes, 1)

	obtainedResult := make(map[string]interface{})
	err = json.Unmarshal([]byte(marshalledAttributes[0]), &obtainedResult)
	s.Require().NoError(err)

	expectedResult := map[string]interface{}{}

	s.Equal(expectedResult, obtainedResult)
}

func (s *KeycloakSuite) TestRevokeUserProductGrantsForOtherProduct() {
	// GIVEN a user with no previous existing grants and two products
	ctx := context.Background()
	user := s.getTestUser()
	product := "test-product"
	product2 := "test-product-2"

	// GIVEN previous existing grants
	err := s.keycloakUserRegistry.UpdateUserProductGrants(
		ctx,
		*user.ID,
		product,
		[]string{"grant1", "grant2"},
	)
	s.Require().NoError(err)

	err = s.keycloakUserRegistry.UpdateUserProductGrants(
		ctx,
		*user.ID,
		product2,
		[]string{"grant3", "grant4"},
	)
	s.Require().NoError(err)

	// WHEN revoking grants for one product
	err = s.keycloakUserRegistry.UpdateUserProductGrants(
		ctx,
		*user.ID,
		product,
		[]string{},
	)
	s.Require().NoError(err)

	// THEN grants for the other product are not revoked
	updatedUser := s.getTestUser()
	marshalledAttributes := (*updatedUser.Attributes)["product_roles"]

	s.Require().NotNil(marshalledAttributes)
	s.Require().Len(marshalledAttributes, 1)

	obtainedResult := make(map[string]interface{})
	err = json.Unmarshal([]byte(marshalledAttributes[0]), &obtainedResult)
	s.Require().NoError(err)

	expectedResult := map[string]interface{}{
		product2: []interface{}{"grant3", "grant4"},
	}

	s.Equal(expectedResult, obtainedResult)
}

func (s *KeycloakSuite) TestRefreshNotExpiredToken() {
	// GIVEN the recently obtained token through setup test
	ctx := context.Background()
	expiredTimeCopy := s.keycloakUserRegistry.tokenExpiresAt

	// WHEN refreshing the token
	err := s.keycloakUserRegistry.refreshToken(ctx)
	s.Require().NoError(err)

	// THEN the token is not refreshed
	s.True(expiredTimeCopy.Equal(s.keycloakUserRegistry.tokenExpiresAt))
}

func (s *KeycloakSuite) TestRefreshExpiredToken() {
	// GIVEN an expired token
	ctx := context.Background()
	s.keycloakUserRegistry.tokenExpiresAt = time.Now().Add(-time.Hour)

	// WHEN refreshing the token
	now := time.Now()
	err := s.keycloakUserRegistry.refreshToken(ctx)
	s.Require().NoError(err)

	// THEN the token is refreshed
	s.True(now.Before(s.keycloakUserRegistry.tokenExpiresAt))
}

func (s *KeycloakSuite) TestRefreshExpiredRefreshToken() {
	// GIVEN both an expired token and its refresh token expired as well
	ctx := context.Background()
	s.keycloakUserRegistry.tokenExpiresAt = time.Now().Add(-time.Hour)
	s.keycloakUserRegistry.refreshTokenExpiresAt = time.Now().Add(-time.Hour)

	// WHEN refreshing the token
	now := time.Now()
	err := s.keycloakUserRegistry.refreshToken(ctx)
	s.Require().NoError(err)

	// THEN a new token is obtained
	s.True(now.Before(s.keycloakUserRegistry.tokenExpiresAt))
	s.True(now.Before(s.keycloakUserRegistry.refreshTokenExpiresAt))
}

func (s *KeycloakSuite) TestRefreshExpiredTokenWithError() {
	// GIVEN an expired token
	ctx := context.Background()
	s.keycloakUserRegistry.tokenExpiresAt = time.Now().Add(-time.Hour)

	// WHEN refreshing the token with invalid credentials
	s.keycloakUserRegistry.token.RefreshToken = "invalid"
	err := s.keycloakUserRegistry.refreshToken(ctx)

	// THEN an error prompts
	s.Require().Error(err)
}

func (s *KeycloakSuite) TestRefreshExpiredRefreshTokenWithError() {
	// GIVEN both an expired token and its refresh token expired as well
	ctx := context.Background()
	s.keycloakUserRegistry.tokenExpiresAt = time.Now().Add(-time.Hour)
	s.keycloakUserRegistry.refreshTokenExpiresAt = time.Now().Add(-time.Hour)

	// WHEN refreshing the token with invalid credentials
	viper.Set(config.KeycloakAdminPasswordKey, "invalid")
	err := s.keycloakUserRegistry.refreshToken(ctx)

	// THEN an error prompts
	s.Require().Error(err)

	viper.Set(config.KeycloakAdminPasswordKey, _adminPassword)
}
