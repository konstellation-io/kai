package bootstrap

import (
	"fmt"
	"net"

	"github.com/go-logr/zapr"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/usecase"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/config"
	internalgrpc "github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/grpc"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/grpc/proto/versionpb"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/kube"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Run() error {
	// Init viper config
	if err := config.Init("config.yaml"); err != nil {
		return err
	}

	port := 50051

	serverAddress := fmt.Sprintf("0.0.0.0:%d", port)

	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		return err
	}

	zapLog, err := zap.NewDevelopment()
	if err != nil {
		return err
	}

	logger := zapr.NewLogger(zapLog)

	client, err := kube.NewClientset()
	if err != nil {
		return err
	}

	s := grpc.NewServer()

	k8sContainerService := kube.NewK8sContainerService(logger, client)
	starter := usecase.NewVersionStarter(logger, k8sContainerService)
	stopper := usecase.NewVersionStopper(logger, k8sContainerService)

	versionService := internalgrpc.NewVersionService(logger, starter, stopper)

	versionpb.RegisterVersionServiceServer(s, versionService)
	reflection.Register(s)

	logger.Info("Server listening", "port", port)

	if err := s.Serve(listener); err != nil {
		logger.Error(err, "Failed to serve")
	}

	return nil
}
