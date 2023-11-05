package versionrepository

import (
	"context"
	"errors"
	"time"

	"github.com/konstellation-io/kai/engine/admin-api/domain/service/logging"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/version"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

const versionsCollectionName = "versions"

type VersionRepoMongoDB struct {
	cfg    *config.Config
	logger logging.Logger
	client *mongo.Client
}

func New(
	cfg *config.Config,
	logger logging.Logger,
	client *mongo.Client,
) *VersionRepoMongoDB {
	versions := &VersionRepoMongoDB{
		cfg,
		logger,
		client,
	}

	return versions
}

func (r *VersionRepoMongoDB) CreateIndexes(ctx context.Context, productID string) error {
	collection := r.client.Database(productID).Collection(versionsCollectionName)
	r.logger.Infof("MongoDB creating indexes for %s collection...", versionsCollectionName)

	indexes := []mongo.IndexModel{
		{
			Keys: bson.M{
				"tag": 1,
			},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)

	return err
}

func (r *VersionRepoMongoDB) Create(userID, productID string, newVersion *entity.Version) (*entity.Version, error) {
	collection := r.client.Database(productID).Collection(versionsCollectionName)

	versionDTO := mapEntityToDTO(newVersion)

	versionDTO.CreationDate = time.Now().UTC()
	versionDTO.CreationAuthor = userID
	versionDTO.Status = entity.VersionStatusCreated.String()

	_, err := collection.InsertOne(context.Background(), versionDTO)
	if err != nil {
		return nil, err
	}

	savedVersion := mapDTOToEntity(versionDTO)

	return savedVersion, nil
}

func (r *VersionRepoMongoDB) GetByTag(ctx context.Context, productID, tag string) (*entity.Version, error) {
	collection := r.client.Database(productID).Collection(versionsCollectionName)

	versionDTO := &versionDTO{}
	filter := bson.M{"tag": tag}

	err := collection.FindOne(ctx, filter).Decode(versionDTO)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, version.ErrVersionNotFound
	}

	return mapDTOToEntity(versionDTO), err
}

func (r *VersionRepoMongoDB) Update(productID string, updatedVersion *entity.Version) error {
	collection := r.client.Database(productID).Collection(versionsCollectionName)

	versionDTO := mapEntityToDTO(updatedVersion)
	updateResult, err := collection.ReplaceOne(context.Background(), bson.M{"tag": updatedVersion.Tag}, versionDTO)

	if updateResult.ModifiedCount == 0 {
		return version.ErrVersionNotFound
	}

	return err
}

func (r *VersionRepoMongoDB) ListVersionsByProduct(ctx context.Context, productID string) ([]*entity.Version, error) {
	collection := r.client.Database(productID).Collection(versionsCollectionName)

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var versions []*entity.Version

	cur, err := collection.Find(ctxWithTimeout, bson.M{})
	if err != nil {
		return versions, err
	}
	defer cur.Close(ctxWithTimeout)

	for cur.Next(ctxWithTimeout) {
		var versionDTO versionDTO

		err = cur.Decode(&versionDTO)
		if err != nil {
			return versions, err
		}

		versions = append(versions, mapDTOToEntity(&versionDTO))
	}

	return versions, nil
}

func (r *VersionRepoMongoDB) SetStatus(
	ctx context.Context,
	productID string,
	versionTag string,
	status entity.VersionStatus,
) error {
	collection := r.client.Database(productID).Collection(versionsCollectionName)

	res := collection.FindOneAndUpdate(
		ctx,
		bson.M{"tag": versionTag},
		bson.M{"$set": bson.M{"status": status, "error": ""}},
	)
	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		return version.ErrVersionNotFound
	}

	return res.Err()
}

func (r *VersionRepoMongoDB) SetCriticalStatusWithError(ctx context.Context, productID, versionTag, errorMessage string) error {
	return r.setStatusWithError(ctx, productID, versionTag, errorMessage, entity.VersionStatusCritical)
}

func (r *VersionRepoMongoDB) SetErrorStatusWithError(ctx context.Context, productID, versionTag, errorMessage string) error {
	return r.setStatusWithError(ctx, productID, versionTag, errorMessage, entity.VersionStatusError)
}

func (r *VersionRepoMongoDB) setStatusWithError(
	ctx context.Context,
	productID, versionTag, errorMessage string,
	status entity.VersionStatus,
) error {
	collection := r.client.Database(productID).Collection(versionsCollectionName)

	elem := bson.M{"$set": bson.M{"status": status.String(), "error": errorMessage}}

	result, err := collection.UpdateOne(ctx, bson.M{"tag": versionTag}, elem)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return version.ErrVersionNotFound
	}

	return nil
}
