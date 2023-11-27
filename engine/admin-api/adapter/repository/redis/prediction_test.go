//go:build integration

package redis_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/repository/redis"
	rdb "github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	_defaultRedisUser = "default"
)

type RedisPredictionRepositorySuite struct {
	suite.Suite
	redisContainer            testcontainers.Container
	redisPredictionRepository *redis.PredictionRepository
	redisClient               *rdb.Client
}

func TestRedisPredictionRepositorySuite(t *testing.T) {
	suite.Run(t, new(RedisPredictionRepositorySuite))
}

func (s *RedisPredictionRepositorySuite) SetupSuite() {
	ctx := context.Background()

	testdataPath, err := filepath.Abs("./testdata")
	s.Require().NoError(err)

	redisConfContainerPath := "/usr/local/etc/redis"

	req := testcontainers.ContainerRequest{
		Image: "redis/redis-stack-server:7.2.0-v6",
		Cmd: []string{
			"redis-stack-server",
			redisConfContainerPath,
		},
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections tcp"),
		Mounts: []testcontainers.ContainerMount{
			{
				Source: testcontainers.DockerBindMountSource{
					HostPath: testdataPath,
				},
				Target: testcontainers.ContainerMountTarget(redisConfContainerPath),
			},
		},
	}

	s.redisContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	s.Require().NoError(err)

	redisEndpoint, err := s.redisContainer.PortEndpoint(ctx, "6379/tcp", "redis")
	s.Require().NoError(err)

	viper.Set(config.RedisEndpointKey, redisEndpoint)
	viper.Set(config.RedisUsernameKey, "default")
	viper.Set(config.RedisPasswordKey, "testpassword")

	s.redisClient, err = redis.NewRedisClient()
	s.Require().NoError(err)

	s.redisPredictionRepository = redis.NewPredictionRepository(s.redisClient)
}

func (s *RedisPredictionRepositorySuite) TearDownSuite() {
	err := s.redisContainer.Terminate(context.Background())
	s.Require().NoError(err)
}

func (s *RedisPredictionRepositorySuite) TearDownTest() {
	ctx := context.Background()
	err := s.redisClient.FlushAll(ctx).Err()
	s.Require().NoError(err)

	res, err := s.redisClient.Do(ctx, "ACL", "USERS").Result()
	s.Require().NoError(err)

	for _, user := range res.([]interface{}) {
		if user != _defaultRedisUser {
			err := s.redisClient.Do(ctx, "ACL", "DELUSER", user).Err()
			s.Require().NoError(err)
		}
	}
}

func (s *RedisPredictionRepositorySuite) TestCreateUser() {
	var (
		ctx      = context.Background()
		testUser = "test-user"
	)

	err := s.redisPredictionRepository.CreateUser(ctx, "test-product", testUser, "test-password")
	s.Require().NoError(err)

	res, err := s.redisClient.Do(ctx, "ACL", "USERS").Result()
	s.Require().NoError(err)

	s.Assert().Contains(res.([]interface{}), testUser)
}

func (s *RedisPredictionRepositorySuite) TestCreateUser_InvalidUsername() {
	ctx := context.Background()

	err := s.redisPredictionRepository.CreateUser(ctx, "test-product", _defaultRedisUser, "test-password")
	s.Require().ErrorIs(err, redis.ErrProtectedUsername)
}

func (s *RedisPredictionRepositorySuite) TestCreateUser_ErrInvalidUsername() {
	ctx := context.Background()

	err := s.redisPredictionRepository.CreateUser(ctx, "test-product", "invalid user", "test-password")
	s.Require().Error(err)
}

func (s *RedisPredictionRepositorySuite) TestEnsureIngressCreated() {
	viper.Set(config.RedisPredictionsIndexKey, "predictionsIdx")

	err := s.redisPredictionRepository.EnsurePredictionIndexCreated()
	s.Require().NoError(err)
}

func (s *RedisPredictionRepositorySuite) TestEnsureIngressCreated_DoestFailIfIndexAlreadyExists() {
	viper.Set(config.RedisPredictionsIndexKey, "predictionsIdx")

	err := s.redisPredictionRepository.EnsurePredictionIndexCreated()
	s.Require().NoError(err)

	err = s.redisPredictionRepository.EnsurePredictionIndexCreated()
	s.Require().NoError(err)
}
