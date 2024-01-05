package nats

import (
	"fmt"
	"regexp"
	"time"

	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/nats-manager/internal/config"
	"github.com/spf13/viper"

	"github.com/konstellation-io/kai/engine/nats-manager/internal"
	"github.com/konstellation-io/kai/engine/nats-manager/internal/entity"
	"github.com/nats-io/nats.go"
)

type NatsClient struct {
	js     nats.JetStreamContext
	logger logr.Logger
}

func New(logger logr.Logger, js nats.JetStreamContext) *NatsClient {
	return &NatsClient{
		logger: logger,
		js:     js,
	}
}

func InitJetStreamConnection(url string) (nats.JetStreamContext, error) {
	natsConn, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("error connecting to NATS: %w", err)
	}

	js, err := natsConn.JetStream()
	if err != nil {
		return nil, fmt.Errorf("error connecting to NATS JetStream: %w", err)
	}

	return js, nil
}

func (n *NatsClient) CreateStream(streamConfig *entity.StreamConfig) error {
	n.logger.Info("Creating stream", "stream", streamConfig.Stream)

	subjects := n.getProcessesSubjects(streamConfig.Processes)

	streamCfg := &nats.StreamConfig{
		Name:        streamConfig.Stream,
		Description: "",
		Subjects:    subjects,
		Retention:   nats.InterestPolicy,
	}

	_, err := n.js.AddStream(streamCfg)

	return err
}

// GetObjectStoreNames returns the list of object stores.
// The optional param `optFilter` accepts 0 or 1 value.
func (n *NatsClient) GetObjectStoreNames(optFilter ...*regexp.Regexp) ([]string, error) {
	if len(optFilter) > 1 {
		return nil, internal.ErrNoOptFilter
	}

	var regexpFilter *regexp.Regexp
	if len(optFilter) == 1 {
		regexpFilter = optFilter[0]
	}

	objectStoresCh := n.js.ObjectStores()
	objectStores := make([]string, 0)

	for objectStore := range objectStoresCh {
		objStoreName := objectStore.Bucket()

		nameMatchFilter := regexpFilter == nil || regexpFilter.MatchString(objStoreName)
		if nameMatchFilter {
			objectStores = append(objectStores, objStoreName)
		}
	}

	return objectStores, nil
}

func (n *NatsClient) CreateObjectStore(objectStore string) error {
	n.logger.Info("Creating object store", "object-store", objectStore)

	_, err := n.js.CreateObjectStore(&nats.ObjectStoreConfig{
		Bucket:  objectStore,
		Storage: nats.FileStorage,
		TTL:     time.Duration(viper.GetInt(config.ObjectStoreDefaultTTLDays)*24) * time.Hour,
	})
	if err != nil {
		return fmt.Errorf("error creating the object store: %w", err)
	}

	return nil
}

func (n *NatsClient) CreateKeyValueStore(keyValueStore string) error {
	n.logger.Info("Creating key-value store", "key-value-store", keyValueStore)

	_, err := n.js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: keyValueStore,
	})
	if err != nil {
		return fmt.Errorf("error creating the key-value store: %w", err)
	}

	return nil
}

func (n *NatsClient) DeleteStream(stream string) error {
	n.logger.Info("Deleting stream", "stream", stream)
	return n.js.DeleteStream(stream)
}

func (n *NatsClient) DeleteObjectStore(objectStore string) error {
	n.logger.Info("Deleting object store", "object-store", objectStore)
	return n.js.DeleteObjectStore(objectStore)
}

func (n *NatsClient) DeleteKeyValueStore(keyValueStore string) error {
	n.logger.Info("Deleting key-value store", "key-value-store", keyValueStore)
	return n.js.DeleteKeyValue(keyValueStore)
}

// GetStreamNames returns the list of streams' names.
// The optional param `optFilter` accepts 0 or 1 value.
func (n *NatsClient) GetStreamNames(optFilter ...*regexp.Regexp) ([]string, error) {
	if len(optFilter) > 1 {
		return nil, internal.ErrNoOptFilter
	}

	var regexpFilter *regexp.Regexp
	if len(optFilter) == 1 {
		regexpFilter = optFilter[0]
	}

	streamsChannel := n.js.StreamNames()
	streams := make([]string, 0)

	for streamName := range streamsChannel {
		nameMatchFilter := regexpFilter == nil || regexpFilter.MatchString(streamName)
		if nameMatchFilter {
			streams = append(streams, streamName)
		}
	}

	return streams, nil
}

func (n *NatsClient) getProcessesSubjects(processes entity.ProcessesStreamConfig) []string {
	subjects := make([]string, 0, len(processes)*2)

	for _, processCfg := range processes {
		subSubject := processCfg.Subject + ".*"
		subjects = append(subjects, processCfg.Subject, subSubject)
	}

	return subjects
}
