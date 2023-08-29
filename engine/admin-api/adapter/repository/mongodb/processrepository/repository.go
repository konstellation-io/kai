package processrepository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/logging"
)

const registeredProcessesCollectionName = "registered_processes"

type RegisteredProcessRepoMongoDB struct {
	cfg    *config.Config
	logger logging.Logger
	client *mongo.Client
}

func New(
	cfg *config.Config,
	logger logging.Logger,
	client *mongo.Client,
) *RegisteredProcessRepoMongoDB {
	registeredProcessRepo := &RegisteredProcessRepoMongoDB{
		cfg,
		logger,
		client,
	}

	return registeredProcessRepo
}

func (r *RegisteredProcessRepoMongoDB) CreateIndexes(ctx context.Context, productID string) error {
	collection := r.client.Database(productID).Collection(registeredProcessesCollectionName)
	r.logger.Infof("MongoDB creating indexes for %s collection...", registeredProcessesCollectionName)

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

func (r *RegisteredProcessRepoMongoDB) Create(
	productID string,
	newRegisteredProcess *entity.RegisteredProcess,
) (*entity.RegisteredProcess, error) {
	collection := r.client.Database(productID).Collection(registeredProcessesCollectionName)

	registeredProcessDTO := mapEntityToDTO(newRegisteredProcess)

	res, err := collection.InsertOne(context.Background(), registeredProcessDTO)
	if err != nil {
		return nil, err
	}

	registeredProcessDTO.ID = res.InsertedID.(string)

	savedRegisteredProcess := mapDTOToEntity(registeredProcessDTO)

	return savedRegisteredProcess, nil
}

func (r *RegisteredProcessRepoMongoDB) ListByProductAndType(
	ctx context.Context,
	productID, processType string,
) ([]*entity.RegisteredProcess, error) {
	collection := r.client.Database(productID).Collection(registeredProcessesCollectionName)

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var registeredProcesses []*entity.RegisteredProcess

	var filter bson.M
	if processType != "" {
		filter = bson.M{"type": processType}
	} else {
		filter = bson.M{}
	}

	cur, err := collection.Find(ctxWithTimeout, filter)
	if err != nil {
		return registeredProcesses, err
	}
	defer cur.Close(ctxWithTimeout)

	for cur.Next(ctxWithTimeout) {
		var registeredProcessDTO registeredProcessDTO

		err = cur.Decode(&registeredProcessDTO)
		if err != nil {
			return registeredProcesses, err
		}

		registeredProcesses = append(registeredProcesses, mapDTOToEntity(&registeredProcessDTO))
	}

	return registeredProcesses, nil
}
