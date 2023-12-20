package processregistry

import (
	"context"

	"github.com/go-logr/logr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

const processRegistryCollectionName = "process_registry"

type ProcessRegistryRepoMongoDB struct {
	logger logr.Logger
	client *mongo.Client
}

func NewProcessRegistryRepoMongoDB(
	logger logr.Logger,
	client *mongo.Client,
) *ProcessRegistryRepoMongoDB {
	processRegistryRepo := &ProcessRegistryRepoMongoDB{
		logger,
		client,
	}

	return processRegistryRepo
}

func (r *ProcessRegistryRepoMongoDB) CreateIndexes(ctx context.Context, productID string) error {
	collection := r.client.Database(productID).Collection(processRegistryCollectionName)
	r.logger.Info("MongoDB creating indexes for %s collection", "collection", processRegistryCollectionName)

	indexes := []mongo.IndexModel{
		{
			Keys: bson.M{
				"name": 1,
			},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)

	return err
}

func (r *ProcessRegistryRepoMongoDB) Create(
	productID string,
	newProcessRegistry *entity.ProcessRegistry,
) (*entity.ProcessRegistry, error) {
	collection := r.client.Database(productID).Collection(processRegistryCollectionName)

	processRegistryDTO := mapEntityToDTO(newProcessRegistry)

	res, err := collection.InsertOne(context.Background(), processRegistryDTO)
	if err != nil {
		return nil, err
	}

	processRegistryDTO.ID = res.InsertedID.(string)

	savedProcessRegistry := mapDTOToEntity(processRegistryDTO)

	return savedProcessRegistry, nil
}
