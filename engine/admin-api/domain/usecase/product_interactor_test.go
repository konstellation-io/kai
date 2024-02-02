//go:build unit

package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase"
	"github.com/konstellation-io/kai/engine/admin-api/mocks"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
	"github.com/sethvargo/go-password/password"
	"github.com/stretchr/testify/suite"
)

const (
	_testPassword     = "test-password"
	_testBucketPolicy = "test-policy"
	_testKvStore      = "test-kv-store"
)

type productSuite struct {
	suite.Suite

	ctrl              *gomock.Controller
	productInteractor *usecase.ProductInteractor

	logger            logr.Logger
	productRepo       *mocks.MockProductRepo
	versionRepo       *mocks.MockVersionRepo
	processRepo       *mocks.MockProcessRepository
	userActivityRepo  *mocks.MockUserActivityRepo
	accessControl     *mocks.MockAccessControl
	objectStorage     *mocks.MockObjectStorage
	userRegistry      *mocks.MockUserRegistry
	passwordGenerator password.PasswordGenerator
	natsService       *mocks.MockNatsManagerService
	predictionRepo    *mocks.MockPredictionRepo
}

func TestProductSuite(t *testing.T) {
	suite.Run(t, new(productSuite))
}

func (s *productSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())

	s.logger = testr.NewWithOptions(s.T(), testr.Options{Verbosity: -1})
	s.productRepo = mocks.NewMockProductRepo(ctrl)
	s.userActivityRepo = mocks.NewMockUserActivityRepo(ctrl)
	s.versionRepo = mocks.NewMockVersionRepo(ctrl)
	s.processRepo = mocks.NewMockProcessRepository(ctrl)
	s.accessControl = mocks.NewMockAccessControl(ctrl)
	s.objectStorage = mocks.NewMockObjectStorage(s.T())
	s.userRegistry = mocks.NewMockUserRegistry(ctrl)
	s.passwordGenerator = password.NewMockGenerator(_testPassword, nil)
	s.natsService = mocks.NewMockNatsManagerService(ctrl)
	s.predictionRepo = mocks.NewMockPredictionRepo(s.T())

	userActivity := usecase.NewUserActivityInteractor(
		s.logger,
		s.userActivityRepo,
		s.accessControl,
	)

	productInteractorOpts := usecase.ProductInteractorOpts{
		Logger:               s.logger,
		ProductRepo:          s.productRepo,
		VersionRepo:          s.versionRepo,
		ProcessRepo:          s.processRepo,
		UserActivity:         userActivity,
		AccessControl:        s.accessControl,
		ObjectStorage:        s.objectStorage,
		UserRegistry:         s.userRegistry,
		PasswordGenerator:    s.passwordGenerator,
		NatsService:          s.natsService,
		PredictionRepository: s.predictionRepo,
	}
	s.productInteractor = usecase.NewProductInteractor(&productInteractorOpts)
}

func (s *productSuite) TestGet() {
	productID := "product1"
	expectedProduct := &entity.Product{
		ID: productID,
	}

	user := testhelpers.NewUserBuilder().Build()

	ctx := context.Background()

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActViewProduct)
	s.productRepo.EXPECT().Get(ctx).Return(expectedProduct, nil)

	product, err := s.productInteractor.Get(ctx, user, productID)
	s.Require().Nil(err)
	s.Require().Equal(expectedProduct, product)
}

func (s *productSuite) TestCreateProduct() {
	ctx := context.Background()

	user := testhelpers.NewUserBuilder().Build()
	productID := "test-product"
	productName := "test-product"
	productDescription := "This is a product description"
	expectedProduct := &entity.Product{
		ID:           productID,
		Name:         productName,
		Description:  productDescription,
		CreationDate: time.Time{},
		Owner:        user.ID,
		MinioConfiguration: entity.MinioConfiguration{
			Bucket: productID,
		},
		ServiceAccount: entity.ServiceAccount{
			Username: productID,
			Group:    productID,
			Password: _testPassword,
		},
		KeyValueStore: _testKvStore,
	}

	s.accessControl.EXPECT().CheckRoleGrants(user, auth.ActCreateProduct).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, productID).Return(nil, usecase.ErrProductNotFound)
	s.productRepo.EXPECT().GetByName(ctx, productName).Return(nil, usecase.ErrProductNotFound)
	s.natsService.EXPECT().CreateGlobalKeyValueStore(ctx, productID).Return(_testKvStore, nil)
	s.objectStorage.EXPECT().CreateBucket(ctx, productID).Times(1).Return(nil)
	s.objectStorage.EXPECT().CreateBucketPolicy(ctx, productID).Times(1).Return(_testBucketPolicy, nil)
	s.userRegistry.EXPECT().CreateGroupWithPolicy(ctx, productID, _testBucketPolicy).Times(1).Return(nil)
	s.userRegistry.EXPECT().CreateUserWithinGroup(ctx, productID, _testPassword, productID).Times(1).Return(nil)
	s.versionRepo.EXPECT().CreateIndexes(ctx, productID).Return(nil)
	s.processRepo.EXPECT().CreateIndexes(ctx, productID).Return(nil)
	s.predictionRepo.EXPECT().
		CreateUser(ctx, productID, expectedProduct.ServiceAccount.Username, expectedProduct.ServiceAccount.Password).
		Return(nil)
	s.productRepo.EXPECT().Create(ctx, expectedProduct).Return(expectedProduct, nil)

	product, err := s.productInteractor.CreateProduct(ctx, user, productName, productDescription)

	s.Require().Nil(err)
	s.Require().Equal(expectedProduct, product)
}

func (s *productSuite) TestCreateProduct_FailsIfUserHasNotPermission() {
	ctx := context.Background()

	user := testhelpers.NewUserBuilder().Build()
	productName := "product-name"
	productDescription := "This is a product description"

	grantError := errors.New("grant error")

	s.accessControl.EXPECT().CheckRoleGrants(user, auth.ActCreateProduct).Return(grantError)

	product, err := s.productInteractor.CreateProduct(ctx, user, productName, productDescription)

	s.Require().Error(grantError, err)
	s.Require().Nil(product)
}

func (s *productSuite) TestCreateProduct_FailsIfProductHasAnInvalidField() {
	ctx := context.Background()

	user := testhelpers.NewUserBuilder().Build()
	// the product name is bigger thant the max length (it should be lte=40)
	productName := "lore ipsum dolor sit amet consectetur adipiscing elit sed do eiusmod tempor incididunt ut labores"
	productDescription := "This is a product description"

	s.accessControl.EXPECT().CheckRoleGrants(user, auth.ActCreateProduct).Return(nil)

	product, err := s.productInteractor.CreateProduct(ctx, user, productName, productDescription)

	s.Require().Error(err)
	s.Require().Nil(product)
}

func (s *productSuite) TestCreateProduct_FailsIfProductWithSameIDAlreadyExists() {
	ctx := context.Background()

	user := testhelpers.NewUserBuilder().Build()

	productID := "test-product"
	productName := "test-product"
	productDescription := "This is a product description"

	existingProduct := &entity.Product{
		ID:          productID,
		Name:        "existing-product-name",
		Description: "existing-product-description",
	}

	s.accessControl.EXPECT().CheckRoleGrants(user, auth.ActCreateProduct).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, productID).Return(existingProduct, nil)

	product, err := s.productInteractor.CreateProduct(ctx, user, productName, productDescription)

	s.Require().Error(err)
	s.Require().Nil(product)
}

func (s *productSuite) TestCreateProduct_FailsIfProductWithSameNameAlreadyExists() {
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()

	productID := "test-product"
	productName := "test product"

	productDescription := "This is a product description"

	existingProduct := &entity.Product{
		ID:          "existing-product-id",
		Name:        productName,
		Description: "existing-product-description",
	}

	s.accessControl.EXPECT().CheckRoleGrants(user, auth.ActCreateProduct).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, productID).Return(nil, usecase.ErrProductNotFound)
	s.productRepo.EXPECT().GetByName(ctx, productName).Return(existingProduct, nil)

	product, err := s.productInteractor.CreateProduct(ctx, user, productName, productDescription)

	s.Require().Error(err)
	s.Require().Nil(product)
}

func (s *productSuite) TestCreateProduct_ErrorCheckingProductIDInRepo() {
	ctx := context.Background()

	user := testhelpers.NewUserBuilder().Build()
	productID := "test-product"
	productName := "test-product"
	productDescription := "This is a product description"
	repoError := errors.New("repo error")

	s.accessControl.EXPECT().CheckRoleGrants(user, auth.ActCreateProduct).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, productID).Return(nil, repoError)

	_, err := s.productInteractor.CreateProduct(ctx, user, productName, productDescription)
	s.Require().ErrorIs(err, repoError)
}

func (s *productSuite) TestCreateProduct_ErrorCheckingProductNameInRepo() {
	ctx := context.Background()

	user := testhelpers.NewUserBuilder().Build()
	productID := "test-product"
	productName := "test-product"
	productDescription := "This is a product description"
	repoError := errors.New("repo error")

	s.accessControl.EXPECT().CheckRoleGrants(user, auth.ActCreateProduct).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, productID).Return(nil, usecase.ErrProductNotFound)
	s.productRepo.EXPECT().GetByName(ctx, productName).Return(nil, repoError)

	_, err := s.productInteractor.CreateProduct(ctx, user, productName, productDescription)
	s.Require().ErrorIs(err, repoError)
}

func (s *productSuite) TestCreateProduct_ErrorCreatingVersionRepoIndexes() {
	ctx := context.Background()

	user := testhelpers.NewUserBuilder().Build()
	productID := "test-product"
	productName := "test-product"
	productDescription := "This is a product description"

	expectedError := errors.New("error creating versions collection indexes")

	s.accessControl.EXPECT().CheckRoleGrants(user, auth.ActCreateProduct).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, productID).Return(nil, usecase.ErrProductNotFound)
	s.productRepo.EXPECT().GetByName(ctx, productName).Return(nil, usecase.ErrProductNotFound)
	s.natsService.EXPECT().CreateGlobalKeyValueStore(ctx, productID).Return(_testKvStore, nil)

	s.objectStorage.EXPECT().CreateBucket(ctx, productID).Return(nil)
	s.objectStorage.EXPECT().CreateBucketPolicy(ctx, productID).Times(1).Return(_testBucketPolicy, nil)
	s.userRegistry.EXPECT().CreateGroupWithPolicy(ctx, productID, _testBucketPolicy).Times(1).Return(nil)
	s.userRegistry.EXPECT().CreateUserWithinGroup(ctx, productID, _testPassword, productID).Times(1).Return(nil)

	s.versionRepo.EXPECT().CreateIndexes(ctx, productID).Return(expectedError)

	_, err := s.productInteractor.CreateProduct(ctx, user, productName, productDescription)
	s.Require().ErrorIs(err, expectedError)
}

func (s *productSuite) TestCreateProduct_ErrorCreatingProcessRegistryRepoIndexes() {
	ctx := context.Background()

	user := testhelpers.NewUserBuilder().Build()
	productID := "test-product"
	productName := "test-product"
	productDescription := "This is a product description"

	expectedError := errors.New("error creating versions collection indexes")

	s.accessControl.EXPECT().CheckRoleGrants(user, auth.ActCreateProduct).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, productID).Return(nil, usecase.ErrProductNotFound)
	s.productRepo.EXPECT().GetByName(ctx, productName).Return(nil, usecase.ErrProductNotFound)
	s.natsService.EXPECT().CreateGlobalKeyValueStore(ctx, productID).Return(_testKvStore, nil)
	s.objectStorage.EXPECT().CreateBucket(ctx, productID).Return(nil)
	s.objectStorage.EXPECT().CreateBucketPolicy(ctx, productID).Times(1).Return(_testBucketPolicy, nil)
	s.userRegistry.EXPECT().CreateGroupWithPolicy(ctx, productID, _testBucketPolicy).Times(1).Return(nil)
	s.userRegistry.EXPECT().CreateUserWithinGroup(ctx, productID, _testPassword, productID).Times(1).Return(nil)

	s.versionRepo.EXPECT().CreateIndexes(ctx, productID).Return(nil)
	s.processRepo.EXPECT().CreateIndexes(ctx, productID).Return(expectedError)

	_, err := s.productInteractor.CreateProduct(ctx, user, productName, productDescription)
	s.Require().ErrorIs(err, expectedError)
}

func (s *productSuite) TestCreateProduct_FailsIfCreateProductFails() {
	ctx := context.Background()

	user := testhelpers.NewUserBuilder().Build()

	productID := "test-product"
	productName := "test-product"
	productDescription := "This is a product description"

	newProduct := &entity.Product{
		ID:           productID,
		Name:         productName,
		Description:  productDescription,
		Owner:        user.ID,
		CreationDate: time.Time{},
		MinioConfiguration: entity.MinioConfiguration{
			Bucket: productID,
		},
		ServiceAccount: entity.ServiceAccount{
			Username: productID,
			Group:    productID,
			Password: _testPassword,
		},
		KeyValueStore: _testKvStore,
	}

	expectedError := errors.New("create product error")

	s.accessControl.EXPECT().CheckRoleGrants(user, auth.ActCreateProduct).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, productID).Return(nil, usecase.ErrProductNotFound)
	s.productRepo.EXPECT().GetByName(ctx, productName).Return(nil, usecase.ErrProductNotFound)
	s.natsService.EXPECT().CreateGlobalKeyValueStore(ctx, productID).Return(_testKvStore, nil)

	s.objectStorage.EXPECT().CreateBucket(ctx, productID).Return(nil)
	s.objectStorage.EXPECT().CreateBucketPolicy(ctx, productID).Times(1).Return(_testBucketPolicy, nil)
	s.userRegistry.EXPECT().CreateGroupWithPolicy(ctx, productID, _testBucketPolicy).Times(1).Return(nil)
	s.userRegistry.EXPECT().CreateUserWithinGroup(ctx, productID, _testPassword, productID).Times(1).Return(nil)

	s.versionRepo.EXPECT().CreateIndexes(ctx, productID).Return(nil)
	s.processRepo.EXPECT().CreateIndexes(ctx, productID).Return(nil)
	s.productRepo.EXPECT().Create(ctx, newProduct).Return(nil, expectedError)

	_, err := s.productInteractor.CreateProduct(ctx, user, productName, productDescription)
	s.Require().ErrorIs(err, expectedError)
}

func (s *productSuite) TestGetByID() {
	ctx := context.Background()

	user := testhelpers.NewUserBuilder().Build()
	productID := "test-product"
	productName := "test-product"

	expected := &entity.Product{
		ID:           productID,
		Name:         productName,
		Description:  "Product description...",
		CreationDate: time.Time{},
	}

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActViewProduct).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, productID).Return(expected, nil)

	actual, err := s.productInteractor.GetByID(ctx, user, productID)

	s.Require().Nil(err)
	s.Require().Equal(expected, actual)
}

func (s *productSuite) TestFindAll() {
	ctx := context.Background()

	productID := "test-product"
	productName := "test-product"

	user := testhelpers.NewUserBuilder().
		WithProductGrants(
			map[string][]string{
				productID: {
					auth.ActViewProduct.String(),
				},
			},
		).Build()

	expected := []*entity.Product{
		{
			ID:           productID,
			Name:         productName,
			Description:  "Product description...",
			CreationDate: time.Time{},
		},
	}

	userProducts := []string{productID}

	s.accessControl.EXPECT().IsAdmin(user).Return(false)
	s.accessControl.EXPECT().GetUserProducts(user).Return(userProducts)
	s.productRepo.EXPECT().FindByIDs(ctx, userProducts).Return(expected, nil)

	actual, err := s.productInteractor.FindAll(ctx, user)

	s.Require().Nil(err)
	s.Require().Equal(expected, actual)
}

func (s *productSuite) TestFindAll_AdminUser() {
	ctx := context.Background()

	productID := "test-product"
	productName := "test-product"

	user := testhelpers.NewUserBuilder().WithRoles([]string{auth.DefaultAdminRole}).Build()

	expected := []*entity.Product{
		{
			ID:           productID,
			Name:         productName,
			Description:  "Product description...",
			CreationDate: time.Time{},
		},
	}

	s.accessControl.EXPECT().IsAdmin(user).Return(true)
	s.productRepo.EXPECT().FindAll(ctx).Return(expected, nil)

	actual, err := s.productInteractor.FindAll(ctx, user)

	s.Require().Nil(err)
	s.Require().Equal(expected, actual)
}
