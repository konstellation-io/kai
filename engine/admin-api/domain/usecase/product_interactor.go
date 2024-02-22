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
	"github.com/konstellation-io/kai/engine/admin-api/pkg/compensator"
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
	versionRepo       repository.VersionRepo
	processRepo       repository.ProcessRepository
	userActivity      UserActivityInteracter
	accessControl     auth.AccessControl
	objectStorage     repository.ObjectStorage
	natsService       service.NatsManagerService
	userRegistry      service.UserRegistry
	passwordGenerator password.PasswordGenerator
	predictionRepo    repository.PredictionRepository
}

type ProductInteractorOpts struct {
	Logger               logr.Logger
	ProductRepo          repository.ProductRepo
	VersionRepo          repository.VersionRepo
	ProcessRepo          repository.ProcessRepository
	UserActivity         UserActivityInteracter
	AccessControl        auth.AccessControl
	ObjectStorage        repository.ObjectStorage
	NatsService          service.NatsManagerService
	UserRegistry         service.UserRegistry
	PasswordGenerator    password.PasswordGenerator
	PredictionRepository repository.PredictionRepository
}

// NewProductInteractor creates a new ProductInteractor.
func NewProductInteractor(ps *ProductInteractorOpts) *ProductInteractor {
	return &ProductInteractor{
		ps.Logger,
		ps.ProductRepo,
		ps.VersionRepo,
		ps.ProcessRepo,
		ps.UserActivity,
		ps.AccessControl,
		ps.ObjectStorage,
		ps.NatsService,
		ps.UserRegistry,
		ps.PasswordGenerator,
		ps.PredictionRepository,
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

	newProduct := i.buildProductFromParams(user, name, description)

	i.logger.Info("Creating product", "name", newProduct.Name, "ID", newProduct.ID)

	// Validation
	err := newProduct.Validate()
	if err != nil {
		return nil, err
	}

	err = i.checkIfProductExists(ctx, newProduct)
	if err != nil {
		return nil, err
	}

	// Create resources
	compensations := compensator.New()

	createdProduct, err := i.createProductResources(ctx, compensations, newProduct)
	if err != nil {
		go i.executeCompensations(compensations)
		return nil, err
	}

	i.logger.Info("Product stored in the database", "name", createdProduct.Name, "ID", createdProduct.ID)

	return createdProduct, nil
}

func (i *ProductInteractor) createProductResources(
	ctx context.Context,
	compensations *compensator.Compensator,
	newProduct *entity.Product,
) (*entity.Product, error) {
	globalKeyValueStore, err := i.natsService.CreateGlobalKeyValueStore(ctx, newProduct.ID)
	if err != nil {
		return nil, fmt.Errorf("creating global key-value store: %w", err)
	}

	compensations.AddCompensation(func() error {
		return i.natsService.DeleteGlobalKeyValueStore(context.Background(), newProduct.ID)
	})

	newProduct.KeyValueStore = globalKeyValueStore

	minioConfiguration := entity.MinioConfiguration{
		Bucket: newProduct.ID,
	}

	serviceAccount := entity.ServiceAccount{
		Username: newProduct.ID,
		Group:    newProduct.ID,
		Password: i.passwordGenerator.MustGenerate(32, 8, 8, true, true),
	}

	err = i.objectStorage.CreateBucket(ctx, minioConfiguration.Bucket)
	if err != nil {
		return nil, fmt.Errorf("creating object storage bucket: %w", err)
	}

	compensations.AddCompensation(func() error {
		return i.objectStorage.DeleteBucket(context.Background(), minioConfiguration.Bucket)
	})

	policyName, err := i.objectStorage.CreateBucketPolicy(ctx, newProduct.ID)
	if err != nil {
		return nil, fmt.Errorf("creating object storage policy: %w", err)
	}

	compensations.AddCompensation(func() error {
		return i.objectStorage.DeleteBucketPolicy(context.Background(), policyName)
	})

	err = i.userRegistry.CreateGroupWithPolicy(ctx, serviceAccount.Group, policyName)
	if err != nil {
		return nil, err
	}

	compensations.AddCompensation(func() error {
		return i.userRegistry.DeleteGroup(context.Background(), newProduct.ID)
	})

	err = i.userRegistry.CreateUserWithinGroup(
		ctx,
		serviceAccount.Username,
		serviceAccount.Password,
		serviceAccount.Group,
	)
	if err != nil {
		return nil, err
	}

	compensations.AddCompensation(func() error {
		return i.userRegistry.DeleteUser(context.Background(), newProduct.ID)
	})

	err = i.predictionRepo.CreateUser(ctx, newProduct.ID, serviceAccount.Username, serviceAccount.Password)
	if err != nil {
		return nil, fmt.Errorf("creating user in prediction's repository: %w", err)
	}

	compensations.AddCompensation(func() error {
		return i.predictionRepo.DeleteUser(context.Background(), newProduct.ID)
	})

	newProduct.MinioConfiguration = minioConfiguration
	newProduct.ServiceAccount = serviceAccount

	err = i.createDatabaseIndexes(ctx, newProduct.Name)
	if err != nil {
		return nil, err
	}

	compensations.AddCompensation(func() error {
		return i.productRepo.DeleteDatabase(context.Background(), newProduct.ID)
	})

	createdProduct, err := i.productRepo.Create(ctx, newProduct)
	if err != nil {
		return nil, err
	}

	return createdProduct, nil
}

func (i *ProductInteractor) checkIfProductExists(ctx context.Context, newProduct *entity.Product) error {
	// Check if the Product already exists
	productFromDB, err := i.productRepo.GetByID(ctx, newProduct.ID)
	if productFromDB != nil {
		return ErrProductDuplicated
	} else if !errors.Is(err, ErrProductNotFound) {
		return err
	}

	// Check if there is another Product with the same name
	productFromDB, err = i.productRepo.GetByName(ctx, newProduct.Name)
	if productFromDB != nil {
		return ErrProductDuplicatedName
	} else if !errors.Is(err, ErrProductNotFound) {
		return err
	}

	return nil
}

func (i *ProductInteractor) createDatabaseIndexes(ctx context.Context, productID string) error {
	err := i.versionRepo.CreateIndexes(ctx, productID)
	if err != nil {
		return err
	}

	err = i.processRepo.CreateIndexes(ctx, productID)
	if err != nil {
		return err
	}

	return nil
}

// GetByID return a Product by its ID.
func (i *ProductInteractor) GetByID(ctx context.Context, user *entity.User, productID string) (*entity.Product, error) {
	if err := i.accessControl.CheckProductGrants(user, productID, auth.ActViewProduct); err != nil {
		return nil, err
	}

	return i.productRepo.GetByID(ctx, productID)
}

// FindAll returns a list of all Products.
func (i *ProductInteractor) FindAll(ctx context.Context, user *entity.User, filter *repository.FindAllFilter) ([]*entity.Product, error) {
	if i.accessControl.IsAdmin(user) {
		return i.productRepo.FindAll(ctx, filter)
	}

	visibleProducts := i.accessControl.GetUserProductsWithViewAccess(user)

	return i.productRepo.FindByIDs(ctx, visibleProducts, filter)
}

func (i *ProductInteractor) generateProductID(name string) string {
	id := strings.TrimSpace(name)
	id = strings.ToLower(id)
	id = _whiteSpacesRE.ReplaceAllString(id, "-")
	id = _validCharactersRE.ReplaceAllString(id, "")

	return id
}

func (i *ProductInteractor) buildProductFromParams(user *entity.User, name, description string) *entity.Product {
	return &entity.Product{
		ID:          i.generateProductID(name),
		Name:        strings.TrimSpace(name),
		Description: strings.TrimSpace(description),
		Owner:       user.ID,
	}
}

func (i *ProductInteractor) executeCompensations(compensations *compensator.Compensator) {
	err := compensations.Execute()
	if err != nil {
		i.logger.Error(err, "Executing compensations on create product request")
	}
}
