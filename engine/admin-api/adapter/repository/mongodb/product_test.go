//go:build integration

package mongodb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-logr/logr/testr"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/repository/mongodb"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductRepositorySuite struct {
	suite.Suite
	mongoDBContainer   testcontainers.Container
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
