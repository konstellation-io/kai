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
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/auth"
	"github.com/konstellation-io/kai/engine/admin-api/mocks"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
	"github.com/stretchr/testify/require"
)

type productSuite struct {
	ctrl              *gomock.Controller
	productInteractor *usecase.ProductInteractor
	mocks             *productSuiteMocks
}

type productSuiteMocks struct {
	logger              logr.Logger
	productRepo         *mocks.MockProductRepo
	measurementRepo     *mocks.MockMeasurementRepo
	versionRepo         *mocks.MockVersionRepo
	metricRepo          *mocks.MockMetricRepo
	processLogRepo      *mocks.MockProcessLogRepository
	processRegistryRepo *mocks.MockProcessRegistryRepo
	userActivityRepo    *mocks.MockUserActivityRepo
	accessControl       *mocks.MockAccessControl
}

func newProductSuite(t *testing.T) *productSuite {
	ctrl := gomock.NewController(t)

	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})
	productRepo := mocks.NewMockProductRepo(ctrl)
	userActivityRepo := mocks.NewMockUserActivityRepo(ctrl)
	measurementRepo := mocks.NewMockMeasurementRepo(ctrl)
	versionRepo := mocks.NewMockVersionRepo(ctrl)
	metricRepo := mocks.NewMockMetricRepo(ctrl)
	processLogRepo := mocks.NewMockProcessLogRepository(ctrl)
	processRegistryRepo := mocks.NewMockProcessRegistryRepo(ctrl)
	accessControl := mocks.NewMockAccessControl(ctrl)

	userActivity := usecase.NewUserActivityInteractor(
		logger,
		userActivityRepo,
		accessControl,
	)

	ps := usecase.ProductInteractorOpts{
		Logger:              logger,
		ProductRepo:         productRepo,
		MeasurementRepo:     measurementRepo,
		VersionRepo:         versionRepo,
		MetricRepo:          metricRepo,
		ProcessLogRepo:      processLogRepo,
		ProcessRegistryRepo: processRegistryRepo,
		UserActivity:        userActivity,
		AccessControl:       accessControl,
	}
	productInteractor := usecase.NewProductInteractor(&ps)

	return &productSuite{
		ctrl:              ctrl,
		productInteractor: productInteractor,
		mocks: &productSuiteMocks{
			logger,
			productRepo,
			measurementRepo,
			versionRepo,
			metricRepo,
			processLogRepo,
			processRegistryRepo,
			userActivityRepo,
			accessControl,
		},
	}
}

func TestGet(t *testing.T) {
	s := newProductSuite(t)
	defer s.ctrl.Finish()

	productID := "product1"
	expectedProduct := &entity.Product{
		ID: productID,
	}

	user := testhelpers.NewUserBuilder().Build()

	ctx := context.Background()

	s.mocks.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActViewProduct)
	s.mocks.productRepo.EXPECT().Get(ctx).Return(expectedProduct, nil)

	product, err := s.productInteractor.Get(ctx, user, productID)
	require.Nil(t, err)
	require.Equal(t, expectedProduct, product)
}

func TestCreateProduct(t *testing.T) {
	s := newProductSuite(t)
	defer s.ctrl.Finish()

	ctx := context.Background()

	user := testhelpers.NewUserBuilder().Build()
	productID := "product-id"
	productName := "product-name"
	productDescription := "This is a product description"
	expectedProduct := &entity.Product{
		ID:           productID,
		Name:         productName,
		Description:  productDescription,
		CreationDate: time.Time{},
		Owner:        user.ID,
	}

	s.mocks.accessControl.EXPECT().CheckRoleGrants(user, auth.ActCreateProduct).Return(nil)
	s.mocks.productRepo.EXPECT().GetByID(ctx, productID).Return(nil, usecase.ErrProductNotFound)
	s.mocks.productRepo.EXPECT().GetByName(ctx, productName).Return(nil, usecase.ErrProductNotFound)
	s.mocks.productRepo.EXPECT().Create(ctx, expectedProduct).Return(expectedProduct, nil)
	s.mocks.measurementRepo.EXPECT().CreateDatabase(productID).Return(nil)
	s.mocks.versionRepo.EXPECT().CreateIndexes(ctx, productID).Return(nil)
	s.mocks.metricRepo.EXPECT().CreateIndexes(ctx, productID).Return(nil)
	s.mocks.processLogRepo.EXPECT().CreateIndexes(ctx, productID).Return(nil)
	s.mocks.processRegistryRepo.EXPECT().CreateIndexes(ctx, productID).Return(nil)

	product, err := s.productInteractor.CreateProduct(ctx, user, productID, productName, productDescription)

	require.Nil(t, err)
	require.Equal(t, expectedProduct, product)
}

func TestCreateProduct_FailsIfUserHasNotPermission(t *testing.T) {
	s := newProductSuite(t)
	defer s.ctrl.Finish()

	ctx := context.Background()

	user := testhelpers.NewUserBuilder().Build()
	productID := "product-id"
	productName := "product-name"
	productDescription := "This is a product description"

	grantError := errors.New("grant error")

	s.mocks.accessControl.EXPECT().CheckRoleGrants(user, auth.ActCreateProduct).Return(grantError)

	product, err := s.productInteractor.CreateProduct(ctx, user, productID, productName, productDescription)

	require.Error(t, grantError, err)
	require.Nil(t, product)
}

func TestCreateProduct_FailsIfProductHasAnInvalidField(t *testing.T) {
	s := newProductSuite(t)
	defer s.ctrl.Finish()

	ctx := context.Background()

	user := testhelpers.NewUserBuilder().Build()
	productID := "product-id"
	// the product name is bigger thant the max length (it should be lte=40)
	productName := "lore ipsum dolor sit amet consectetur adipiscing elit sed do eiusmod tempor incididunt ut labores"
	productDescription := "This is a product description"

	s.mocks.accessControl.EXPECT().CheckRoleGrants(user, auth.ActCreateProduct).Return(nil)

	product, err := s.productInteractor.CreateProduct(ctx, user, productID, productName, productDescription)

	require.Error(t, err)
	require.Nil(t, product)
}

func TestCreateProduct_FailsIfProductWithSameIDAlreadyExists(t *testing.T) {
	s := newProductSuite(t)
	defer s.ctrl.Finish()

	ctx := context.Background()

	user := testhelpers.NewUserBuilder().Build()

	productID := "product-id"
	productName := "product-name"
	productDescription := "This is a product description"

	existingProduct := &entity.Product{
		ID:          productID,
		Name:        "existing-product-name",
		Description: "existing-product-description",
	}

	s.mocks.accessControl.EXPECT().CheckRoleGrants(user, auth.ActCreateProduct).Return(nil)
	s.mocks.productRepo.EXPECT().GetByID(ctx, productID).Return(existingProduct, nil)

	product, err := s.productInteractor.CreateProduct(ctx, user, productID, productName, productDescription)

	require.Error(t, err)
	require.Nil(t, product)
}

func TestCreateProduct_FailsIfProductWithSameNameAlreadyExists(t *testing.T) {
	s := newProductSuite(t)
	defer s.ctrl.Finish()

	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()

	productName := "product-name"
	productID := "new-product-id"
	productDescription := "This is a product description"

	existingProduct := &entity.Product{
		ID:          "existing-product-id",
		Name:        productName,
		Description: "existing-product-description",
	}

	s.mocks.accessControl.EXPECT().CheckRoleGrants(user, auth.ActCreateProduct).Return(nil)
	s.mocks.productRepo.EXPECT().GetByID(ctx, productID).Return(nil, usecase.ErrProductNotFound)
	s.mocks.productRepo.EXPECT().GetByName(ctx, productName).Return(existingProduct, nil)

	product, err := s.productInteractor.CreateProduct(ctx, user, productID, productName, productDescription)

	require.Error(t, err)
	require.Nil(t, product)
}

func TestCreateProduct_FailsIfCreateProductFails(t *testing.T) {
	s := newProductSuite(t)
	defer s.ctrl.Finish()

	ctx := context.Background()

	user := testhelpers.NewUserBuilder().Build()
	productName := "product-name"
	productID := "new-product-id"
	productDescription := "This is a product description"

	newProduct := &entity.Product{
		ID:           productID,
		Name:         productName,
		Description:  productDescription,
		Owner:        user.ID,
		CreationDate: time.Time{},
	}

	s.mocks.accessControl.EXPECT().CheckRoleGrants(user, auth.ActCreateProduct).Return(nil)
	s.mocks.productRepo.EXPECT().GetByID(ctx, productID).Return(nil, usecase.ErrProductNotFound)
	s.mocks.productRepo.EXPECT().GetByName(ctx, productName).Return(nil, usecase.ErrProductNotFound)
	s.mocks.productRepo.EXPECT().Create(ctx, newProduct).Return(nil, errors.New("create product error"))

	product, err := s.productInteractor.CreateProduct(ctx, user, productID, productName, productDescription)

	require.Error(t, err)
	require.Nil(t, product)
}

func TestCreateProduct_ErrorCheckingProductIDInRepo(t *testing.T) {
	s := newProductSuite(t)
	defer s.ctrl.Finish()

	ctx := context.Background()

	user := testhelpers.NewUserBuilder().Build()
	productID := "product-id"
	productName := "product-name"
	productDescription := "This is a product description"
	repoError := errors.New("repo error")

	s.mocks.accessControl.EXPECT().CheckRoleGrants(user, auth.ActCreateProduct).Return(nil)
	s.mocks.productRepo.EXPECT().GetByID(ctx, productID).Return(nil, repoError)

	_, err := s.productInteractor.CreateProduct(ctx, user, productID, productName, productDescription)
	require.ErrorIs(t, err, repoError)
}

func TestCreateProduct_ErrorCheckingProductNameInRepo(t *testing.T) {
	s := newProductSuite(t)
	defer s.ctrl.Finish()

	ctx := context.Background()

	user := testhelpers.NewUserBuilder().Build()
	productID := "product-id"
	productName := "product-name"
	productDescription := "This is a product description"
	repoError := errors.New("repo error")

	s.mocks.accessControl.EXPECT().CheckRoleGrants(user, auth.ActCreateProduct).Return(nil)
	s.mocks.productRepo.EXPECT().GetByID(ctx, productID).Return(nil, usecase.ErrProductNotFound)
	s.mocks.productRepo.EXPECT().GetByName(ctx, productName).Return(nil, repoError)

	_, err := s.productInteractor.CreateProduct(ctx, user, productID, productName, productDescription)
	require.ErrorIs(t, err, repoError)
}

func TestCreateProduct_ErrorCreatingDatabase(t *testing.T) {
	s := newProductSuite(t)
	defer s.ctrl.Finish()

	ctx := context.Background()

	user := testhelpers.NewUserBuilder().Build()
	productID := "product-id"
	productName := "product-name"
	productDescription := "This is a product description"
	expectedProduct := &entity.Product{
		ID:           productID,
		Name:         productName,
		Description:  productDescription,
		CreationDate: time.Time{},
		Owner:        user.ID,
	}

	expectedError := errors.New("error creating database")

	s.mocks.accessControl.EXPECT().CheckRoleGrants(user, auth.ActCreateProduct).Return(nil)
	s.mocks.productRepo.EXPECT().GetByID(ctx, productID).Return(nil, usecase.ErrProductNotFound)
	s.mocks.productRepo.EXPECT().GetByName(ctx, productName).Return(nil, usecase.ErrProductNotFound)
	s.mocks.productRepo.EXPECT().Create(ctx, expectedProduct).Return(expectedProduct, nil)
	s.mocks.measurementRepo.EXPECT().CreateDatabase(productID).Return(expectedError)

	_, err := s.productInteractor.CreateProduct(ctx, user, productID, productName, productDescription)
	require.ErrorIs(t, err, expectedError)
}

func TestCreateProduct_ErrorCreatingMetricsRepoIndexes(t *testing.T) {
	s := newProductSuite(t)
	defer s.ctrl.Finish()

	ctx := context.Background()

	user := testhelpers.NewUserBuilder().Build()
	productID := "product-id"
	productName := "product-name"
	productDescription := "This is a product description"
	expectedProduct := &entity.Product{
		ID:           productID,
		Name:         productName,
		Description:  productDescription,
		CreationDate: time.Time{},
		Owner:        user.ID,
	}

	expectedError := errors.New("error creating collection indexes")

	s.mocks.accessControl.EXPECT().CheckRoleGrants(user, auth.ActCreateProduct).Return(nil)
	s.mocks.productRepo.EXPECT().GetByID(ctx, productID).Return(nil, usecase.ErrProductNotFound)
	s.mocks.productRepo.EXPECT().GetByName(ctx, productName).Return(nil, usecase.ErrProductNotFound)
	s.mocks.productRepo.EXPECT().Create(ctx, expectedProduct).Return(expectedProduct, nil)
	s.mocks.measurementRepo.EXPECT().CreateDatabase(productID).Return(nil)
	s.mocks.metricRepo.EXPECT().CreateIndexes(ctx, productID).Return(expectedError)

	_, err := s.productInteractor.CreateProduct(ctx, user, productID, productName, productDescription)
	require.ErrorIs(t, err, expectedError)
}

func TestCreateProduct_ErrorCreatingProcessLogRepoIndexes(t *testing.T) {
	s := newProductSuite(t)
	defer s.ctrl.Finish()

	ctx := context.Background()

	user := testhelpers.NewUserBuilder().Build()
	productID := "product-id"
	productName := "product-name"
	productDescription := "This is a product description"
	expectedProduct := &entity.Product{
		ID:           productID,
		Name:         productName,
		Description:  productDescription,
		CreationDate: time.Time{},
		Owner:        user.ID,
	}

	expectedError := errors.New("error creating collection indexes")

	s.mocks.accessControl.EXPECT().CheckRoleGrants(user, auth.ActCreateProduct).Return(nil)
	s.mocks.productRepo.EXPECT().GetByID(ctx, productID).Return(nil, usecase.ErrProductNotFound)
	s.mocks.productRepo.EXPECT().GetByName(ctx, productName).Return(nil, usecase.ErrProductNotFound)
	s.mocks.productRepo.EXPECT().Create(ctx, expectedProduct).Return(expectedProduct, nil)
	s.mocks.measurementRepo.EXPECT().CreateDatabase(productID).Return(nil)
	s.mocks.metricRepo.EXPECT().CreateIndexes(ctx, productID).Return(nil)
	s.mocks.processLogRepo.EXPECT().CreateIndexes(ctx, productID).Return(expectedError)

	_, err := s.productInteractor.CreateProduct(ctx, user, productID, productName, productDescription)
	require.ErrorIs(t, err, expectedError)
}

func TestCreateProduct_ErrorCreatingVersionRepoIndexes(t *testing.T) {
	s := newProductSuite(t)
	defer s.ctrl.Finish()

	ctx := context.Background()

	user := testhelpers.NewUserBuilder().Build()
	productID := "product-id"
	productName := "product-name"
	productDescription := "This is a product description"
	expectedProduct := &entity.Product{
		ID:           productID,
		Name:         productName,
		Description:  productDescription,
		CreationDate: time.Time{},
		Owner:        user.ID,
	}

	expectedError := errors.New("error creating versions collection indexes")

	s.mocks.accessControl.EXPECT().CheckRoleGrants(user, auth.ActCreateProduct).Return(nil)
	s.mocks.productRepo.EXPECT().GetByID(ctx, productID).Return(nil, usecase.ErrProductNotFound)
	s.mocks.productRepo.EXPECT().GetByName(ctx, productName).Return(nil, usecase.ErrProductNotFound)
	s.mocks.productRepo.EXPECT().Create(ctx, expectedProduct).Return(expectedProduct, nil)
	s.mocks.measurementRepo.EXPECT().CreateDatabase(productID).Return(nil)
	s.mocks.metricRepo.EXPECT().CreateIndexes(ctx, productID).Return(nil)
	s.mocks.processLogRepo.EXPECT().CreateIndexes(ctx, productID).Return(nil)
	s.mocks.versionRepo.EXPECT().CreateIndexes(ctx, productID).Return(expectedError)

	_, err := s.productInteractor.CreateProduct(ctx, user, productID, productName, productDescription)
	require.ErrorIs(t, err, expectedError)
}

func TestCreateProduct_ErrorCreatingProcessRegistryRepoIndexes(t *testing.T) {
	s := newProductSuite(t)
	defer s.ctrl.Finish()

	ctx := context.Background()

	user := testhelpers.NewUserBuilder().Build()
	productID := "product-id"
	productName := "product-name"
	productDescription := "This is a product description"
	expectedProduct := &entity.Product{
		ID:           productID,
		Name:         productName,
		Description:  productDescription,
		CreationDate: time.Time{},
		Owner:        user.ID,
	}

	expectedError := errors.New("error creating versions collection indexes")

	s.mocks.accessControl.EXPECT().CheckRoleGrants(user, auth.ActCreateProduct).Return(nil)
	s.mocks.productRepo.EXPECT().GetByID(ctx, productID).Return(nil, usecase.ErrProductNotFound)
	s.mocks.productRepo.EXPECT().GetByName(ctx, productName).Return(nil, usecase.ErrProductNotFound)
	s.mocks.productRepo.EXPECT().Create(ctx, expectedProduct).Return(expectedProduct, nil)
	s.mocks.measurementRepo.EXPECT().CreateDatabase(productID).Return(nil)
	s.mocks.metricRepo.EXPECT().CreateIndexes(ctx, productID).Return(nil)
	s.mocks.processLogRepo.EXPECT().CreateIndexes(ctx, productID).Return(nil)
	s.mocks.versionRepo.EXPECT().CreateIndexes(ctx, productID).Return(nil)
	s.mocks.processRegistryRepo.EXPECT().CreateIndexes(ctx, productID).Return(expectedError)

	_, err := s.productInteractor.CreateProduct(ctx, user, productID, productName, productDescription)
	require.ErrorIs(t, err, expectedError)
}

func TestGetByID(t *testing.T) {
	s := newProductSuite(t)
	defer s.ctrl.Finish()

	ctx := context.Background()

	user := testhelpers.NewUserBuilder().Build()
	productID := "product-id"
	productName := "product-name"

	expected := &entity.Product{
		ID:           productID,
		Name:         productName,
		Description:  "Product description...",
		CreationDate: time.Time{},
	}

	s.mocks.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActViewProduct).Return(nil)
	s.mocks.productRepo.EXPECT().GetByID(ctx, productID).Return(expected, nil)

	actual, err := s.productInteractor.GetByID(ctx, user, productID)

	require.Nil(t, err)
	require.Equal(t, expected, actual)
}

func TestFindAll(t *testing.T) {
	s := newProductSuite(t)
	defer s.ctrl.Finish()

	ctx := context.Background()

	productID := "product-id"
	productName := "product-name"

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

	s.mocks.accessControl.EXPECT().IsAdmin(user).Return(false)
	s.mocks.accessControl.EXPECT().GetUserProducts(user).Return(userProducts)
	s.mocks.productRepo.EXPECT().FindByIDs(ctx, userProducts).Return(expected, nil)

	actual, err := s.productInteractor.FindAll(ctx, user)

	require.Nil(t, err)
	require.Equal(t, expected, actual)
}

func TestFindAll_AdminUser(t *testing.T) {
	s := newProductSuite(t)
	defer s.ctrl.Finish()

	ctx := context.Background()

	productID := "product-id"
	productName := "product-name"

	user := testhelpers.NewUserBuilder().WithRoles([]string{auth.DefaultAdminRole}).Build()

	expected := []*entity.Product{
		{
			ID:           productID,
			Name:         productName,
			Description:  "Product description...",
			CreationDate: time.Time{},
		},
	}

	s.mocks.accessControl.EXPECT().IsAdmin(user).Return(true)
	s.mocks.productRepo.EXPECT().FindAll(ctx).Return(expected, nil)

	actual, err := s.productInteractor.FindAll(ctx, user)

	require.Nil(t, err)
	require.Equal(t, expected, actual)
}
