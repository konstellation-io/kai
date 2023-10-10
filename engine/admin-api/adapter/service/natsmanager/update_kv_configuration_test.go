//go:build unit

package natsmanager_test

import (
	"context"
	"errors"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/natspb"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

func (s *NatsManagerTestSuite) TestUpdateKeyValueStoresConfiguration() {
	var (
		ctx           = context.Background()
		kvStore       = "keyValueStore"
		configuration = entity.ConfigurationVariable{
			Key: "key1", Value: "val1",
		}

		configurations = []entity.KeyValueConfiguration{{
			Store:         kvStore,
			Configuration: []entity.ConfigurationVariable{configuration},
		}}
	)

	expectedRequest := &natspb.UpdateKeyValueConfigurationRequest{
		KeyValueStoresConfig: []*natspb.KeyValueConfiguration{
			{
				KeyValueStore: kvStore,
				Configuration: map[string]string{
					configuration.Key: configuration.Value,
				},
			},
		},
	}

	s.mockService.EXPECT().UpdateKeyValueConfiguration(ctx, expectedRequest).
		Return(&natspb.UpdateKeyValueConfigurationResponse{}, nil)

	err := s.natsManagerClient.UpdateKeyValueConfiguration(ctx, configurations)
	s.Assert().NoError(err)
}

func (s *NatsManagerTestSuite) TestUpdateKeyValueStoresConfiguration_ServiceError() {
	var (
		ctx           = context.Background()
		expectedError = errors.New("nats service error")
	)

	s.mockService.EXPECT().UpdateKeyValueConfiguration(ctx, gomock.Any()).
		Return(nil, expectedError)

	err := s.natsManagerClient.UpdateKeyValueConfiguration(ctx, []entity.KeyValueConfiguration{})
	s.Assert().ErrorIs(err, expectedError)
}
