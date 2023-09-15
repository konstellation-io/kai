package processrepository

import (
	"context"
	"errors"
	"time"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/repository"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/logging"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/version"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const registeredProcessesCollectionName = "registered_processes"

type ProcessRepositoryMongoDB struct {
	cfg    *config.Config
	logger logging.Logger
	client *mongo.Client
}

var _ repository.ProcessRepository = (*ProcessRepositoryMongoDB)(nil)

func New(
	cfg *config.Config,
	logger logging.Logger,
	client *mongo.Client,
) *ProcessRepositoryMongoDB {
	return &ProcessRepositoryMongoDB{
		cfg,
		logger,
		client,
	}
}

func (r *ProcessRepositoryMongoDB) CreateIndexes(ctx context.Context, productID string) error {
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

func (r *ProcessRepositoryMongoDB) Create(
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

func (r *ProcessRepositoryMongoDB) ListByProductAndType(
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

func (r *ProcessRepositoryMongoDB) Update(ctx context.Context, productID string, process *entity.RegisteredProcess) error {
	collection := r.client.Database(productID).Collection(registeredProcessesCollectionName)

	versionDTO := mapEntityToDTO(process)
	updateResult, err := collection.ReplaceOne(ctx, bson.M{"_id": process.ID}, versionDTO)

	if updateResult.ModifiedCount == 0 {
		return version.ErrVersionNotFound
	}

	return err
}

func (r *ProcessRepositoryMongoDB) GetByID(ctx context.Context, productID, imageID string) (*entity.RegisteredProcess, error) {
	collection := r.client.Database(productID).Collection(registeredProcessesCollectionName)

	res := collection.FindOne(ctx, bson.M{"_id": imageID})
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return nil, usecase.ErrRegisteredProcessNotFound
		}

		return nil, res.Err()
	}

	var registeredProcess registeredProcessDTO

	err := res.Decode(&registeredProcess)
	if err != nil {
		return nil, err
	}

	return mapDTOToEntity(&registeredProcess), nil
}
