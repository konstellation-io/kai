package mongodb

import (
	"context"
	"gitlab.com/konstellation/konstellation-ce/kre/admin-api/adapter/config"
	"gitlab.com/konstellation/konstellation-ce/kre/admin-api/domain/entity"
	"gitlab.com/konstellation/konstellation-ce/kre/admin-api/domain/usecase/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type RuntimeRepoMongoDB struct {
	cfg        *config.Config
	logger     logging.Logger
	collection *mongo.Collection
}

func NewRuntimeRepoMongoDB(cfg *config.Config, logger logging.Logger, client *mongo.Client) *RuntimeRepoMongoDB {
	collection := client.Database(cfg.MongoDB.DBName).Collection("runtimes")
	return &RuntimeRepoMongoDB{
		cfg,
		logger,
		collection,
	}
}

func (r *RuntimeRepoMongoDB) Create(name string, userID string) (*entity.Runtime, error) {
	runtime := &entity.Runtime{
		ID:           primitive.NewObjectID().Hex(),
		Name:         name,
		CreationDate: time.Now().UTC(),
		Owner:        userID,
	}
	res, err := r.collection.InsertOne(context.Background(), runtime)
	if err != nil {
		return nil, err
	}

	runtime.ID = res.InsertedID.(string)
	return runtime, nil
}

func (r *RuntimeRepoMongoDB) FindAll() ([]entity.Runtime, error) {
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	var runtimes []entity.Runtime
	cur, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return runtimes, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var runtime entity.Runtime
		err = cur.Decode(&runtime)
		if err != nil {
			return runtimes, err
		}
		runtimes = append(runtimes, runtime)
	}

	return runtimes, nil
}
