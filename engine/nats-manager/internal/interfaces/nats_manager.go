package interfaces

import "github.com/konstellation-io/kai/engine/nats-manager/internal/entity"

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/${GOFILE} -package=mocks

type NatsManager interface {
	CreateStreams(productID, versionTag string, workflows []entity.Workflow) (entity.WorkflowsStreamsConfig, error)
	CreateObjectStores(productID, versionTag string, workflows []entity.Workflow) (entity.WorkflowsObjectStoresConfig, error)
	DeleteStreams(productID, versionTag string) error
	DeleteObjectStores(productID, versionTag string) error
	CreateKeyValueStores(productID, versionTag string, workflows []entity.Workflow) (*entity.VersionKeyValueStores, error)
}
