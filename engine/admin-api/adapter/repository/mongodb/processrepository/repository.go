package processrepository

import (
	"context"
	"errors"
	"time"

	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/repository"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/process"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/version"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const registeredProcessesCollectionName = "registered_processes"

type MongoDBProcessRepository struct {
	logger logr.Logger
	client *mongo.Client
}

var _ repository.ProcessRepository = (*MongoDBProcessRepository)(nil)

func New(
	logger logr.Logger,
	client *mongo.Client,
) *MongoDBProcessRepository {
	return &MongoDBProcessRepository{
		logger,
		client,
	}
}

func (r *MongoDBProcessRepository) CreateIndexes(ctx context.Context, productID string) error {
	collection := r.client.Database(productID).Collection(registeredProcessesCollectionName)
	r.logger.Info("MongoDB creating indexes", "collection", registeredProcessesCollectionName)

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

func (r *MongoDBProcessRepository) Create(
	ctx context.Context,
	productID string,
	newRegisteredProcess *entity.RegisteredProcess,
) error {
	collection := r.client.Database(productID).Collection(registeredProcessesCollectionName)

	registeredProcessDTO := mapEntityToDTO(newRegisteredProcess)

	_, err := collection.InsertOne(ctx, registeredProcessDTO)
	if err != nil {
		return err
	}

	return nil
}

func (r *MongoDBProcessRepository) SearchByProduct(
	ctx context.Context,
	product string,
	filter repository.SearchFilter,
) ([]*entity.RegisteredProcess, error) {
	return r.searchInDatabaseWithFilter(ctx, product, r.getSearchMongoFilter(filter))
}

func (r *MongoDBProcessRepository) GlobalSearch(ctx context.Context, filter repository.SearchFilter) ([]*entity.RegisteredProcess, error) {
	return r.searchInDatabaseWithFilter(ctx, viper.GetString(config.MongoDBKaiDatabaseKey), r.getSearchMongoFilter(filter))
}

func (r *MongoDBProcessRepository) Update(ctx context.Context, productID string, p *entity.RegisteredProcess) error {
	collection := r.client.Database(productID).Collection(registeredProcessesCollectionName)

	versionDTO := mapEntityToDTO(p)
	updateResult, err := collection.ReplaceOne(ctx, bson.M{"_id": p.ID}, versionDTO)

	if updateResult.ModifiedCount == 0 {
		return version.ErrVersionNotFound
	}

	return err
}

func (r *MongoDBProcessRepository) GetByID(ctx context.Context, productID, imageID string) (*entity.RegisteredProcess, error) {
	collection := r.client.Database(productID).Collection(registeredProcessesCollectionName)

	res := collection.FindOne(ctx, bson.M{"_id": imageID})
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return nil, process.ErrRegisteredProcessNotFound
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

func (r *MongoDBProcessRepository) searchInDatabaseWithFilter(
	ctx context.Context,
	database string,
	filter bson.M,
) ([]*entity.RegisteredProcess, error) {
	collection := r.client.Database(database).Collection(registeredProcessesCollectionName)

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var registeredProcesses []*entity.RegisteredProcess

	cur, err := collection.Find(ctxWithTimeout, filter)
	if err != nil {
		return registeredProcesses, err
	}
	defer cur.Close(ctxWithTimeout)

	for cur.Next(ctxWithTimeout) {
		var dto registeredProcessDTO

		err = cur.Decode(&dto)
		if err != nil {
			return registeredProcesses, err
		}

		registeredProcesses = append(registeredProcesses, mapDTOToEntity(&dto))
	}

	return registeredProcesses, nil
}

func (r *MongoDBProcessRepository) getSearchMongoFilter(searchFilter repository.SearchFilter) bson.M {
	var filter bson.M

	if searchFilter.ProcessType != "" {
		filter["type"] = searchFilter.ProcessType
	}

	if searchFilter.ProcessID != "" {
		filter["_id"] = searchFilter.ProcessID
	}

	return filter
}
