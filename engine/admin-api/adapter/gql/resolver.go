package gql

//go:generate go run github.com/99designs/gqlgen --verbose

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/logging"
)

//nolint:gochecknoglobals // needs to be global to be used in the resolver
var versionStatusChannels map[string]chan *entity.Version

//nolint:gochecknoglobals // needs to be global to be used in the resolver
var mux *sync.RWMutex

func initialize() {
	versionStatusChannels = map[string]chan *entity.Version{}
	mux = &sync.RWMutex{}
}

type Resolver struct {
	logger                 logging.Logger
	productInteractor      *usecase.ProductInteractor
	userInteractor         *usecase.UserInteractor
	userActivityInteractor usecase.UserActivityInteracter
	versionInteractor      *usecase.VersionInteractor
	metricsInteractor      *usecase.MetricsInteractor
	serverInfoGetter       *usecase.ServerInfoGetter
	processService         *usecase.ProcessService
	cfg                    *config.Config
}

func NewGraphQLResolver(params Params) *Resolver {
	initialize()

	return &Resolver{
		params.Logger,
		params.ProductInteractor,
		params.UserInteractor,
		params.UserActivityInteractor,
		params.VersionInteractor,
		params.MetricsInteractor,
		params.ServerInfoGetter,
		params.ProcessService,
		params.Cfg,
	}
}

func (r *mutationResolver) CreateProduct(ctx context.Context, input CreateProductInput) (*entity.Product, error) {
	loggedUser := ctx.Value("user").(*entity.User)

	r.logger.Debug("Creating product with id " + input.ID)

	product, err := r.productInteractor.CreateProduct(ctx, loggedUser, input.ID, input.Name, input.Description)
	if err != nil {
		r.logger.Error("Error creating product: " + err.Error())
		return nil, err
	}

	return product, nil
}

func (r *mutationResolver) CreateVersion(ctx context.Context, input CreateVersionInput) (*entity.Version, error) {
	loggedUser := ctx.Value("user").(*entity.User)

	version, notifyCh, err := r.versionInteractor.Create(ctx, loggedUser, input.ProductID, input.File.File)
	if err != nil {
		return nil, err
	}

	go r.notifyVersionStatus(notifyCh)

	return version, nil
}

func (r *mutationResolver) RegisterProcess(ctx context.Context, input RegisterProcessInput) (*RegisteredImage, error) {
	_ = ctx.Value("user").(*entity.User)

	registeredImage, err := r.processService.RegisterProcess(ctx, input.ProductID, input.Version, input.ProcessID, input.File.File)
	if err != nil {
		return nil, err
	}

	return &RegisteredImage{
		ProcessedImageID: registeredImage,
	}, nil
}

func (r *mutationResolver) StartVersion(ctx context.Context, input StartVersionInput) (*entity.Version, error) {
	loggedUser := ctx.Value("user").(*entity.User)

	v, notifyCh, err := r.versionInteractor.Start(ctx, loggedUser, input.ProductID, input.VersionTag, input.Comment)
	if err != nil {
		r.logger.Errorf("[mutationResolver.StartVersion] errors starting version: %s", err)
		return nil, err
	}

	go r.notifyVersionStatus(notifyCh)

	return v, err
}

func (r *mutationResolver) StopVersion(ctx context.Context, input StopVersionInput) (*entity.Version, error) {
	loggedUser := ctx.Value("user").(*entity.User)

	v, notifyCh, err := r.versionInteractor.Stop(ctx, loggedUser, input.ProductID, input.VersionTag, input.Comment)
	if err != nil {
		return nil, err
	}

	go r.notifyVersionStatus(notifyCh)

	return v, err
}

func (r *mutationResolver) notifyVersionStatus(notifyCh chan *entity.Version) {
	//nolint:gosimple // legacy code
	for {
		select {
		case v, ok := <-notifyCh:
			if !ok {
				r.logger.Debugf("[notifyVersionStatus] received nil on notifyCh. closing notifier")
				return
			}

			mux.RLock()
			for _, vs := range versionStatusChannels {
				vs <- v
			}
			mux.RUnlock()
		}
	}
}

func (r *mutationResolver) UnpublishVersion(ctx context.Context, input UnpublishVersionInput) (*entity.Version, error) {
	loggedUser := ctx.Value("user").(*entity.User)
	return r.versionInteractor.Unpublish(ctx, loggedUser, input.ProductID, input.VersionTag, input.Comment)
}

func (r *mutationResolver) PublishVersion(ctx context.Context, input PublishVersionInput) (*entity.Version, error) {
	loggedUser := ctx.Value("user").(*entity.User)
	return r.versionInteractor.Publish(ctx, loggedUser, input.ProductID, input.VersionTag, input.Comment)
}

func (r *mutationResolver) UpdateUserProductGrants(
	ctx context.Context,
	input UpdateUserProductGrantsInput,
) (*entity.User, error) {
	user := ctx.Value("user").(*entity.User)

	err := r.userInteractor.UpdateUserProductGrants(
		user,
		input.TargetID,
		input.Product,
		input.Grants,
		*input.Comment,
	)
	if err != nil {
		return nil, err
	}

	return &entity.User{ID: input.TargetID}, nil
}

func (r *mutationResolver) RevokeUserProductGrants(
	ctx context.Context,
	input RevokeUserProductGrantsInput,
) (*entity.User, error) {
	loggedUser := ctx.Value("user").(*entity.User)

	err := r.userInteractor.RevokeUserProductGrants(loggedUser, input.TargetID, input.Product, *input.Comment)
	if err != nil {
		return nil, err
	}

	return &entity.User{ID: input.TargetID}, nil
}

func (r *queryResolver) Metrics(ctx context.Context, productID, versionTag, startDate, endDate string) (*entity.Metrics, error) {
	loggedUser := ctx.Value("user").(*entity.User)
	return r.metricsInteractor.GetMetrics(ctx, loggedUser, productID, versionTag, startDate, endDate)
}

func (r *queryResolver) Product(ctx context.Context, id string) (*entity.Product, error) {
	loggedUser := ctx.Value("user").(*entity.User)
	return r.productInteractor.GetByID(ctx, loggedUser, id)
}

func (r *queryResolver) Products(ctx context.Context) ([]*entity.Product, error) {
	loggedUser := ctx.Value("user").(*entity.User)
	return r.productInteractor.FindAll(ctx, loggedUser)
}

func (r *queryResolver) Version(ctx context.Context, tag, productID string) (*entity.Version, error) {
	loggedUser := ctx.Value("user").(*entity.User)
	return r.versionInteractor.GetByTag(ctx, loggedUser, productID, tag)
}

func (r *queryResolver) Versions(ctx context.Context, productID string) ([]*entity.Version, error) {
	loggedUser := ctx.Value("user").(*entity.User)
	return r.versionInteractor.ListVersionsByProduct(ctx, loggedUser, productID)
}

func (r *queryResolver) UserActivityList(
	ctx context.Context,
	userEmail *string,
	types []entity.UserActivityType,
	versionIds []string,
	fromDate *string,
	toDate *string,
	lastID *string,
) ([]*entity.UserActivity, error) {
	loggedUser := ctx.Value("user").(*entity.User)
	return r.userActivityInteractor.Get(ctx, loggedUser, userEmail, types, versionIds, fromDate, toDate, lastID)
}

func (r *queryResolver) Logs(
	ctx context.Context,
	productID string,
	filters entity.LogFilters,
	cursor *string,
) (*LogPage, error) {
	loggedUser := ctx.Value("user").(*entity.User)

	searchResult, err := r.versionInteractor.SearchLogs(ctx, loggedUser, productID, filters, cursor)
	if err != nil {
		return nil, err
	}

	nextCursor := new(string)
	if searchResult.Cursor != "" {
		*nextCursor = searchResult.Cursor
	}

	return &LogPage{
		Cursor: nextCursor,
		Items:  searchResult.Logs,
	}, nil
}

func (r *queryResolver) ServerInfo(ctx context.Context) (*entity.ServerInfo, error) {
	user, _ := ctx.Value("user").(*entity.User)
	return r.serverInfoGetter.GetKAIServerInfo(ctx, user)
}

func (r *productResolver) CreationAuthor(_ context.Context, product *entity.Product) (string, error) {
	return product.Owner, nil
}

func (r *productResolver) CreationDate(_ context.Context, obj *entity.Product) (string, error) {
	return obj.CreationDate.Format(time.RFC3339), nil
}

func (r *productResolver) PublishedVersion(_ context.Context, obj *entity.Product) (*entity.Version, error) {
	if obj.PublishedVersion != "" {
		return r.versionInteractor.GetByID(obj.ID, obj.PublishedVersion)
	}

	return nil, nil
}

func (r *productResolver) MeasurementsURL(_ context.Context, _ *entity.Product) (string, error) {
	return fmt.Sprintf("%s/measurements/%s", r.cfg.Admin.BaseURL, r.cfg.K8s.Namespace), nil
}

func (r *productResolver) DatabaseURL(_ context.Context, _ *entity.Product) (string, error) {
	return fmt.Sprintf("%s/database/%s", r.cfg.Admin.BaseURL, r.cfg.K8s.Namespace), nil
}

func (r *productResolver) EntrypointAddress(_ context.Context, _ *entity.Product) (string, error) {
	return fmt.Sprintf("entrypoint.%s", r.cfg.BaseDomainName), nil
}

func (r *subscriptionResolver) WatchProcessLogs(ctx context.Context, productID, versionTag string,
	filters entity.LogFilters) (<-chan *entity.ProcessLog, error) {
	loggedUser := ctx.Value("user").(*entity.User)
	return r.versionInteractor.WatchProcessLogs(ctx, loggedUser, productID, versionTag, filters)
}

func (r *userActivityResolver) Date(_ context.Context, obj *entity.UserActivity) (string, error) {
	return obj.Date.Format(time.RFC3339), nil
}

func (r *userActivityResolver) User(_ context.Context, obj *entity.UserActivity) (string, error) {
	return obj.UserID, nil
}

func (r *versionResolver) CreationDate(_ context.Context, obj *entity.Version) (string, error) {
	return obj.CreationDate.Format(time.RFC3339), nil
}

func (r *versionResolver) CreationAuthor(_ context.Context, obj *entity.Version) (string, error) {
	return obj.CreationAuthor, nil
}

func (r *versionResolver) PublicationDate(_ context.Context, obj *entity.Version) (*string, error) {
	if obj.PublicationDate == nil {
		return nil, nil
	}

	result := obj.PublicationDate.Format(time.RFC3339)

	return &result, nil
}

func (r *versionResolver) PublicationAuthor(_ context.Context, obj *entity.Version) (*string, error) {
	if obj.PublicationAuthor == nil {
		return nil, nil
	}

	return obj.PublicationAuthor, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Product returns ProductResolver implementation.
func (r *Resolver) Product() ProductResolver { return &productResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

// UserActivity returns UserActivityResolver implementation.
func (r *Resolver) UserActivity() UserActivityResolver { return &userActivityResolver{r} }

// Version returns VersionResolver implementation.
func (r *Resolver) Version() VersionResolver { return &versionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type productResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
type userActivityResolver struct{ *Resolver }
type versionResolver struct{ *Resolver }
