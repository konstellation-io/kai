package mongodb

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
)

type MongoDB struct {
	logger logr.Logger
	client *mongo.Client
}

func NewMongoDB(logger logr.Logger) *MongoDB {
	return &MongoDB{
		logger,
		nil,
	}
}

func (m *MongoDB) Connect() (*mongo.Client, error) {
	m.logger.V(2).Info("Connecting to MongoDB", "endpoint", viper.GetString(config.MongoDBEndpointKey))

	client, err := mongo.NewClient(options.Client().ApplyURI(viper.GetString(config.MongoDBEndpointKey)))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	// Call Ping to verify that the deployment is up and the Client was configured successfully.
	ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	m.logger.V(2).Info("MongoDB ping...")

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	m.logger.Info("MongoDB connected")
	m.client = client

	return client, nil
}

func (m *MongoDB) Disconnect() {
	m.logger.Info("MongoDB disconnecting...")

	if m.client == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err := m.client.Disconnect(ctx)
	if err != nil {
		m.logger.Error(err, "Error closing MongoDB connection")
		return
	}

	m.logger.Info("Connection to MongoDB closed.")
}
