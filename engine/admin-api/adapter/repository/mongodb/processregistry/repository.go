package processregistry

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/logging"
)

const processRegistryCollectionName = "process_registry"

type ProcessRegistryRepoMongoDB struct {
	cfg    *config.Config
	logger logging.Logger
	client *mongo.Client
}

func NewProcessRegistryRepoMongoDB(
	cfg *config.Config,
	logger logging.Logger,
	client *mongo.Client,
) *ProcessRegistryRepoMongoDB {
	processRegistryRepo := &ProcessRegistryRepoMongoDB{
		cfg,
		logger,
		client,
	}

	return processRegistryRepo
}

func (r *ProcessRegistryRepoMongoDB) CreateIndexes(ctx context.Context, productID string) error {
	collection := r.client.Database(productID).Collection(processRegistryCollectionName)

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
