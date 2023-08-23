package processregistry

import (
	"context"
	"time"

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
	r.logger.Infof("MongoDB creating indexes for %s collection...", processRegistryCollectionName)

	_, err := collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "name", Value: 1},
				{Key: "version", Value: 1},
			},
		},
		{
			Keys: bson.M{"type": 1},
		},
	})

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

func (r *ProcessRegistryRepoMongoDB) ListByProductWithTypeFilter(
	ctx context.Context,
	productID, processType string,
) ([]*entity.ProcessRegistry, error) {
	collection := r.client.Database(productID).Collection(processRegistryCollectionName)

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var processRegistries []*entity.ProcessRegistry

	var filter bson.M
	if processType != "" {
		filter = bson.M{"type": processType}
	} else {
		filter = bson.M{}
	}

	cur, err := collection.Find(ctxWithTimeout, filter)
	if err != nil {
		return processRegistries, err
	}
	defer cur.Close(ctxWithTimeout)

	for cur.Next(ctxWithTimeout) {
		var processRegistryDTO processRegistryDTO

		err = cur.Decode(&processRegistryDTO)
		if err != nil {
			return processRegistries, err
		}

		processRegistries = append(processRegistries, mapDTOToEntity(&processRegistryDTO))
	}

	return processRegistries, nil
}
