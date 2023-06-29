package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/natspb"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
//
//nolint:dupl // this is not being duplicated
func (n *NatsManagerClient) CreateStreams(
	ctx context.Context,
	productID string,
	version *entity.Version,
) (*entity.VersionStreamsConfig, error) {
	workflows, err := n.mapWorkflowsToDTO(version.Workflows)
	if err != nil {
		return nil, err
	}

	req := natspb.CreateStreamsRequest{
		ProductId:   productID,
		VersionName: version.ID,
		Workflows:   workflows,
	}

	res, err := n.client.CreateStreams(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error creating streams: %w", err)
	}

	return n.dtoToVersionStreamConfig(res.Workflows), err
}

// CreateObjectStores calls nats-manager to create NATS Object Stores for given version.
//
//nolint:dupl // this is not being duplicated
func (n *NatsManagerClient) CreateObjectStores(
	ctx context.Context,
	productID string,
	version *entity.Version,
) (*entity.VersionObjectStoresConfig, error) {
	workflows, err := n.mapWorkflowsToDTO(version.Workflows)
	if err != nil {
		return nil, err
	}

	req := natspb.CreateObjectStoresRequest{
		ProductId:   productID,
		VersionName: version.ID,
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
	productID string,
	version *entity.Version,
) (*entity.KeyValueStoresConfig, error) {
	workflows, err := n.mapWorkflowsToDTO(version.Workflows)
	if err != nil {
		return nil, err
	}

	req := natspb.CreateKeyValueStoresRequest{
		ProductId:   productID,
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
func (n *NatsManagerClient) DeleteStreams(ctx context.Context, productID, versionID string) error {
	req := natspb.DeleteStreamsRequest{
		ProductId:   productID,
		VersionName: versionID,
	}

	_, err := n.client.DeleteStreams(ctx, &req)
	if err != nil {
		return fmt.Errorf("error deleting version %q NATS streams: %w", versionID, err)
	}

	return nil
}

// DeleteObjectStores calls nats-manager to delete NATS Object Stores for given version.
func (n *NatsManagerClient) DeleteObjectStores(ctx context.Context, productID, versionID string) error {
	req := natspb.DeleteObjectStoresRequest{
		ProductId:   productID,
		VersionName: versionID,
	}

	_, err := n.client.DeleteObjectStores(ctx, &req)
	if err != nil {
		return fmt.Errorf("error deleting version %q NATS object stores: %w", versionID, err)
	}

	return nil
}

func (n *NatsManagerClient) mapWorkflowsToDTO(workflows []entity.Workflow) ([]*natspb.Workflow, error) {
	workflowsDTO := make([]*natspb.Workflow, 0, len(workflows))

	for _, w := range workflows {
		processes := make([]*natspb.Process, 0, len(w.Processes))

		for _, process := range w.Processes {
			nodeToAppend := natspb.Process{
				Id:            process.Name,
				Subscriptions: process.Subscriptions,
			}

			if process.ObjectStore != nil {
				scope, err := mapObjectStoreScopeToDTO(process.ObjectStore.Scope)
				if err != nil {
					return nil, err
				}

				nodeToAppend.ObjectStore = &natspb.ObjectStore{
					Name:  process.ObjectStore.Name,
					Scope: scope,
				}
			}

			processes = append(processes, &nodeToAppend)
		}

		workflowsDTO = append(workflowsDTO, &natspb.Workflow{
			Id:        w.Name,
			Processes: processes,
		})
	}

	return workflowsDTO, nil
}

func (n *NatsManagerClient) dtoToVersionStreamConfig(
	workflowsDTO map[string]*natspb.WorkflowStreamConfig,
) *entity.VersionStreamsConfig {
	workflows := make(map[string]entity.WorkflowStreamConfig, len(workflowsDTO))

	for workflow, streamCfg := range workflowsDTO {
		workflows[workflow] = entity.WorkflowStreamConfig{
			Stream:    streamCfg.Stream,
			Processes: n.dtoToNodesStreamConfig(streamCfg.Processes),
		}
	}

	return &entity.VersionStreamsConfig{
		Workflows: workflows,
	}
}

func (n *NatsManagerClient) dtoToNodesStreamConfig(
	processes map[string]*natspb.ProcessStreamConfig,
) map[string]entity.ProcessStreamConfig {
	processesStreamCfg := map[string]entity.ProcessStreamConfig{}

	for process, subjectCfg := range processes {
		processesStreamCfg[process] = entity.ProcessStreamConfig{
			Subject:       subjectCfg.Subject,
			Subscriptions: subjectCfg.Subscriptions,
		}
	}

	return processesStreamCfg
}

func (n *NatsManagerClient) dtoToVersionObjectStoreConfig(
	workflowsDTO map[string]*natspb.WorkflowObjectStoreConfig,
) *entity.VersionObjectStoresConfig {
	workflows := make(map[string]entity.WorkflowObjectStoresConfig, len(workflowsDTO))

	for workflow, objStoreCfg := range workflowsDTO {
		workflows[workflow] = entity.WorkflowObjectStoresConfig{
			Processes: objStoreCfg.Processes,
		}
	}

	return &entity.VersionObjectStoresConfig{
		Workflows: workflows,
	}
}

func (n *NatsManagerClient) dtoToVersionKeyValueStoreConfig(
	projectKeyValueStore string,
	workflows map[string]*natspb.WorkflowKeyValueStoreConfig,
) *entity.KeyValueStoresConfig {
	workflowsKVConfig := make(map[string]*entity.WorkflowKeyValueStores, len(workflows))

	for workflow, kvStoreCfg := range workflows {
		workflowsKVConfig[workflow] = &entity.WorkflowKeyValueStores{
			WorkflowKeyValueStore:   kvStoreCfg.KeyValueStore,
			ProcessesKeyValueStores: kvStoreCfg.Processes,
		}
	}

	return &entity.KeyValueStoresConfig{
		ProductKeyValueStore:    projectKeyValueStore,
		WorkflowsKeyValueStores: workflowsKVConfig,
	}
}

func mapObjectStoreScopeToDTO(scope entity.ObjectStoreScope) (natspb.ObjectStoreScope, error) {
	//nolint:exhaustive // wrong lint rule
	switch scope {
	case "project":
		return natspb.ObjectStoreScope_SCOPE_PROJECT, nil
	case "workflow":
		return natspb.ObjectStoreScope_SCOPE_WORKFLOW, nil
	default:
		//nolint:goerr113 // error needs to be wrapped
		return natspb.ObjectStoreScope_SCOPE_WORKFLOW, errors.New("invalid object store scope")
	}
}
