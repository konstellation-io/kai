package bootstrap

import (
	"fmt"
	"net"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/usecase"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/config"
	internalgrpc "github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/grpc"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/grpc/proto/versionpb"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/kube"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Run() error {
	if err := config.Init("config.yaml"); err != nil {
		return err
	}

	logger, err := initLogger()
	if err != nil {
		return err
	}

	s, err := initGrpcServer(logger)
	if err != nil {
		return err
	}

	return startServer(logger, s)
}

func initLogger() (logr.Logger, error) {
	zapLog, err := zap.NewDevelopment()
	if err != nil {
		return logr.Logger{}, err
	}

	return zapr.NewLogger(zapLog), nil
}

func initGrpcServer(logger logr.Logger) (*grpc.Server, error) {
	client, err := kube.NewClientset()
	if err != nil {
		return nil, err
	}

	s := grpc.NewServer()

	k8sContainerService := kube.NewK8sContainerService(logger, client)
	starter := usecase.NewVersionStarter(logger, k8sContainerService)
	stopper := usecase.NewVersionStopper(logger, k8sContainerService)

	versionService := internalgrpc.NewVersionService(logger, starter, stopper)

	versionpb.RegisterVersionServiceServer(s, versionService)
	reflection.Register(s)

	return s, nil
}

func startServer(logger logr.Logger, s *grpc.Server) error {
	serverAddress := fmt.Sprintf("0.0.0.0:%d", viper.GetInt("server.port"))

	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		return err
	}

	logger.Info("Server listening", "port", viper.GetInt("server.port"))

	if err := s.Serve(listener); err != nil {
		logger.Error(err, "Failed to serve")
		return err
	}

	return nil
}
