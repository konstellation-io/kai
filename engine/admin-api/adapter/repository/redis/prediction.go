package redis

import (
	"context"
	"fmt"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

const (
	_defaultRedisUser = "default"
)

var (
	ErrProtectedUsername = fmt.Errorf("username %q is protected", _defaultRedisUser)
)

type PredictionRepository struct {
	client *redis.Client
}

func NewRedisClient() (*redis.Client, error) {
	opts, err := redis.ParseURL(viper.GetString(config.RedisEndpointKey))
	if err != nil {
		return nil, fmt.Errorf("parsing Redis URL: %w", err)
	}

	opts.Username = "default"
	opts.Password = viper.GetString(config.RedisPasswordKey)

	return redis.NewClient(opts), nil
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
		allowedKeys = fmt.Sprintf("~%s:*", product)
		passwordArg = fmt.Sprintf(">%s", password)
	)

	command := r.client.Do(ctx, "ACL", "SETUSER", username, allowedKeys, passwordArg, "+@all", "-@dangerous", "on")

	if err := command.Err(); err != nil {
		return err
	}

	return nil
}
