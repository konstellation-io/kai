//go:build integration

package mongodb_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-logr/logr/testr"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/repository/mongodb"
	"github.com/konstellation-io/kai/engine/admin-api/domain/repository"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductRepositorySuite struct {
	suite.Suite
	mongoDBContainer   testcontainers.Container
	mongoClient        *mongo.Client
	productsCollection *mongo.Collection
	productRepo        *mongodb.ProductRepoMongoDB
}

func TestProductRepositorySuite(t *testing.T) {
	suite.Run(t, new(ProductRepositorySuite))
}

func (s *ProductRepositorySuite) SetupSuite() {
	ctx := context.Background()
	logger := testr.NewWithOptions(s.T(), testr.Options{Verbosity: -1})

	req := testcontainers.ContainerRequest{
		Image:        "mongo:latest",
		ExposedPorts: []string{"27017/tcp", "27018/tcp"},
		Env: map[string]string{
			"MONGO_INITDB_ROOT_USERNAME": "root",
			"MONGO_INITDB_ROOT_PASSWORD": "root",
		},
		WaitingFor: wait.ForLog("MongoDB starting"),
	}

	mongoDBContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	s.Require().NoError(err)

	host, err := mongoDBContainer.Host(context.Background())
	s.Require().NoError(err)
	p, err := mongoDBContainer.MappedPort(context.Background(), "27017/tcp")
	s.Require().NoError(err)

	port := p.Int()
	uri := fmt.Sprintf("mongodb://root:root@%v:%v/", host, port) //NOSONAR not used in secure contexts
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	s.Require().NoError(err)

	viper.Set(config.MongoDBKaiDatabaseKey, "kai")

	s.mongoDBContainer = mongoDBContainer
	s.productsCollection = client.Database(viper.GetString(config.MongoDBKaiDatabaseKey)).Collection("products")
	s.productRepo = mongodb.NewProductRepoMongoDB(logger, client)
	s.mongoClient = client

	s.Require().NoError(err)
}

func (s *ProductRepositorySuite) TearDownSuite() {
	s.Require().NoError(s.mongoDBContainer.Terminate(context.Background()))
}

func (s *ProductRepositorySuite) TearDownTest() {
	err := s.productsCollection.Drop(context.Background())
	s.Require().NoError(err)
}

func (s *ProductRepositorySuite) TestUpdate() {
	product := testhelpers.NewProductBuilder().WithPublishedVersion(testhelpers.StrPointer("test")).Build()
	ctx := context.Background()

	_, err := s.productsCollection.InsertOne(ctx, product)
	s.Require().NoError(err)

	expectedPublishedVersion := "another-version"
	product.UpdatePublishedVersion(expectedPublishedVersion)

	err = s.productRepo.Update(ctx, product)
	s.Require().NoError(err)

	actualProduct, err := s.productRepo.GetByID(ctx, product.ID)
	s.Require().NoError(err)

	s.Equal(product, actualProduct)
}

func (s *ProductRepositorySuite) TestUpdate_ProductDoesntExist() {
	product := testhelpers.NewProductBuilder().WithPublishedVersion(testhelpers.StrPointer("test")).Build()
	ctx := context.Background()

	err := s.productRepo.Update(ctx, product)
	s.Require().Error(err, mongodb.ErrUpdateProductNotFound)
}

func (s *ProductRepositorySuite) TestDeleteDatabase() {
	ctx := context.Background()
	name := "test-database"

	_, err := s.mongoClient.Database(name).Collection("test").InsertOne(ctx, bson.M{"foo": "bar"})
	s.Require().NoError(err)

	err = s.productRepo.DeleteDatabase(ctx, name)
	s.NoError(err)

	res, err := s.mongoClient.Database(name).Collection("test").CountDocuments(ctx, bson.M{})
	s.Require().NoError(err)
	s.Zero(res)
}

func (s *ProductRepositorySuite) TestCreate() {
	product := testhelpers.NewProductBuilder().Build()
	ctx := context.Background()

	createdProduct, err := s.productRepo.Create(ctx, product)
	s.Require().NoError(err)

	actualProduct, err := s.productRepo.GetByID(ctx, createdProduct.ID)
	s.Require().NoError(err)

	product.CreationDate = time.Time{}
	actualProduct.CreationDate = time.Time{}

	s.Equal(createdProduct, actualProduct)
}

func (s *ProductRepositorySuite) TestGetByID() {
	product := testhelpers.NewProductBuilder().Build()
	ctx := context.Background()

	_, err := s.productsCollection.InsertOne(ctx, product)
	s.Require().NoError(err)

	actualProduct, err := s.productRepo.GetByID(ctx, product.ID)
	s.Require().NoError(err)

	product.CreationDate = time.Time{}
	actualProduct.CreationDate = time.Time{}

	s.Equal(product, actualProduct)
}

func (s *ProductRepositorySuite) TestGetByID_ProductDoesntExist() {
	ctx := context.Background()

	_, err := s.productRepo.GetByID(ctx, "non-existing-id")
	s.ErrorIs(err, usecase.ErrProductNotFound)
}

func (s *ProductRepositorySuite) TestGetByName() {
	product := testhelpers.NewProductBuilder().Build()
	ctx := context.Background()

	_, err := s.productsCollection.InsertOne(ctx, product)
	s.Require().NoError(err)

	actualProduct, err := s.productRepo.GetByName(ctx, product.Name)
	s.Require().NoError(err)

	product.CreationDate = time.Time{}
	actualProduct.CreationDate = time.Time{}

	s.Equal(product, actualProduct)
}

func (s *ProductRepositorySuite) TestGetByName_ProductDoesntExist() {
	ctx := context.Background()

	_, err := s.productRepo.GetByName(ctx, "non-existing-name")
	s.ErrorIs(err, usecase.ErrProductNotFound)
}

func (s *ProductRepositorySuite) TestFindAll() {
	product1 := testhelpers.NewProductBuilder().Build()
	product2 := testhelpers.NewProductBuilder().WithName("product2").Build()
	ctx := context.Background()

	_, err := s.productsCollection.InsertMany(ctx, []interface{}{product1, product2})
	s.Require().NoError(err)

	products, err := s.productRepo.FindAll(ctx, nil)
	s.Require().NoError(err)

	s.Len(products, 2)
}

func (s *ProductRepositorySuite) TestFindAllWithNameFilter() {
	product1 := testhelpers.NewProductBuilder().Build()
	product2 := testhelpers.NewProductBuilder().WithName("product2").Build()
	ctx := context.Background()

	_, err := s.productsCollection.InsertMany(ctx, []interface{}{product1, product2})
	s.Require().NoError(err)

	filter := repository.FindAllFilter{ProductName: product2.Name}

	products, err := s.productRepo.FindAll(ctx, &filter)
	s.Require().NoError(err)

	s.Require().Len(products, 1)

	product2.CreationDate = time.Time{}
	products[0].CreationDate = time.Time{}

	s.Equal(product2, products[0])
}

func (s *ProductRepositorySuite) TestFindAll_Empty() {
	ctx := context.Background()

	products, err := s.productRepo.FindAll(ctx, nil)
	s.Require().NoError(err)

	s.Len(products, 0)
}

func (s *ProductRepositorySuite) TestFindByIDs() {
	product1 := testhelpers.NewProductBuilder().Build()
	product2 := testhelpers.NewProductBuilder().WithName("product2").Build()
	ctx := context.Background()

	_, err := s.productsCollection.InsertMany(ctx, []interface{}{product1, product2})
	s.Require().NoError(err)

	products, err := s.productRepo.FindByIDs(ctx, []string{product1.ID, product2.ID}, nil)
	s.Require().NoError(err)

	s.Len(products, 2)
}

func (s *ProductRepositorySuite) TestFindByIDs_Empty() {
	ctx := context.Background()

	products, err := s.productRepo.FindByIDs(ctx, []string{"non-existing-id"}, nil)
	s.Require().NoError(err)

	s.Len(products, 0)
}
