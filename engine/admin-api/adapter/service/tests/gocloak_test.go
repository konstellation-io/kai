package service_test

import (
	"context"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/Nerzal/gocloak/v13"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/service"
)

type GocloakTestSuite struct {
	suite.Suite
	keycloakContainer testcontainers.Container
	cfg               *service.KeycloakConfig
	gocloakService    *service.GocloakService
	gocloakClient     *gocloak.GoCloak
	gocloakToken      *gocloak.JWT
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

	s.cfg = &service.KeycloakConfig{
		Realm:         "example",
		MasterRealm:   "master",
		URL:           keycloakEndpoint,
		AdminUsername: "admin",
		AdminPassword: "admin",
	}

	gocloakClient := service.WithClient(s.cfg.URL)
	gocloakService, err := service.NewGocloakService(gocloakClient, s.cfg)
	s.Require().NoError(err)

	token, err := gocloakClient.LoginAdmin(
		ctx,
		s.cfg.AdminUsername,
		s.cfg.AdminPassword,
		s.cfg.MasterRealm,
	)
	s.Require().NoError(err)

	s.keycloakContainer = keycloakContainer
	s.gocloakService = gocloakService
	s.gocloakClient = gocloakClient
	s.gocloakToken = token
}

func (s *GocloakTestSuite) TearDownSuite() {
	err := s.keycloakContainer.Terminate(context.Background())
	s.Require().NoError(err)
}

func (s *GocloakTestSuite) getTestUser() *gocloak.User {
	users, err := s.gocloakClient.GetUsers(
		context.Background(),
		s.gocloakToken.AccessToken,
		s.cfg.Realm,
		gocloak.GetUsersParams{},
	)
	s.Require().NoError(err)

	return users[0]
}

func (s *GocloakTestSuite) TestUpdateUserProductGrants() {
	user := s.getTestUser()
	product := "test-product"

	err := s.gocloakService.UpdateUserProductGrants(
		*user.ID,
		product,
		[]string{"grant1", "grant2"},
	)
	s.Require().NoError(err)

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

// TODO: test for revoke grants
// Tests for failing scenarios
// Think about how to maintain keycloak and users clean after each test
