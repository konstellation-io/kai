package usecase

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/repository"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/sethvargo/go-password/password"
)

var (
	ErrProductNotFound       = errors.New("error product not found")
	ErrProductDuplicated     = errors.New("there is already a product with the same id")
	ErrProductDuplicatedName = errors.New("there is already a product with the same name")

	_whiteSpacesRE     = regexp.MustCompile(" +")
	_validCharactersRE = regexp.MustCompile("[^a-z0-9-]")
)

// ProductInteractor contains app logic to handle Product entities.
type ProductInteractor struct {
	logger            logr.Logger
	productRepo       repository.ProductRepo
	measurementRepo   repository.MeasurementRepo
	versionRepo       repository.VersionRepo
	metricRepo        repository.MetricRepo
	processLogRepo    repository.ProcessLogRepository
	processRepo       repository.ProcessRepository
	userActivity      UserActivityInteracter
	accessControl     auth.AccessControl
	objectStorage     repository.ObjectStorage
	natsService       service.NatsManagerService
	userRegistry      service.UserRegistry
	passwordGenerator password.PasswordGenerator
}

type ProductInteractorOpts struct {
	Logger            logr.Logger
	ProductRepo       repository.ProductRepo
	MeasurementRepo   repository.MeasurementRepo
	VersionRepo       repository.VersionRepo
	MetricRepo        repository.MetricRepo
	ProcessLogRepo    repository.ProcessLogRepository
	ProcessRepo       repository.ProcessRepository
	UserActivity      UserActivityInteracter
	AccessControl     auth.AccessControl
	ObjectStorage     repository.ObjectStorage
	NatsService       service.NatsManagerService
	UserRegistry      service.UserRegistry
	PasswordGenerator password.PasswordGenerator
}

// NewProductInteractor creates a new ProductInteractor.
func NewProductInteractor(ps *ProductInteractorOpts) *ProductInteractor {
	return &ProductInteractor{
		ps.Logger,
		ps.ProductRepo,
		ps.MeasurementRepo,
		ps.VersionRepo,
		ps.MetricRepo,
		ps.ProcessLogRepo,
		ps.ProcessRepo,
		ps.UserActivity,
		ps.AccessControl,
		ps.ObjectStorage,
		ps.NatsService,
		ps.UserRegistry,
		ps.PasswordGenerator,
	}
}

// CreateProduct adds a new Product.
func (i *ProductInteractor) CreateProduct(
	ctx context.Context,
	user *entity.User,
	name,
	description string,
) (*entity.Product, error) {
	if err := i.accessControl.CheckRoleGrants(user, auth.ActCreateProduct); err != nil {
		return nil, err
	}

	// Sanitize input params
	productID := i.generateProductID(name)
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)

	i.logger.Info("Creating product", "name", name, "ID", productID)

	newProduct := &entity.Product{
		ID:          productID,
		Name:        name,
		Description: description,
		Owner:       user.ID,
	}

	// Validation
	err := newProduct.Validate()
	if err != nil {
		return nil, err
	}

	// Check if the Product already exists
	productFromDB, err := i.productRepo.GetByID(ctx, productID)
	if productFromDB != nil {
		return nil, ErrProductDuplicated
	} else if !errors.Is(err, ErrProductNotFound) {
		return nil, err
	}

	// Check if there is another Product with the same name
	productFromDB, err = i.productRepo.GetByName(ctx, name)
	if productFromDB != nil {
		return nil, ErrProductDuplicatedName
	} else if !errors.Is(err, ErrProductNotFound) {
		return nil, err
	}

	globalKeyValueStore, err := i.natsService.CreateGlobalKeyValueStore(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("creating global key-value store: %w", err)
	}

	newProduct.KeyValueStore = globalKeyValueStore

	minioConfiguration := entity.MinioConfiguration{
		User:     productID,
		Group:    productID,
		Password: i.passwordGenerator.MustGenerate(32, 8, 8, true, true),
		Bucket:   productID,
	}

	err = i.objectStorage.CreateBucket(ctx, minioConfiguration.Bucket)
	if err != nil {
		return nil, fmt.Errorf("creating object storage bucket: %w", err)
	}

	policyName, err := i.objectStorage.CreateBucketPolicy(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("creating object storage policy: %w", err)
	}

	err = i.userRegistry.CreateGroupWithPolicy(ctx, minioConfiguration.Group, policyName)
	if err != nil {
		return nil, err
	}

	err = i.userRegistry.CreateUserWithinGroup(
		ctx,
		minioConfiguration.User,
		minioConfiguration.Password,
		minioConfiguration.Group,
	)
	if err != nil {
		return nil, err
	}

	newProduct.MinioConfiguration = minioConfiguration

	err = i.measurementRepo.CreateDatabase(newProduct.Name)
	if err != nil {
		return nil, err
	}

	i.logger.V(0).Info("Measurement database created for product")

	err = i.createDatabaseIndexes(ctx, name)
	if err != nil {
		return nil, err
	}

	createdProduct, err := i.productRepo.Create(ctx, newProduct)
	if err != nil {
		return nil, err
	}

	i.logger.Info("Product stored in the database", "name", createdProduct.Name, "ID", createdProduct.ID)

	return createdProduct, nil
}

func (i *ProductInteractor) createDatabaseIndexes(ctx context.Context, productID string) error {
	err := i.metricRepo.CreateIndexes(ctx, productID)
	if err != nil {
		return err
	}

	err = i.processLogRepo.CreateIndexes(ctx, productID)
	if err != nil {
		return err
	}

	err = i.versionRepo.CreateIndexes(ctx, productID)
	if err != nil {
		return err
	}

	err = i.processRepo.CreateIndexes(ctx, productID)
	if err != nil {
		return err
	}

	return nil
}

// Get return product by ID.
func (i *ProductInteractor) Get(ctx context.Context, user *entity.User, productID string) (*entity.Product, error) {
	if err := i.accessControl.CheckProductGrants(user, productID, auth.ActViewProduct); err != nil {
		return nil, err
	}

	return i.productRepo.Get(ctx)
}

// GetByID return a Product by its ID.
func (i *ProductInteractor) GetByID(ctx context.Context, user *entity.User, productID string) (*entity.Product, error) {
	if err := i.accessControl.CheckProductGrants(user, productID, auth.ActViewProduct); err != nil {
		return nil, err
	}

	return i.productRepo.GetByID(ctx, productID)
}

// FindAll returns a list of all Products.
func (i *ProductInteractor) FindAll(ctx context.Context, user *entity.User) ([]*entity.Product, error) {
	if i.accessControl.IsAdmin(user) {
		return i.productRepo.FindAll(ctx)
	}

	visibleProducts := i.accessControl.GetUserProducts(user)

	return i.productRepo.FindByIDs(ctx, visibleProducts)
}

func (i *ProductInteractor) generateProductID(name string) string {
	id := strings.TrimSpace(name)
	id = strings.ToLower(id)
	id = _whiteSpacesRE.ReplaceAllString(id, "-")
	id = _validCharactersRE.ReplaceAllString(id, "")

	return id
}
