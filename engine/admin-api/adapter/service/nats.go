package service

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/grpc"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/natspb"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/logging"
	"github.com/konstellation-io/krt/pkg/krt"
)

type NatsManagerClient struct {
	cfg    *config.Config
	client natspb.NatsManagerServiceClient
	logger logging.Logger
}

func NewNatsManagerClient(cfg *config.Config, logger logging.Logger) (*NatsManagerClient, error) {
	cc, err := grpc.Dial(cfg.Services.NatsManager, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	client := natspb.NewNatsManagerServiceClient(cc)

	return &NatsManagerClient{
		cfg,
		client,
		logger,
	}, nil
}

// CreateStreams calls nats-manager to create NATS streams for given version.
func (n *NatsManagerClient) CreateStreams(
	ctx context.Context,
	runtimeID string,
	version *entity.Version,
) (*entity.VersionStreamsConfig, error) {
	workflows, err := n.getWorkflowsFromVersion(version)
	if err != nil {
		return nil, err
	}

	req := natspb.CreateStreamsRequest{
		RuntimeId:   runtimeID,
		VersionName: version.Name,
		Workflows:   workflows,
	}

	res, err := n.client.CreateStreams(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error creating streams: %w", err)
	}

	return n.dtoToVersionStreamConfig(res.Workflows), err
}

// CreateObjectStores calls nats-manager to create NATS Object Stores for given version.
func (n *NatsManagerClient) CreateObjectStores(
	ctx context.Context,
	runtimeID string,
	version *entity.Version,
) (*entity.VersionObjectStoresConfig, error) {
	workflows, err := n.getWorkflowsFromVersion(version)
	if err != nil {
		return nil, err
	}

	req := natspb.CreateObjectStoresRequest{
		RuntimeId:   runtimeID,
		VersionName: version.Name,
		Workflows:   workflows,
	}

	res, err := n.client.CreateObjectStores(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error creating object stores: %w", err)
	}

	return n.dtoToVersionObjectStoreConfig(res.Workflows), err
}

// CreateKeyValueStores calls nats-manager to create NATS Key Value Stores for given version.
func (n *NatsManagerClient) CreateKeyValueStores(
	ctx context.Context,
	runtimeID string,
	version *entity.Version,
) (*entity.KeyValueStoresConfig, error) {
	workflows, err := n.getWorkflowsFromVersion(version)
	if err != nil {
		return nil, err
	}

	req := natspb.CreateKeyValueStoresRequest{
		RuntimeId:   runtimeID,
		VersionName: version.Name,
		Workflows:   workflows,
	}

	res, err := n.client.CreateKeyValueStores(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error creating key value stores: %w", err)
	}

	return n.dtoToVersionKeyValueStoreConfig(res.KeyValueStore, res.Workflows), err
}

// DeleteStreams calls nats-manager to delete NATS streams for given version.
func (n *NatsManagerClient) DeleteStreams(ctx context.Context, runtimeID, versionName string) error {
	req := natspb.DeleteStreamsRequest{
		RuntimeId:   runtimeID,
		VersionName: versionName,
	}

	_, err := n.client.DeleteStreams(ctx, &req)
	if err != nil {
		return fmt.Errorf("error deleting version %q NATS streams: %w", versionName, err)
	}

	return nil
}

// DeleteObjectStores calls nats-manager to delete NATS Object Stores for given version.
func (n *NatsManagerClient) DeleteObjectStores(ctx context.Context, runtimeID, versionName string) error {
	req := natspb.DeleteObjectStoresRequest{
		RuntimeId:   runtimeID,
		VersionName: versionName,
	}

	_, err := n.client.DeleteObjectStores(ctx, &req)
	if err != nil {
		return fmt.Errorf("error deleting version %q NATS object stores: %w", versionName, err)
	}

	return nil
}

// TODO migrate to new krt once proto is updated
func (n *NatsManagerClient) getWorkflowsFromVersion(version *entity.Version) ([]*natspb.Workflow, error) {
	var workflows = make([]*natspb.Workflow, 0, len(version.Workflows))

	for _, w := range version.Workflows {
		processes := make([]*natspb.Node, 0, len(w.Processes))

		for _, process := range w.Processes {
			processToAppend := natspb.Node{
				Name:          process.Name,
				Subscriptions: process.Subscriptions,
			}

			if process.ObjectStore != nil {
				scope, err := translateObjectStoreEnum(process.ObjectStore.Scope)
				if err != nil {
					return nil, err
				}

				processToAppend.ObjectStore = &natspb.ObjectStore{
					Name:  process.ObjectStore.Name,
					Scope: scope,
				}
			}

			processes = append(processes, &processToAppend)
		}

		workflows = append(workflows, &natspb.Workflow{
			Name:  w.Name,
			Nodes: processes,
		})
	}

	return workflows, nil
}

func (n *NatsManagerClient) dtoToVersionStreamConfig(
	workflows map[string]*natspb.CreateStreamsResponse_WorkflowStreamConfig,
) *entity.VersionStreamsConfig {
	workflowsConfig := map[string]*entity.WorkflowStreamConfig{}
	for workflow, streamCfg := range workflows {
		workflowsConfig[workflow] = &entity.WorkflowStreamConfig{
			Stream:            streamCfg.Stream,
			Processs:          n.dtoToProcesssStreamConfig(streamCfg.Nodes),
			EntrypointSubject: streamCfg.EntrypointSubject,
		}
	}

	return &entity.VersionStreamsConfig{
		Workflows: workflowsConfig,
	}
}

func (n *NatsManagerClient) dtoToProcesssStreamConfig(
	processs map[string]*natspb.CreateStreamsResponse_NodeStreamConfig,
) map[string]*entity.ProcessStreamConfig {
	processsStreamCfg := map[string]*entity.ProcessStreamConfig{}

	for process, subjectCfg := range processs {
		processsStreamCfg[process] = &entity.ProcessStreamConfig{
			Subject:       subjectCfg.Subject,
			Subscriptions: subjectCfg.Subscriptions,
		}
	}

	return processsStreamCfg
}

func (n *NatsManagerClient) dtoToVersionObjectStoreConfig(
	workflows map[string]*natspb.CreateObjectStoresResponse_WorkflowObjectStoreConfig,
) *entity.VersionObjectStoresConfig {
	workflowsObjStoreConfig := entity.WorkflowsObjectStoresConfig{}

	for workflow, objStoreCfg := range workflows {
		workflowsObjStoreConfig[workflow] = objStoreCfg.Nodes
	}

	return &entity.VersionObjectStoresConfig{
		Workflows: workflowsObjStoreConfig,
	}
}

func (n *NatsManagerClient) dtoToVersionKeyValueStoreConfig(
	projectKeyValueStore string,
	workflows map[string]*natspb.CreateKeyValueStoreResponse_WorkflowKeyValueStoreConfig,
) *entity.KeyValueStoresConfig {
	workflowsKVConfig := map[string]*entity.WorkflowKeyValueStores{}

	for workflow, kvStoreCfg := range workflows {
		workflowsKVConfig[workflow] = &entity.WorkflowKeyValueStores{
			WorkflowKeyValueStore:  kvStoreCfg.KeyValueStore,
			ProcesssKeyValueStores: kvStoreCfg.Nodes,
		}
	}

	return &entity.KeyValueStoresConfig{
		ProjectKeyValueStore:    projectKeyValueStore,
		WorkflowsKeyValueStores: workflowsKVConfig,
	}
}

// TODO, migrate project to product in proto
func translateObjectStoreEnum(scope krt.ObjectStoreScope) (natspb.ObjectStoreScope, error) {
	switch scope {
	case krt.ObjectStoreScopeProduct:
		return natspb.ObjectStoreScope_SCOPE_PROJECT, nil
	case krt.ObjectStoreScopeWorkflow:
		return natspb.ObjectStoreScope_SCOPE_WORKFLOW, nil
	default:
		//nolint:goerr113 // error needs to be wrapped
		return natspb.ObjectStoreScope_SCOPE_WORKFLOW, errors.New("invalid object store scope")
	}
}
