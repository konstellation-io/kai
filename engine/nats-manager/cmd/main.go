package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net"

	"github.com/konstellation-io/kai/libs/simplelogger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/konstellation-io/kai/engine/nats-manager/internal/config"
	"github.com/konstellation-io/kai/engine/nats-manager/internal/manager"
	"github.com/konstellation-io/kai/engine/nats-manager/internal/service"
	"github.com/konstellation-io/kai/engine/nats-manager/nats"
	"github.com/konstellation-io/kai/engine/nats-manager/proto/natspb"
)

func main() {
	logger := simplelogger.New(simplelogger.LevelDebug)
	logger.Info("Starting NATS manager")

	config.Initialize()

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", viper.GetInt(config.NATS_MANAGER_PORT)))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	logger.Info("Connecting to NATS...")

	js, err := nats.InitJetStreamConnection(viper.GetString(config.NATS_URL))
	if err != nil {
		log.Fatal(err)
	}

	natsClient := nats.New(logger, js)

	grpcServer := grpc.NewServer()

	natsManager := manager.NewNatsManager(logger, natsClient)
	natsService := service.NewNatsService(logger, natsManager)
	natspb.RegisterNatsManagerServiceServer(grpcServer, natsService)
	reflection.Register(grpcServer)

	logger.Infof("Server listening on port: %d", viper.GetInt(config.NATS_MANAGER_PORT))

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
