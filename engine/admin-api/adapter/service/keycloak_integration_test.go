//go:build integration

// This test file needs to be in service package, otherwise it won't be able to access the private
// fields of the service package. As some tests alter manually GocloakUserRegistry struct variables.
package service

import (
	"context"
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type GocloakTestSuite struct {
	suite.Suite
	keycloakContainer   testcontainers.Container
	cfg                 *KeycloakConfig
	gocloakUserRegistry *GocloakUserRegistry
	gocloakClient       *gocloak.GoCloak
}

func TestGocloakTestSuite(t *testing.T) {
	suite.Run(t, new(GocloakTestSuite))
}

func (s *GocloakTestSuite) SetupSuite() {
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
			"KEYCLOAK_ADMIN":          "admin",
			"KEYCLOAK_ADMIN_PASSWORD": "admin",
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

	s.cfg = &KeycloakConfig{
		Realm:         "example",
		MasterRealm:   "master",
		AdminUsername: "admin",
		AdminPassword: "admin",
		AdminClientID: "admin-cli",
	}

	s.keycloakContainer = keycloakContainer
	s.gocloakClient = WithClient(keycloakEndpoint)
}

func (s *GocloakTestSuite) TearDownSuite() {
	err := s.keycloakContainer.Terminate(context.Background())
	s.Require().NoError(err)
}

func (s *GocloakTestSuite) SetupTest() {
	gocloakUserRegistry, err := NewGocloakUserRegistry(s.gocloakClient, s.cfg)
	s.Require().NoError(err)

	s.gocloakUserRegistry = gocloakUserRegistry
}

func (s *GocloakTestSuite) TearDownTest() {
	testUser := s.getTestUser()
	testUser.Attributes = &map[string][]string{}
	err := s.gocloakClient.UpdateUser(
		context.Background(),
		s.gocloakUserRegistry.token.AccessToken,
		s.cfg.Realm,
		*testUser,
	)
	s.Require().NoError(err)
}

func (s *GocloakTestSuite) getTestUser() *gocloak.User {
	users, err := s.gocloakClient.GetUsers(
		context.Background(),
		s.gocloakUserRegistry.token.AccessToken,
		s.cfg.Realm,
		gocloak.GetUsersParams{},
	)
	s.Require().NoError(err)

	return users[0]
}

func (s *GocloakTestSuite) TestUpdateUserProductGrantsNoPreviousExisting() {
	// GIVEN a user with no previous existing grants and a product
	user := s.getTestUser()
	product := "test-product"

	// WHEN updating grants for a product for the first time
	err := s.gocloakUserRegistry.UpdateUserProductGrants(
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

func (s *GocloakTestSuite) TestUpdateUserProductGrantsWithPreviousExisting() {
	// GIVEN a user with no previous existing grants and a product
	user := s.getTestUser()
	product := "test-product"

	// GIVEN previous existing grants
	err := s.gocloakUserRegistry.UpdateUserProductGrants(
		*user.ID,
		product,
		[]string{"grant1", "grant2"},
	)
	s.Require().NoError(err)

	// WHEN updating grants for a product already existing
	err = s.gocloakUserRegistry.UpdateUserProductGrants(
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

func (s *GocloakTestSuite) TestUpdateUserProductGrantsForOtherProduct() {
	// GIVEN a user with no previous existing grants and two products
	user := s.getTestUser()
	product := "test-product"
	product2 := "test-product-2"

	// GIVEN previous existing grants
	err := s.gocloakUserRegistry.UpdateUserProductGrants(
		*user.ID,
		product,
		[]string{"grant1", "grant2"},
	)
	s.Require().NoError(err)

	// WHEN adding grants for other product
	err = s.gocloakUserRegistry.UpdateUserProductGrants(
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

func (s *GocloakTestSuite) TestRevokeUserProductGrants() {
	// GIVEN a user with no previous existing grants and a product
	user := s.getTestUser()
	product := "test-product"

	// GIVEN previous existing grants
	err := s.gocloakUserRegistry.UpdateUserProductGrants(
		*user.ID,
		product,
		[]string{"grant1", "grant2"},
	)
	s.Require().NoError(err)

	// WHEN revoking grants
	err = s.gocloakUserRegistry.UpdateUserProductGrants(
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

func (s *GocloakTestSuite) TestRevokeUserProductGrantsForOtherProduct() {
	// GIVEN a user with no previous existing grants and two products
	user := s.getTestUser()
	product := "test-product"
	product2 := "test-product-2"

	// GIVEN previous existing grants
	err := s.gocloakUserRegistry.UpdateUserProductGrants(
		*user.ID,
		product,
		[]string{"grant1", "grant2"},
	)
	s.Require().NoError(err)

	err = s.gocloakUserRegistry.UpdateUserProductGrants(
		*user.ID,
		product2,
		[]string{"grant3", "grant4"},
	)
	s.Require().NoError(err)

	// WHEN revoking grants for one product
	err = s.gocloakUserRegistry.UpdateUserProductGrants(
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

func (s *GocloakTestSuite) TestRefreshNotExpiredToken() {
	// GIVEN the recently obtained token through setup test
	expiredTimeCopy := s.gocloakUserRegistry.tokenExpiresAt

	// WHEN refreshing the token
	err := s.gocloakUserRegistry.refreshToken()
	s.Require().NoError(err)

	// THEN the token is not refreshed
	s.True(expiredTimeCopy.Equal(s.gocloakUserRegistry.tokenExpiresAt))
}

func (s *GocloakTestSuite) TestRefreshExpiredToken() {
	// GIVEN an expired token
	s.gocloakUserRegistry.tokenExpiresAt = time.Now().Add(-time.Hour)

	// WHEN refreshing the token
	now := time.Now()
	err := s.gocloakUserRegistry.refreshToken()
	s.Require().NoError(err)

	// THEN the token is refreshed
	s.True(now.Before(s.gocloakUserRegistry.tokenExpiresAt))
}

func (s *GocloakTestSuite) TestRefreshExpiredRefreshToken() {
	// GIVEN both an expired token and its refresh token expired as well
	s.gocloakUserRegistry.tokenExpiresAt = time.Now().Add(-time.Hour)
	s.gocloakUserRegistry.refreshtokenExpiresAt = time.Now().Add(-time.Hour)

	// WHEN refreshing the token
	now := time.Now()
	err := s.gocloakUserRegistry.refreshToken()
	s.Require().NoError(err)

	// THEN a new token is obtained
	s.True(now.Before(s.gocloakUserRegistry.tokenExpiresAt))
	s.True(now.Before(s.gocloakUserRegistry.refreshtokenExpiresAt))
}
