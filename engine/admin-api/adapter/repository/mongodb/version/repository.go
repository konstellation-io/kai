package version

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	usecaseErrors "github.com/konstellation-io/kai/engine/admin-api/domain/usecase/errors"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/logging"
)

const versionsCollectionName = "versions"

type VersionRepoMongoDB struct {
	cfg    *config.Config
	logger logging.Logger
	client *mongo.Client
}

func NewVersionRepoMongoDB(
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

	indexes := []mongo.IndexModel{
		{
			Keys: bson.M{
				"name": 1,
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

	versionDTO.ID = primitive.NewObjectID().Hex()
	versionDTO.CreationDate = time.Now().UTC()
	versionDTO.CreationAuthor = userID
	versionDTO.Status = entity.VersionStatusCreating.String()

	res, err := collection.InsertOne(context.Background(), versionDTO)
	if err != nil {
		return nil, err
	}

	versionDTO.ID = res.InsertedID.(string)

	savedVersion := mapDTOToEntity(versionDTO)

	return savedVersion, nil
}

func (r *VersionRepoMongoDB) GetByID(productID, versionID string) (*entity.Version, error) {
	collection := r.client.Database(productID).Collection(versionsCollectionName)

	versionDTO := &versionDTO{}
	filter := bson.M{"_id": versionID}

	err := collection.FindOne(context.Background(), filter).Decode(versionDTO)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, usecaseErrors.ErrVersionNotFound
	}

	return mapDTOToEntity(versionDTO), err
}

func (r *VersionRepoMongoDB) GetByName(ctx context.Context, productID, name string) (*entity.Version, error) {
	collection := r.client.Database(productID).Collection(versionsCollectionName)

	versionDTO := &versionDTO{}
	filter := bson.M{"name": name}

	err := collection.FindOne(ctx, filter).Decode(versionDTO)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, usecaseErrors.ErrVersionNotFound
	}

	return mapDTOToEntity(versionDTO), err
}

func (r *VersionRepoMongoDB) Update(productID string, version *entity.Version) error {
	collection := r.client.Database(productID).Collection(versionsCollectionName)

	versionDTO := mapEntityToDTO(version)
	_, err := collection.ReplaceOne(context.Background(), bson.M{"_id": version.ID}, versionDTO)

	return err
}

func (r *VersionRepoMongoDB) GetVersionsByProduct(ctx context.Context, productID string) ([]*entity.Version, error) {
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
	versionID string,
	status entity.VersionStatus,
) error {
	collection := r.client.Database(productID).Collection(versionsCollectionName)

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": versionID},
		bson.M{"$set": bson.M{"status": status}},
	)
	if err != nil {
		return err
	}

	if result.ModifiedCount != 1 {
		return usecaseErrors.ErrVersionNotFound
	}

	return nil
}

func (r *VersionRepoMongoDB) SetErrors(
	ctx context.Context,
	productID string,
	version *entity.Version,
	errorMessages []string,
) (*entity.Version, error) {
	collection := r.client.Database(productID).Collection(versionsCollectionName)

	versionDTO := mapEntityToDTO(version)

	versionDTO.Status = entity.VersionStatusError.String()
	versionDTO.Errors = errorMessages

	elem := bson.M{"$set": bson.M{"status": versionDTO.Status, "errors": versionDTO.Errors}}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": versionDTO.ID}, elem)
	if err != nil {
		return nil, err
	}

	if result.ModifiedCount != 1 {
		return nil, usecaseErrors.ErrVersionNotFound
	}

	return mapDTOToEntity(versionDTO), nil
}

// TODO: Rethink this to upload krt YAML files.
//
//nolint:godox // To be done.
func (r *VersionRepoMongoDB) UploadKRTFile(productID string, version *entity.Version, file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("reading KRT file at %s: %w", file, err)
	}

	bucket, err := gridfs.NewBucket(
		r.client.Database(productID),
		options.GridFSBucket().SetName(r.cfg.MongoDB.KRTBucket),
	)
	if err != nil {
		return fmt.Errorf("creating bucket %q to store KRT files: %w", r.cfg.MongoDB.DBName, err)
	}

	versionDTO := mapEntityToDTO(version)

	filename := fmt.Sprintf("%s.krt", versionDTO.Name)

	uploadStream, err := bucket.OpenUploadStreamWithID(
		versionDTO.ID,
		filename,
	)
	if err != nil {
		return fmt.Errorf("opening KRT upload stream: %w", err)
	}
	defer uploadStream.Close()

	fileSize, err := uploadStream.Write(data)
	if err != nil {
		return fmt.Errorf("writing into the KRT upload stream: %w", err)
	}

	r.logger.Infof("Uploaded %d bytes of %q to GridFS successfully", filename, fileSize)

	return nil
}

func (r *VersionRepoMongoDB) ClearPublishedVersion(ctx context.Context, productID string) (*entity.Version, error) {
	collection := r.client.Database(productID).Collection(versionsCollectionName)

	oldPublishedVersion := &versionDTO{}

	filter := bson.M{"status": entity.VersionStatusPublished}

	upd := bson.M{
		"$set": bson.M{
			"status":            entity.VersionStatusStarted,
			"publicationDate":   nil,
			"publicationAuthor": nil,
		},
	}

	upsert := true
	after := options.After
	opt := &options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	result := collection.FindOneAndUpdate(ctx, filter, upd, opt)
	err := result.Decode(oldPublishedVersion)

	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	return mapDTOToEntity(oldPublishedVersion), nil
}

// TODO: Delete version method.
//
//nolint:godox // To be done.
