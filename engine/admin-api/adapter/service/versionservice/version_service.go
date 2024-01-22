package versionservice

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/versionpb"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service"
)

//go:generate mockgen -source=../proto/versionpb/version_grpc.pb.go -destination=../../../mocks/${GOFILE} -package=mocks

// TODO: move this to viper config
//
//nolint:godox // To be done.
const _requestTimeout = 5 * time.Minute

type K8sVersionService struct {
	client versionpb.VersionServiceClient
	logger logr.Logger
}

var _ service.VersionService = (*K8sVersionService)(nil)

func New(logger logr.Logger, client versionpb.VersionServiceClient) (*K8sVersionService, error) {
	return &K8sVersionService{
		client,
		logger,
	}, nil
}

// Start creates the version resources in k8s.
func (k *K8sVersionService) Start(
	ctx context.Context,
	product *entity.Product,
	version *entity.Version,
	versionConfig *entity.VersionStreamingResources,
) error {
	wf, err := mapWorkflowsToDTO(version.Workflows, versionConfig)
	if err != nil {
		return fmt.Errorf("map workflows to DTO: %w", err)
	}

	req := versionpb.StartRequest{
		ProductId:            product.ID,
		VersionTag:           version.Tag,
		Workflows:            wf,
		GlobalKeyValueStore:  versionConfig.KeyValueStores.GlobalKeyValueStore,
		VersionKeyValueStore: versionConfig.KeyValueStores.VersionKeyValueStore,
		MinioConfiguration: &versionpb.MinioConfiguration{
			Bucket: product.MinioConfiguration.Bucket,
		},
		ServiceAccount: &versionpb.ServiceAccount{
			Username: product.ServiceAccount.Username,
			Password: product.ServiceAccount.Password,
		},
	}

	_, err = k.client.Start(ctx, &req)

	return err
}

func (k *K8sVersionService) Stop(ctx context.Context, productID string, version *entity.Version) error {
	req := versionpb.StopRequest{
		Product:    productID,
		VersionTag: version.Tag,
	}

	_, err := k.client.Stop(ctx, &req)
	if err != nil {
		return fmt.Errorf("stop version %q in product %q error: %w", version.Tag, productID, err)
	}

	return nil
}

func (k *K8sVersionService) Unpublish(ctx context.Context, productID string, version *entity.Version) error {
	req := versionpb.UnpublishRequest{
		Product:    productID,
		VersionTag: version.Tag,
	}

	ctx, cancel := context.WithTimeout(ctx, _requestTimeout)
	defer cancel()

	_, err := k.client.Unpublish(ctx, &req)

	return err
}

func (k *K8sVersionService) Publish(ctx context.Context, productID, version string) error {
	req := versionpb.PublishRequest{
		Product:    productID,
		VersionTag: version,
	}

	ctx, cancel := context.WithTimeout(ctx, _requestTimeout)
	defer cancel()

	_, err := k.client.Publish(ctx, &req)

	return err
}

func (k *K8sVersionService) WatchProcessStatus(ctx context.Context, productID, versionTag string) (<-chan *entity.Process, error) {
	stream, err := k.client.WatchProcessStatus(ctx, &versionpb.ProcessStatusRequest{
		VersionTag: versionTag,
		ProductId:  productID,
	})
	if err != nil {
		return nil, fmt.Errorf("version status opening stream: %w", err)
	}

	ch := make(chan *entity.Process, 1)

	go func() {
		defer close(ch)

		for {
			k.logger.V(2).Info("[VersionService.WatchProcessStatus] waiting for stream.Recv()...")

			msg, err := stream.Recv()

			if errors.Is(stream.Context().Err(), context.Canceled) {
				k.logger.V(2).Info("[VersionService.WatchProcessStatus] Context canceled.")
				return
			}

			if errors.Is(err, io.EOF) {
				k.logger.V(2).Info("[VersionService.WatchProcessStatus] EOF msg received.")
				return
			}

			if err != nil {
				k.logger.Error(err, "[VersionService.WatchProcessStatus] Unexpected error")
				return
			}

			k.logger.V(2).Info("[VersionService.WatchProcessStatus] Message received")

			status := entity.ProcessStatus(msg.GetStatus())
			if !status.IsValid() {
				k.logger.Error(err, "[VersionService.WatchProcessStatus] Invalid node status", "status", status)
				continue
			}

			ch <- &entity.Process{
				Name:   msg.Name,
				Status: status,
			}
		}
	}()

	return ch, nil
}

func (k *K8sVersionService) RegisterProcess(ctx context.Context, productID, processID, processImage string) (string, error) {
	res, err := k.client.RegisterProcess(ctx, &versionpb.RegisterProcessRequest{
		ProductId:    productID,
		ProcessId:    processID,
		ProcessImage: processImage,
	})
	if err != nil {
		return "", err
	}

	return res.ImageId, nil
}
