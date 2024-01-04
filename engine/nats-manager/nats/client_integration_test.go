//go:build integration

package nats_test

import (
	"context"
	"log"
	"regexp"
	"testing"

	"github.com/go-logr/logr/testr"
	"github.com/konstellation-io/kai/engine/nats-manager/internal"
	natslib "github.com/nats-io/nats.go"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/konstellation-io/kai/engine/nats-manager/internal/entity"
	"github.com/konstellation-io/kai/engine/nats-manager/nats"
)

type ClientTestSuite struct {
	suite.Suite
	natsContainer testcontainers.Container
	natsClient    *nats.NatsClient
	js            natslib.JetStreamContext
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

func (s *ClientTestSuite) SetupSuite() {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "nats:2.8.1",
		Cmd:          []string{"-js"},
		ExposedPorts: []string{"4222/tcp", "8222/tcp"},
		WaitingFor:   wait.ForLog("Server is ready"),
	}

	natsContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal(err)
	}

	natsEndpoint, err := natsContainer.Endpoint(ctx, "")
	if err != nil {
		log.Fatalf("error getting nats container endpoint: %s", err.Error())
	}

	js, err := nats.InitJetStreamConnection(natsEndpoint)
	if err != nil {
		log.Fatalf("error connecting to NATS JetStream: %s", err)
	}

	logger := testr.NewWithOptions(s.T(), testr.Options{Verbosity: -1})
	natsClient := nats.New(logger, js)

	s.js = js
	s.natsContainer = natsContainer
	s.natsClient = natsClient
}

func (s *ClientTestSuite) TearDownSuite() {
	if err := s.natsContainer.Terminate(context.Background()); err != nil {
		log.Fatalf("failed to terminate container: %s", err.Error())
	}
}

func (s *ClientTestSuite) AfterTest(_, _ string) {
	streams := s.js.StreamNames()

	for stream := range streams {
		err := s.js.DeleteStream(stream)
		if err != nil {
			log.Fatalf("error deleting stream %q: %s", stream, err)
		}
	}

	objStores := s.js.ObjectStoreNames()
	for objStore := range objStores {
		err := s.js.DeleteObjectStore(objStore)
		if err != nil {
			log.Fatalf("error deleting object store %q: %s", objStore, err)
		}
	}

	kvStores := s.js.KeyValueStoreNames()
	for kvStore := range kvStores {
		err := s.js.DeleteKeyValue(kvStore)
		if err != nil {
			log.Fatalf("error deleting key-value store %q: %s", kvStore, err)
		}
	}
}

func (s *ClientTestSuite) TestNatsClient_CreateStream() {
	testProcesSubject := "test-stream-test-subject"

	streamConfig := &entity.StreamConfig{
		Stream: "test-stream",
		Processes: entity.ProcessesStreamConfig{
			"test-process": entity.ProcessStreamConfig{
				Subject: testProcesSubject,
			},
		},
	}

	err := s.natsClient.CreateStream(streamConfig)
	s.Require().NoError(err)

	expectedSubjects := []string{
		testProcesSubject,
		testProcesSubject + ".*",
	}

	streams := s.js.Streams()

	amountOfStreams := 0
	for stream := range streams {
		s.Assert().Equal(streamConfig.Stream, stream.Config.Name)
		s.Assert().Equal(expectedSubjects, stream.Config.Subjects)
		amountOfStreams++
	}

	s.Assert().Equal(1, amountOfStreams)
}

func (s *ClientTestSuite) TestNatsClient_CreateStream_ErrorIfStreamAlreadyExists() {
	testStream := "test-stream"
	testProcesSubject := "test-stream-test-subject"

	streamConfig := &entity.StreamConfig{
		Stream: testStream,
		Processes: entity.ProcessesStreamConfig{
			"test-process": entity.ProcessStreamConfig{
				Subject: testProcesSubject,
			},
		},
	}

	_, err := s.js.AddStream(&natslib.StreamConfig{
		Name:      testStream,
		Retention: natslib.InterestPolicy,
	})
	s.Require().NoError(err)

	err = s.natsClient.CreateStream(streamConfig)
	s.Assert().ErrorIs(err, natslib.ErrStreamNameAlreadyInUse)
}

func (s *ClientTestSuite) TestNatsClient_DeleteStream() {
	testStream := "test-stream"

	_, err := s.js.AddStream(&natslib.StreamConfig{
		Name:      testStream,
		Retention: natslib.InterestPolicy,
	})
	s.Require().NoError(err)

	err = s.natsClient.DeleteStream(testStream)
	s.Assert().NoError(err)

	streams := s.js.Streams()

	s.Assert().Nil(<-streams)
}

func (s *ClientTestSuite) TestNatsClient_DeleteStream_ErrorIfStreamDoesntExist() {
	testStream := "test-stream"

	err := s.natsClient.DeleteStream(testStream)
	s.Assert().Error(err, natslib.ErrStreamNotFound)
}

func (s *ClientTestSuite) TestNatsClient_CreateObjectStore() {
	testObjectStore := "test-objstore"

	err := s.natsClient.CreateObjectStore(testObjectStore)
	s.Assert().NoError(err)

	objectStores := s.js.ObjectStores()

	objectStore := <-objectStores
	s.Assert().Equal(testObjectStore, objectStore.Bucket())

	s.Assert().Equal(nil, <-objectStores)
}

func (s *ClientTestSuite) TestNatsClient_CreateObjectStore_NoErrorWhenObjectStoreAlreadyExistsWithSameConfig() {
	testObjectStore := "test-objstore"

	_, err := s.js.CreateObjectStore(&natslib.ObjectStoreConfig{
		Bucket:  testObjectStore,
		Storage: natslib.FileStorage,
	})
	s.Require().NoError(err)

	err = s.natsClient.CreateObjectStore(testObjectStore)
	s.Assert().NoError(err)
}

func (s *ClientTestSuite) TestNatsClient_CreateObjectStore_ErrorWhenObjectStoreAlreadyExistsWithDiffConfig() {
	testObjectStore := "test-objstore"

	_, err := s.js.CreateObjectStore(&natslib.ObjectStoreConfig{
		Bucket:  testObjectStore,
		Storage: natslib.MemoryStorage,
	})
	s.Require().NoError(err)

	err = s.natsClient.CreateObjectStore(testObjectStore)
	s.Assert().ErrorIs(err, natslib.ErrStreamNameAlreadyInUse)
}

func (s *ClientTestSuite) TestNatsClient_DeleteObjectStore() {
	testObjectStore := "test-objstore"

	_, err := s.js.CreateObjectStore(&natslib.ObjectStoreConfig{
		Bucket:  testObjectStore,
		Storage: natslib.FileStorage,
	})
	s.Require().NoError(err)

	err = s.natsClient.DeleteObjectStore(testObjectStore)
	s.Assert().NoError(err)

	objectStores := s.js.ObjectStores()
	s.Assert().Nil(<-objectStores)
}

func (s *ClientTestSuite) TestNatsClient_DeleteObjectStore_ObjectStoreNotFound() {
	testObjectStore := "test-objstore"

	err := s.natsClient.DeleteObjectStore(testObjectStore)
	s.Assert().ErrorIs(err, natslib.ErrStreamNotFound)
}

func (s *ClientTestSuite) TestNatsClient_GetObjectStoresNames() {
	testObjectStore1 := "test-object-store-1"
	testObjectStore2 := "test-object-store-2"
	objectStoreWithOtherFormat := "another-obj-store"

	testCases := []struct {
		name                 string
		optFilter            []*regexp.Regexp
		existingObjectStores []string
		expectedObjectStores []string
		wantError            bool
		expectedError        error
	}{
		{
			name: "Get object store names without filter",
			existingObjectStores: []string{
				testObjectStore1,
				testObjectStore2,
				objectStoreWithOtherFormat,
			},
			expectedObjectStores: []string{
				testObjectStore1,
				testObjectStore2,
			},
			optFilter: []*regexp.Regexp{regexp.MustCompile("test-object-.*")},
			wantError: false,
		},
		{
			name: "Get object store names with regex filter",
			existingObjectStores: []string{
				testObjectStore1,
				testObjectStore2,
				objectStoreWithOtherFormat,
			},
			expectedObjectStores: []string{
				testObjectStore1,
				testObjectStore2,
				objectStoreWithOtherFormat,
			},
			optFilter: nil,
			wantError: false,
		},
		{
			name: "Get object store names with regex filter",
			existingObjectStores: []string{
				testObjectStore1,
				testObjectStore2,
				objectStoreWithOtherFormat,
			},
			expectedObjectStores: nil,
			optFilter:            []*regexp.Regexp{regexp.MustCompile(""), regexp.MustCompile("")},
			wantError:            true,
			expectedError:        internal.ErrNoOptFilter,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			for _, objStore := range tc.existingObjectStores {
				_, err := s.js.CreateObjectStore(&natslib.ObjectStoreConfig{
					Bucket:  objStore,
					Storage: natslib.FileStorage,
				})
				s.Require().NoError(err)
			}

			var objectStores []string
			var err error
			if tc.optFilter == nil {
				objectStores, err = s.natsClient.GetObjectStoreNames()
			} else {
				objectStores, err = s.natsClient.GetObjectStoreNames(tc.optFilter...)
			}

			if tc.wantError {
				s.Assert().ErrorIs(err, tc.expectedError)
				return
			}

			s.Assert().NoError(err)
			s.Assert().ElementsMatch(tc.expectedObjectStores, objectStores)
		})
	}
}

func (s *ClientTestSuite) TestNatsClient_GetStreamNames() {
	testStream1 := "test-stream-1"
	testStream2 := "test-stream-2"
	streamWithOtherFormat := "another-stream"

	testCases := []struct {
		name            string
		optFilter       []*regexp.Regexp
		existingStreams []string
		expectedStreams []string
		wantError       bool
		expectedError   error
	}{
		{
			name: "Get stream names without filter",
			existingStreams: []string{
				testStream1,
				testStream2,
				streamWithOtherFormat,
			},
			expectedStreams: []string{
				testStream1,
				testStream2,
			},
			optFilter: []*regexp.Regexp{regexp.MustCompile("test-stream-.*")},
			wantError: false,
		},
		{
			name: "Get stream names with regex filter",
			existingStreams: []string{
				testStream1,
				testStream2,
				streamWithOtherFormat,
			},
			expectedStreams: []string{
				testStream1,
				testStream2,
				streamWithOtherFormat,
			},
			optFilter: nil,
			wantError: false,
		},
		{
			name: "Get stream names with regex filter",
			existingStreams: []string{
				testStream1,
				testStream2,
				streamWithOtherFormat,
			},
			expectedStreams: nil,
			optFilter:       []*regexp.Regexp{regexp.MustCompile(""), regexp.MustCompile("")},
			wantError:       true,
			expectedError:   internal.ErrNoOptFilter,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			for _, streamName := range tc.existingStreams {
				_, err := s.js.AddStream(&natslib.StreamConfig{
					Name: streamName,
				})
				s.Require().NoError(err)
			}

			var streams []string
			var err error
			if tc.optFilter == nil {
				streams, err = s.natsClient.GetStreamNames()
			} else {
				streams, err = s.natsClient.GetStreamNames(tc.optFilter...)
			}

			if tc.wantError {
				s.Assert().ErrorIs(err, tc.expectedError)
				return
			}

			s.Assert().NoError(err)
			s.Assert().ElementsMatch(tc.expectedStreams, streams)
		})
	}
}

func (s *ClientTestSuite) TestNatsClient_GetStreamNames_DoesntReturnObjectStores() {
	testStreamName := "product-id_version-id_workflows-id"

	_, err := s.js.CreateObjectStore(&natslib.ObjectStoreConfig{
		Bucket:  testStreamName,
		Storage: natslib.MemoryStorage,
	})
	s.Require().NoError(err)

	_, err = s.js.AddStream(&natslib.StreamConfig{
		Name:     testStreamName,
		Subjects: []string{testStreamName + ".*"},
	})
	s.Require().NoError(err)

	actualStreams, err := s.natsClient.GetStreamNames(regexp.MustCompile("^product-id_version-id_.*"))
	s.Assert().NoError(err)
	expectedStreams := []string{testStreamName}

	s.Assert().ElementsMatch(expectedStreams, actualStreams)
}

func (s *ClientTestSuite) TestNatsClient_CreateCreateKeyValueStore() {
	testKeyValueStore := "test-kv-store"

	err := s.natsClient.CreateKeyValueStore(testKeyValueStore)
	s.Assert().NoError(err)

	keyValueStores := s.js.KeyValueStores()

	keyValueStore := <-keyValueStores
	s.Assert().Equal(testKeyValueStore, keyValueStore.Bucket())

	s.Assert().Equal(nil, <-keyValueStores)
}

func (s *ClientTestSuite) TestNatsClient_CreateCreateKeyValueStore_NoErrorWhenKVStoreHaveSameNameAndConfig() {
	testKeyValueStore := "test-kv-store"

	_, err := s.js.CreateKeyValue(&natslib.KeyValueConfig{
		Bucket: testKeyValueStore,
	})
	s.Require().NoError(err)

	err = s.natsClient.CreateKeyValueStore(testKeyValueStore)
	s.Assert().NoError(err)
}

func (s *ClientTestSuite) TestNatsClient_CreateCreateKeyValueStore_ErrorWhenDuplicatedNameHasDiffConfig() {
	testKeyValueStore := "test-kv-store"

	_, err := s.js.CreateKeyValue(&natslib.KeyValueConfig{
		Bucket:  testKeyValueStore,
		Storage: natslib.MemoryStorage,
	})
	s.Require().NoError(err)

	err = s.natsClient.CreateKeyValueStore(testKeyValueStore)
	s.Assert().ErrorIs(err, natslib.ErrStreamNameAlreadyInUse)
}

func (s *ClientTestSuite) TestNatsClient_DeleteKeyValueStore() {
	testKeyValueStore := "test-kv-store"

	_, err := s.js.CreateKeyValue(&natslib.KeyValueConfig{
		Bucket: testKeyValueStore,
	})
	s.Require().NoError(err)

	err = s.natsClient.DeleteKeyValueStore(testKeyValueStore)
	s.Assert().NoError(err)

	keyValueStores := s.js.KeyValueStores()
	s.Assert().Nil(<-keyValueStores)
}

func (s *ClientTestSuite) TestNatsClient_DeleteKeyValueStore_KeyValueStoreNotFound() {
	testKeyValueStore := "test-kv-store"

	err := s.natsClient.DeleteKeyValueStore(testKeyValueStore)
	s.Assert().ErrorIs(err, natslib.ErrStreamNotFound)
}
