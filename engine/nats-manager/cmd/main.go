package main

import (
	"fmt"
	"log"
	"net"

	"github.com/go-logr/zapr"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/konstellation-io/kai/engine/nats-manager/internal/config"
	"github.com/konstellation-io/kai/engine/nats-manager/internal/manager"
	"github.com/konstellation-io/kai/engine/nats-manager/internal/service"
	"github.com/konstellation-io/kai/engine/nats-manager/nats"
	"github.com/konstellation-io/kai/engine/nats-manager/proto/natspb"
)

func main() {
	zapLog, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}

	logger := zapr.NewLogger(zapLog)

	logger.Info("Starting NATS manager")

	config.Initialize()

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", viper.GetInt(config.NatsManagerPort)))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	logger.Info("Connecting to NATS...")

	js, err := nats.InitJetStreamConnection(viper.GetString(config.NatsURL))
	if err != nil {
		log.Fatal(err)
	}

	natsClient := nats.New(logger, js)

	grpcServer := grpc.NewServer()

	natsManager := manager.NewNatsManager(logger, natsClient)
	natsService := service.NewNatsService(logger, natsManager)
	natspb.RegisterNatsManagerServiceServer(grpcServer, natsService)
	reflection.Register(grpcServer)

	logger.Info("Server listening", "port", viper.GetInt(config.NatsManagerPort))

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
