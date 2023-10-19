package interfaces

import (
	"regexp"

	"github.com/konstellation-io/kai/engine/nats-manager/internal/entity"
)

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/${GOFILE} -package=mocks

type NatsClient interface {
	CreateStream(streamConfig *entity.StreamConfig) error
	CreateObjectStore(objectStore string) error
	GetObjectStoreNames(optFilter ...*regexp.Regexp) ([]string, error)
	GetStreamNames(optFilter ...*regexp.Regexp) ([]string, error)
	DeleteStream(stream string) error
	DeleteObjectStore(stream string) error
	CreateKeyValueStore(keyValueStore string) error
	UpdateConfiguration(keyValueStore string, keyValueConfig map[string]string) error
}
