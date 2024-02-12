package redis

import (
	"context"
	"fmt"
	"strings"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

const (
	_defaultRedisUser   = "default"
	_indexAlreadyExists = "Index already exists"
)

var (
	ErrProtectedUsername = fmt.Errorf("username %q is protected", _defaultRedisUser)
)

type PredictionRepository struct {
	client *redis.Client
}

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     viper.GetString(config.RedisEndpointKey),
		Username: viper.GetString(config.RedisUsernameKey),
		Password: viper.GetString(config.RedisPasswordKey),
	})
}

func NewPredictionRepository(client *redis.Client) *PredictionRepository {
	return &PredictionRepository{
		client: client,
	}
}

func (r *PredictionRepository) CreateUser(ctx context.Context, product, username, password string) error {
	if username == _defaultRedisUser {
		return ErrProtectedUsername
	}

	var (
		allowedKeys       = fmt.Sprintf("~%s:*", product)
		passwordArg       = fmt.Sprintf(">%s", password)
		predictionsIdxArg = fmt.Sprintf("~%s", viper.GetString(config.RedisPredictionsIndexKey))
	)

	command := r.client.Do(ctx, "ACL", "SETUSER", username, allowedKeys, passwordArg, "+@all", "-@dangerous", "on",
		"+ft.search", predictionsIdxArg)

	if err := command.Err(); err != nil {
		return err
	}

	return nil
}

func (r *PredictionRepository) DeleteUser(ctx context.Context, username string) error {
	err := r.client.Do(ctx, "ACL", "DELUSER", username).Err()
	if err != nil {
		return fmt.Errorf("deleting user %q from Redis: %w", username, err)
	}

	return nil
}

func (r *PredictionRepository) EnsurePredictionIndexCreated() error {
	command := r.client.Do(context.Background(),
		"FT.CREATE",
		viper.GetString(config.RedisPredictionsIndexKey),
		"ON",
		"JSON",
		"SCHEMA",
		"$.metadata.product", "AS", "product", "TAG",
		"$.metadata.version", "AS", "version", "TAG",
		"$.metadata.workflow", "AS", "workflow", "TAG",
		"$.metadata.process", "AS", "process", "TAG",
		"$.metadata.workflow_type", "AS", "workflow_type", "TAG",
		"$.metadata.request_id", "AS", "request_id", "TAG",
		"$.creation_date", "AS", "creation_date", "NUMERIC",
	)
	if err := command.Err(); err != nil {
		if r.isIndexAlreadyExistsError(err) {
			return nil
		}

		return fmt.Errorf("creating predictions index: %w", err)
	}

	return nil
}

func (r *PredictionRepository) isIndexAlreadyExistsError(err error) bool {
	return strings.Contains(err.Error(), "Index already exists")
}
