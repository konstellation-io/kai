package interfaces

import "github.com/konstellation-io/kai/engine/nats-manager/internal/entity"

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/${GOFILE} -package=mocks

type NatsManager interface {
	CreateStreams(productID, versionName string, workflows []entity.Workflow) (entity.WorkflowsStreamsConfig, error)
	CreateObjectStores(productID, versionName string, workflows []entity.Workflow) (entity.WorkflowsObjectStoresConfig, error)
	DeleteStreams(productID, versionName string) error
	DeleteObjectStores(productID, versionName string) error
	CreateKeyValueStores(productID, versionName string, workflows []entity.Workflow) (*entity.VersionKeyValueStores, error)
}
