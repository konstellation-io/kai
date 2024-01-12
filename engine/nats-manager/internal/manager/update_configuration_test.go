//go:build unit

package manager_test

import (
	"errors"
	"testing"

	"github.com/go-logr/logr/testr"
	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/nats-manager/internal/entity"
	"github.com/konstellation-io/kai/engine/nats-manager/internal/manager"
	"github.com/konstellation-io/kai/engine/nats-manager/mocks"
	"github.com/stretchr/testify/suite"
)

type UpdateConfigurationSuite struct {
	suite.Suite

	client      *mocks.MockNatsClient
	natsManager *manager.NatsManager
}

func TestUpdateConfigurationSuite(t *testing.T) {
	suite.Run(t, new(UpdateConfigurationSuite))
}

func (s *UpdateConfigurationSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())

	logger := testr.NewWithOptions(s.T(), testr.Options{Verbosity: -1})
	s.client = mocks.NewMockNatsClient(ctrl)
	s.natsManager = manager.NewNatsManager(logger, s.client)
}

func (s *UpdateConfigurationSuite) TestUpdateKeyValueStoresConfiguration() {
	var (
		expectedConfiguration = entity.KeyValueConfiguration{
			KeyValueStore: "version-kv-store",
			Configuration: map[string]string{
				"key1": "val1",
			},
		}

		configurations = []entity.KeyValueConfiguration{
			expectedConfiguration,
		}
	)

	s.client.EXPECT().UpdateConfiguration(expectedConfiguration.KeyValueStore, expectedConfiguration.Configuration).
		Return(nil)

	err := s.natsManager.UpdateKeyValueStoresConfiguration(configurations)
	s.Assert().NoError(err)
}

func (s *UpdateConfigurationSuite) TestUpdateKeyValueStoresConfiguration_NoConfigurations() {
	err := s.natsManager.UpdateKeyValueStoresConfiguration(nil)
	s.Assert().NoError(err)
}

func (s *UpdateConfigurationSuite) TestUpdateKeyValueStoresConfiguration_NatsClientError() {
	var (
		expectedConfiguration = entity.KeyValueConfiguration{
			KeyValueStore: "version-kv-store",
			Configuration: map[string]string{
				"key1": "val1",
			},
		}

		configurations = []entity.KeyValueConfiguration{
			expectedConfiguration,
		}

		expectedError = errors.New("nats client error")
	)

	s.client.EXPECT().UpdateConfiguration(expectedConfiguration.KeyValueStore, expectedConfiguration.Configuration).
		Return(expectedError)

	err := s.natsManager.UpdateKeyValueStoresConfiguration(configurations)
	s.Assert().ErrorIs(err, expectedError)
}
