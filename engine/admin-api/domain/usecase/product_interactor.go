package usecase

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/repository"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/auth"
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
	logger              logr.Logger
	productRepo         repository.ProductRepo
	measurementRepo     repository.MeasurementRepo
	versionRepo         repository.VersionRepo
	metricRepo          repository.MetricRepo
	processLogRepo      repository.ProcessLogRepository
	processRegistryRepo repository.ProcessRegistryRepo
	processRepo     repository.ProcessRepository
	userActivity        UserActivityInteracter
	accessControl       auth.AccessControl
}

type ProductInteractorOpts struct {
	Logger              logr.Logger
	ProductRepo         repository.ProductRepo
	MeasurementRepo     repository.MeasurementRepo
	VersionRepo         repository.VersionRepo
	MetricRepo          repository.MetricRepo
	ProcessLogRepo      repository.ProcessLogRepository
	ProcessRegistryRepo repository.ProcessRegistryRepo
	ProcessRepo     repository.ProcessRepository
	UserActivity        UserActivityInteracter
	AccessControl       auth.AccessControl
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
	}
}

// CreateProduct adds a new Product.
func (i *ProductInteractor) CreateProduct(
	ctx context.Context,
	user *entity.User,
	productID,
	name,
	description string,
) (*entity.Product, error) {
	if err := i.accessControl.CheckRoleGrants(user, auth.ActCreateProduct); err != nil {
		return nil, err
	}

	// Sanitize input params
	productID = i.generateProductID(name)
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)

	i.logger.Info("Creating product", "name", name, "id", productID)

	r := &entity.Product{
		ID:          productID,
		Name:        name,
		Description: description,
		Owner:       user.ID,
	}

	// Validation
	err := r.Validate()
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

	createdProduct, err := i.productRepo.Create(ctx, r)
	if err != nil {
		return nil, err
	}

	i.logger.Info("Product stored in the database with ID:" + createdProduct.Name)

	err = i.measurementRepo.CreateDatabase(createdProduct.Name)
	if err != nil {
		return nil, err
	}

	i.logger.Info("Measurement database created for product with ID:" + createdProduct.Name)

	err = i.createDatabaseIndexes(ctx, name)
	if err != nil {
		return nil, err
	}

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
	id = _whiteSpacesRE.ReplaceAllString(id, "_")
	id = _validCharactersRE.ReplaceAllString(id, "")

	return id
}
