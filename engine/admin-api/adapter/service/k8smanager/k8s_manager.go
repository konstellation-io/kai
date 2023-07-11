package k8smanager

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/versionpb"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/logging"
)

//go:generate mockgen -source=../proto/versionpb/version_grpc.pb.go -destination=../../../mocks/${GOFILE} -package=mocks

const _requestTimeout = 1 * time.Minute

type K8sVersionClient struct {
	cfg    *config.Config
	client versionpb.VersionServiceClient
	logger logging.Logger
}

var _ service.K8sService = (*K8sVersionClient)(nil)

func NewK8sVersionClient(cfg *config.Config, logger logging.Logger, client versionpb.VersionServiceClient) (*K8sVersionClient, error) {
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
	wf, err := mapWorkflowsToDTO(version.Workflows, versionConfig)
	if err != nil {
		return fmt.Errorf("map workflows to DTO: %w", err)
	}

	req := versionpb.StartRequest{
		ProductId:     productID,
		VersionTag:    version.Version,
		Workflows:     wf,
		KeyValueStore: versionConfig.KeyValueStoresConfig.KeyValueStore,
	}

	_, err = k.client.Start(ctx, &req)

	return err
}

func (k *K8sVersionClient) Stop(ctx context.Context, productID string, version *entity.Version) error {
	req := versionpb.StopRequest{
		Product:    productID,
		VersionTag: version.Version,
	}

	_, err := k.client.Stop(ctx, &req)
	if err != nil {
		return fmt.Errorf("stop version %q in product %q error: %w", version.Version, productID, err)
	}

	return nil
}

func (k *K8sVersionClient) Unpublish(ctx context.Context, productID string, version *entity.Version) error {
	req := versionpb.UnpublishRequest{
		Product:    productID,
		VersionTag: version.Version,
	}

	ctx, cancel := context.WithTimeout(ctx, _requestTimeout)
	defer cancel()

	_, err := k.client.Unpublish(ctx, &req)

	return err
}

func (k *K8sVersionClient) Publish(ctx context.Context, productID string, version *entity.Version) error {
	req := versionpb.PublishRequest{
		Product:    productID,
		VersionTag: version.Version,
	}

	ctx, cancel := context.WithTimeout(ctx, _requestTimeout)
	defer cancel()

	_, err := k.client.Publish(ctx, &req)

	return err
}

func (k *K8sVersionClient) WatchProcessStatus(ctx context.Context, productID, VersionTag string) (<-chan *entity.Process, error) {
	stream, err := k.client.WatchProcessStatus(ctx, &versionpb.ProcessStatusRequest{
		VersionTag: VersionTag,
		ProductId:  productID,
	})
	if err != nil {
		return nil, fmt.Errorf("version status opening stream: %w", err)
	}

	ch := make(chan *entity.Process, 1)

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

			status := entity.ProcessStatus(msg.GetStatus())
			if !status.IsValid() {
				k.logger.Errorf("[VersionService.WatchProcessStatus] Invalid node status: %s", status)
				continue
			}

			ch <- &entity.Process{
				ID:     msg.ProcessId,
				Name:   msg.Name,
				Status: status,
			}
		}
	}()

	return ch, nil
}
