//go:build integration

package user

import (
	"context"
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
	_testUsername  = "user"
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

	keycloakEndpoint, err := keycloakContainer.PortEndpoint(ctx, "8080/tcp", "http")
	s.Require().NoError(err)

	s.keycloakContainer = keycloakContainer

	s.keycloakClient = WithClient(keycloakEndpoint)

	viper.Set(config.KeycloakURLKey, keycloakEndpoint)
	viper.Set(config.KeycloakRealmKey, "example")
	viper.Set(config.KeycloakMasterRealmKey, "master")
	viper.Set(config.KeycloakAdminUserKey, _adminUser)
	viper.Set(config.KeycloakAdminPasswordKey, _adminPassword)
	viper.Set(config.KeycloakAdminClientIDKey, "admin-cli")
	viper.Set(config.KeycloakPolicyAttributeKey, "policy")
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
	ctx := context.Background()

	testUser := s.getTestUser()
	testUser.Attributes = &map[string][]string{}

	err := s.keycloakClient.UpdateUser(
		ctx,
		s.keycloakUserRegistry.token.AccessToken,
		viper.GetString(config.KeycloakRealmKey),
		*testUser,
	)
	s.Require().NoError(err)

	groups, err := s.keycloakClient.GetGroups(
		ctx,
		s.keycloakUserRegistry.token.AccessToken,
		viper.GetString(config.KeycloakRealmKey),
		gocloak.GetGroupsParams{},
	)
	s.Require().NoError(err)

	for _, group := range groups {
		s.keycloakClient.DeleteGroup(
			ctx,
			s.keycloakUserRegistry.token.AccessToken,
			viper.GetString(config.KeycloakRealmKey),
			*group.ID,
		)
		s.Require().NoError(err)
	}

	users, err := s.keycloakClient.GetUsers(
		ctx,
		s.keycloakUserRegistry.token.AccessToken,
		viper.GetString(config.KeycloakRealmKey),
		gocloak.GetUsersParams{},
	)

	s.Require().NoError(err)

	for _, user := range users {
		if *user.Username != _testUsername {
			err = s.keycloakClient.DeleteUser(
				ctx,
				s.keycloakUserRegistry.token.AccessToken,
				viper.GetString(config.KeycloakRealmKey),
				*user.ID,
			)
			s.Require().NoError(err)
		}
	}
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

func (s *KeycloakSuite) TestRefreshToken_NotExpiredToken() {
	// GIVEN the recently obtained token through setup test
	ctx := context.Background()
	expiredTimeCopy := s.keycloakUserRegistry.tokenExpiresAt

	// WHEN refreshing the token
	err := s.keycloakUserRegistry.refreshToken(ctx)
	s.Require().NoError(err)

	// THEN the token is not refreshed
	s.True(expiredTimeCopy.Equal(s.keycloakUserRegistry.tokenExpiresAt))
}

func (s *KeycloakSuite) TestRefreshToken_ExpiredToken() {
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

func (s *KeycloakSuite) TestRefreshToken_ExpiredRefreshToken() {
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

func (s *KeycloakSuite) TestRefreshToken_ExpiredTokenWithError() {
	// GIVEN an expired token
	ctx := context.Background()
	s.keycloakUserRegistry.tokenExpiresAt = time.Now().Add(-time.Hour)

	// WHEN refreshing the token with invalid credentials
	s.keycloakUserRegistry.token.RefreshToken = "invalid"
	err := s.keycloakUserRegistry.refreshToken(ctx)

	// THEN an error prompts
	s.Require().Error(err)
}

func (s *KeycloakSuite) TestRefreshToken_ExpiredRefreshTokenWithError() {
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
