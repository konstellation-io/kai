package mongodb

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

type UserActivityRepoMongoDB struct {
	logger     logr.Logger
	collection *mongo.Collection
}

func NewUserActivityRepoMongoDB(logger logr.Logger, client *mongo.Client) *UserActivityRepoMongoDB {
	collection := client.Database(viper.GetString(config.MongoDBKaiDatabaseKey)).Collection("userActivity")

	return &UserActivityRepoMongoDB{
		logger,
		collection,
	}
}

//nolint:nestif // legacy code
func (r *UserActivityRepoMongoDB) Get(
	ctx context.Context,
	userEmail *string,
	types []entity.UserActivityType,
	versionIds []string,
	fromDate *string,
	toDate *string,
	lastID *string,
) ([]*entity.UserActivity, error) {
	const limit = 30

	filter := bson.M{}
	if lastID != nil {
		filter["_id"] = bson.M{"$lt": lastID}
	}

	if len(types) > 0 {
		filter["type"] = bson.M{"$in": types}
	}

	if userEmail != nil {
		filter["userId"] = userEmail
	}

	if len(versionIds) > 0 {
		filter["vars.value"] = bson.M{"$in": versionIds}
	}

	if fromDate != nil || toDate != nil {
		filterDate := bson.M{}

		if fromDate != nil {
			from, err := time.Parse(time.RFC3339, *fromDate)
			if err != nil {
				return nil, err
			}

			filterDate["$gte"] = from
		}

		if toDate != nil {
			to, err := time.Parse(time.RFC3339, *toDate)
			if err != nil {
				return nil, err
			}

			to = time.Date(to.Year(), to.Month(), to.Day(), 23, 59, 59, 999999999, to.Location())
			filterDate["$lte"] = to
		}

		filter["date"] = filterDate
	}

	var activities []*entity.UserActivity

	opts := options.Find().SetSort(bson.M{"_id": -1}).SetLimit(limit)

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return activities, err
	}

	err = cursor.All(ctx, &activities)
	if err != nil {
		return nil, err
	}

	return activities, nil
}

func (r *UserActivityRepoMongoDB) Create(activity entity.UserActivity) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, activity)
	if err != nil {
		return err
	}

	return nil
}
