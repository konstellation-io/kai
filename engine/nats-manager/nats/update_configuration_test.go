//go:build integration

package nats_test

import "github.com/nats-io/nats.go"

func (s *ClientTestSuite) TestNatsClient_UpdateConfiguration() {
	var (
		testKeyValueStore = "test-kv-store"
		expectedKey1      = "key1"
		expectedValue1    = "value1"
		expectedKey2      = "key2"
		expectedValue2    = "value2"
	)

	newConfig := map[string]string{
		expectedKey1: expectedValue1,
		expectedKey2: expectedValue2,
	}

	err := s.natsClient.CreateKeyValueStore(testKeyValueStore)
	s.Require().NoError(err)

	err = s.natsClient.UpdateConfiguration(testKeyValueStore, newConfig)
	s.Require().NoError(err)

	kvStoreBucket, err := s.js.KeyValue(testKeyValueStore)
	s.Require().NoError(err)

	entry, err := kvStoreBucket.Get(expectedKey1)
	s.Require().NoError(err)
	s.Assert().Equal(expectedValue1, string(entry.Value()))

	entry, err = kvStoreBucket.Get(expectedKey2)
	s.Require().NoError(err)
	s.Assert().Equal(expectedValue2, string(entry.Value()))
}

func (s *ClientTestSuite) TestNatsClient_UpdateConfiguration_NatsClientError() {
	var (
		testKeyValueStore = "test-kv-store"
		invalidKey        = ".key1"
	)

	newConfig := map[string]string{
		invalidKey: "test value",
	}

	err := s.natsClient.CreateKeyValueStore(testKeyValueStore)
	s.Require().NoError(err)

	err = s.natsClient.UpdateConfiguration(testKeyValueStore, newConfig)
	s.Assert().ErrorIs(err, nats.ErrInvalidKey)
}

func (s *ClientTestSuite) TestNatsClient_UpdateConfiguration_ErrorGettingBucket() {
	nonExistentKeyValueStore := "test-kv-store"

	newConfig := map[string]string{
		"valid-key": "test value",
	}

	err := s.natsClient.UpdateConfiguration(nonExistentKeyValueStore, newConfig)
	s.Assert().ErrorIs(err, nats.ErrBucketNotFound)
}
