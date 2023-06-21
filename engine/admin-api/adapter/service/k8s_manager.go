package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/versionpb"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/logging"
	"github.com/konstellation-io/krt/pkg/krt"
)

type K8sVersionClient struct {
	cfg    *config.Config
	client versionpb.VersionServiceClient
	logger logging.Logger
}

func NewK8sVersionClient(cfg *config.Config, logger logging.Logger) (*K8sVersionClient, error) {
	cc, err := grpc.Dial(cfg.Services.K8sManager, grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := versionpb.NewVersionServiceClient(cc)

	if err != nil {
		return nil, err
	}

	return &K8sVersionClient{
		cfg,
		client,
		logger,
	}, nil
}

// Start creates the version resources in k8s.
func (k *K8sVersionClient) Start(
	ctx context.Context,
	productID string,
	version *entity.Version,
	versionConfig *entity.VersionConfig,
) error {
	configVars := versionToConfig(version)
	wf, err := versionToWorkflows(version, versionConfig)

	if err != nil {
		return err
	}

	//nolint:godox // To be done.
	// TODO: Update to new proto
	req := versionpb.StartRequest{
		ProductId:      productID,
		VersionId:      version.ID,
		VersionName:    version.Name,
		Config:         configVars,
		Workflows:      wf,
		MongoUri:       k.cfg.MongoDB.Address,
		MongoDbName:    k.cfg.MongoDB.DBName,
		MongoKrtBucket: k.cfg.MongoDB.KRTBucket,
		InfluxUri:      fmt.Sprintf("http://%s-influxdb:8086", k.cfg.ReleaseName),
		KeyValueStore:  versionConfig.KeyValueStoresConfig.ProjectKeyValueStore,
	}

	_, err = k.client.Start(ctx, &req)

	return err
}

func (k *K8sVersionClient) Stop(ctx context.Context, productID string, version *entity.Version) error {
	workflowEntrypoints := make([]string, 0)

	for _, w := range version.Workflows {
		for _, p := range w.Processes {
			if p.Type == krt.ProcessTypeTrigger {
				workflowEntrypoints = append(workflowEntrypoints, p.Name)
			}
		}
	}

	req := versionpb.VersionInfo{
		Name:      version.Name,
		ProductId: productID,
		Workflows: workflowEntrypoints,
	}

	_, err := k.client.Stop(ctx, &req)
	if err != nil {
		return fmt.Errorf("stop version %q error: %w", version.Name, err)
	}

	return nil
}

func (k *K8sVersionClient) UpdateConfig(productID string, version *entity.Version) error {
	configVars := versionToConfig(version)

	req := versionpb.UpdateConfigRequest{
		ProductId:   productID,
		VersionName: version.Name,
		Config:      configVars,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	_, err := k.client.UpdateConfig(ctx, &req)

	return err
}

func (k *K8sVersionClient) Unpublish(productID string, version *entity.Version) error {
	req := versionpb.VersionInfo{
		Name:      version.Name,
		ProductId: productID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	_, err := k.client.Unpublish(ctx, &req)

	return err
}

func (k *K8sVersionClient) Publish(productID string, version *entity.Version) error {
	req := versionpb.VersionInfo{
		Name:      version.Name,
		ProductId: productID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	_, err := k.client.Publish(ctx, &req)

	return err
}

func versionToConfig(version *entity.Version) []*versionpb.Config {
	configVars := make([]*versionpb.Config, len(version.Config))
	idx := 0

	for k, v := range version.Config {
		configVars[idx] = &versionpb.Config{
			Key:   k,
			Value: v,
		}
		idx++
	}

	return configVars
}

// TODO: Transform to new krt when protobuf is updated.
//
//nolint:godox // To be done.
func versionToWorkflows(version *entity.Version, versionConfig *entity.VersionConfig) ([]*versionpb.Workflow, error) {
	wf := make([]*versionpb.Workflow, len(version.Workflows))

	for i, w := range version.Workflows {
		workflowStreamConfig, err := versionConfig.GetWorkflowStreamConfig(w.Name)

		if err != nil {
			return nil, fmt.Errorf("error translating version in workflow %q: %w", w.Name, err)
		}

		workflowKeyValueStoresConfig, err := versionConfig.GetWorkflowKeyValueStoresConfig(w.Name)

		if err != nil {
			return nil, fmt.Errorf("error getting workflow %q key-value store: %w", w.Name, err)
		}

		process := make([]*versionpb.Workflow_Node, len(w.Processes))

		for j, n := range w.Processes {
			processStreamCfg, err := workflowStreamConfig.GetProcessStreamConfig(n.Name)
			if err != nil {
				return nil, fmt.Errorf("error getting stream configuration from process %q: %w", n.Name, err)
			}

			processKeyValueStore, err := workflowKeyValueStoresConfig.GetProcessKeyValueStore(n.Name)
			if err != nil {
				return nil, fmt.Errorf("error translating version in workflow %q: %w", w.Name, err)
			}

			if err != nil {
				return nil, fmt.Errorf("error getting process key-value store config: %w", err)
			}

			process[j] = &versionpb.Workflow_Node{
				Name:          n.Name,
				Image:         n.Image,
				Subscriptions: processStreamCfg.Subscriptions,
				Subject:       processStreamCfg.Subject,
				ObjectStore:   versionConfig.GetProcessObjectStoreConfig(w.Name, n.Name),
				KeyValueStore: processKeyValueStore,
			}
		}

		wf[i] = &versionpb.Workflow{
			Name:          w.Name,
			Stream:        workflowStreamConfig.Stream,
			KeyValueStore: workflowKeyValueStoresConfig.WorkflowKeyValueStore,
		}
	}

	return wf, nil
}

func (k *K8sVersionClient) WatchProcessStatus(
	ctx context.Context,
	productID,
	versionName string,
) (<-chan *krt.Process, error) {
	stream, err := k.client.WatchNodeStatus(ctx, &versionpb.NodeStatusRequest{
		VersionName: versionName,
		ProductId:   productID,
	})
	if err != nil {
		return nil, fmt.Errorf("version status opening stream: %w", err)
	}

	ch := make(chan *krt.Process, 1)

	go func() {
		defer close(ch)

		for {
			k.logger.Debug("[VersionService.WatchProcessStatus] waiting for stream.Recv()...")

			msg, err := stream.Recv()

			if errors.Is(stream.Context().Err(), context.Canceled) {
				k.logger.Debug("[VersionService.WatchProcessStatus] Context canceled.")
				return
			}

			if errors.Is(err, io.EOF) {
				k.logger.Debug("[VersionService.WatchProcessStatus] EOF msg received.")
				return
			}

			if err != nil {
				k.logger.Errorf("[VersionService.WatchProcessStatus] Unexpected error: %s", err)
				return
			}

			k.logger.Debug("[VersionService.WatchProcessStatus] Message received")

			status := krt.ProcessStatus(msg.GetStatus())
			if !status.IsValid() {
				k.logger.Errorf("[VersionService.WatchProcessStatus] Invalid process status: %s", status)
				continue
			}

			ch <- &krt.Process{
				Name:   msg.GetName(),
				Status: status,
			}
		}
	}()

	return ch, nil
}
