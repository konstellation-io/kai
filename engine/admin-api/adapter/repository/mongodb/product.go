package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/repository"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase"
)

const (
	_productRepoTimeout = 60 * time.Second
)

var (
	ErrUpdateProductNotFound = errors.New("the product to be updated does not exist")
)

type ProductRepoMongoDB struct {
	logger     logr.Logger
	collection *mongo.Collection
	client     *mongo.Client
}

func NewProductRepoMongoDB(logger logr.Logger, client *mongo.Client) *ProductRepoMongoDB {
	collection := client.Database(viper.GetString(config.MongoDBKaiDatabaseKey)).Collection("products")

	productRepo := &ProductRepoMongoDB{
		logger,
		collection,
		client,
	}

	productRepo.createIndexes()

	return productRepo
}

func (r *ProductRepoMongoDB) createIndexes() {
	_, err := r.collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.M{
			"name": 1,
		},
	})
	if err != nil {
		r.logger.Error(err, "Error creating products collection indexes")
	}
}

func (r *ProductRepoMongoDB) Create(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	product.CreationDate = time.Now().UTC()

	_, err := r.collection.InsertOne(ctx, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *ProductRepoMongoDB) GetByID(ctx context.Context, productID string) (*entity.Product, error) {
	product := &entity.Product{}
	filter := bson.M{"_id": productID}

	err := r.collection.FindOne(ctx, filter).Decode(product)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, usecase.ErrProductNotFound
	}

	return product, err
}

func (r *ProductRepoMongoDB) GetByName(ctx context.Context, name string) (*entity.Product, error) {
	product := &entity.Product{}
	filter := bson.M{"name": name}

	err := r.collection.FindOne(ctx, filter).Decode(product)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, usecase.ErrProductNotFound
	}

	return product, err
}

func (r *ProductRepoMongoDB) FindAll(ctx context.Context, filter *repository.FindAllFilter) ([]*entity.Product, error) {
	var products []*entity.Product

	queryFilter := r.getFindAllMongoFilter(filter)

	cursor, err := r.collection.Find(ctx, queryFilter)
	if err != nil {
		return products, err
	}

	err = cursor.All(ctx, &products)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepoMongoDB) FindByIDs(ctx context.Context, ids []string, filter *repository.FindAllFilter) ([]*entity.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, _productRepoTimeout)
	defer cancel()

	queryFilter := bson.M{"_id": bson.M{"$in": ids}}
	if filter != nil && filter.ProductName != "" {
		queryFilter["name"] = filter.ProductName
	}

	cursor, err := r.collection.Find(ctx, queryFilter)
	if err != nil {
		return nil, err
	}

	var products []*entity.Product

	err = cursor.All(ctx, &products)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepoMongoDB) Update(ctx context.Context, product *entity.Product) error {
	ctx, cancel := context.WithTimeout(ctx, _productRepoTimeout)
	defer cancel()

	resp, err := r.collection.ReplaceOne(
		ctx,
		bson.M{"_id": product.ID},
		product,
	)
	if err != nil {
		return err
	}

	if resp.ModifiedCount == 0 {
		return ErrUpdateProductNotFound
	}

	return nil
}

func (r *ProductRepoMongoDB) DeleteDatabase(ctx context.Context, name string) error {
	err := r.client.Database(name).Drop(ctx)
	if err != nil {
		return fmt.Errorf("deleting product's %q database: %w", name, err)
	}

	return nil
}

func (r *ProductRepoMongoDB) Delete(ctx context.Context, productID string) error {
	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": productID})
	if res.DeletedCount == 0 {
		return usecase.ErrProductNotFound
	}

	return err
}

func (r *ProductRepoMongoDB) getFindAllMongoFilter(findAllFilter *repository.FindAllFilter) bson.M {
	filter := make(bson.M, 1)

	if findAllFilter == nil {
		return filter
	}

	if findAllFilter.ProductName != "" {
		filter["name"] = findAllFilter.ProductName
	}

	return filter
}
