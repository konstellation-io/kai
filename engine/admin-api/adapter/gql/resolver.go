package gql

//go:generate go run github.com/99designs/gqlgen --verbose

import (
	"context"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/logs"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/process"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/version"
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
	logger                 logr.Logger
	productInteractor      *usecase.ProductInteractor
	userInteractor         *usecase.UserInteractor
	userActivityInteractor usecase.UserActivityInteracter
	versionInteractor      *version.Handler
	serverInfoGetter       *usecase.ServerInfoGetter
	processHandler         *process.Handler
	logsService            logs.LogsUsecase
}

func NewGraphQLResolver(params Params) *Resolver {
	initialize()

	return &Resolver{
		params.Logger,
		params.ProductInteractor,
		params.UserInteractor,
		params.UserActivityInteractor,
		params.VersionInteractor,
		params.ServerInfoGetter,
		params.ProcessHandler,
		params.LogsUsecase,
	}
}

func (r *mutationResolver) CreateProduct(ctx context.Context, input CreateProductInput) (*entity.Product, error) {
	loggedUser := ctx.Value("user").(*entity.User)

	product, err := r.productInteractor.CreateProduct(ctx, loggedUser, input.Name, input.Description)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *mutationResolver) CreateVersion(ctx context.Context, input CreateVersionInput) (*entity.Version, error) {
	loggedUser := ctx.Value("user").(*entity.User)

	return r.versionInteractor.Create(ctx, loggedUser, input.ProductID, input.File.File)
}

func (r *mutationResolver) RegisterProcess(ctx context.Context, input RegisterProcessInput) (*entity.RegisteredProcess, error) {
	loggedUser := ctx.Value("user").(*entity.User)

	return r.processHandler.RegisterProcess(
		ctx, loggedUser, process.RegisterProcessOpts{
			Product:     input.ProductID,
			Version:     input.Version,
			Process:     input.ProcessID,
			ProcessType: entity.ProcessType(input.ProcessType),
			Sources:     input.File.File,
		},
	)
}

func (r *mutationResolver) RegisterPublicProcess(ctx context.Context, input RegisterPublicProcessInput) (*entity.RegisteredProcess, error) {
	loggedUser := ctx.Value("user").(*entity.User)

	return r.processHandler.RegisterProcess(
		ctx, loggedUser,
		process.RegisterProcessOpts{
			Version:     input.Version,
			Process:     input.ProcessID,
			ProcessType: entity.ProcessType(input.ProcessType),
			IsPublic:    true,
			Sources:     input.File.File,
		},
	)
}

func (r *mutationResolver) DeleteProcess(ctx context.Context, input DeleteProcessInput) (string, error) {
	loggedUser := ctx.Value("user").(*entity.User)

	return r.processHandler.DeleteProcess(
		ctx, loggedUser,
		process.DeleteProcessOpts{
			Product:  input.ProductID,
			Version:  input.Version,
			Process:  input.ProcessID,
			IsPublic: false,
		},
	)
}

func (r *mutationResolver) DeletePublicProcess(ctx context.Context, input DeletePublicProcessInput) (string, error) {
	loggedUser := ctx.Value("user").(*entity.User)

	return r.processHandler.DeleteProcess(
		ctx, loggedUser,
		process.DeleteProcessOpts{
			Version:  input.Version,
			Process:  input.ProcessID,
			IsPublic: true,
		},
	)
}

func (r *mutationResolver) StartVersion(ctx context.Context, input StartVersionInput) (*entity.Version, error) {
	loggedUser := ctx.Value("user").(*entity.User)

	v, _, err := r.versionInteractor.Start(ctx, loggedUser, input.ProductID, input.VersionTag, input.Comment)
	if err != nil {
		r.logger.Error(err, "Unable to start version",
			"productID", input.ProductID,
			"versionTag", input.VersionTag,
		)

		return nil, err
	}

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
				r.logger.V(2).Info("[notifyVersionStatus] received nil on notifyCh. closing notifier")
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

func (r *mutationResolver) PublishVersion(ctx context.Context, input PublishVersionInput) ([]*entity.PublishedTrigger, error) {
	loggedUser := ctx.Value("user").(*entity.User)

	urls, err := r.versionInteractor.Publish(ctx, loggedUser, version.PublishOpts{
		ProductID:  input.ProductID,
		VersionTag: input.VersionTag,
		Comment:    input.Comment,
		Force:      input.Force,
	})
	if err != nil {
		return nil, err
	}

	publishedTriggers := make([]*entity.PublishedTrigger, 0, len(urls))
	for trigger, url := range urls {
		publishedTriggers = append(publishedTriggers, &entity.PublishedTrigger{
			Trigger: trigger,
			URL:     url,
		})
	}

	return publishedTriggers, nil
}

func (r *mutationResolver) UpdateUserProductGrants(
	ctx context.Context,
	input UpdateUserProductGrantsInput,
) (*entity.User, error) {
	user := ctx.Value("user").(*entity.User)

	err := r.userInteractor.UpdateUserProductGrants(
		ctx,
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

	err := r.userInteractor.RevokeUserProductGrants(ctx, loggedUser, input.TargetID, input.Product, *input.Comment)
	if err != nil {
		return nil, err
	}

	return &entity.User{ID: input.TargetID}, nil
}

func (r *queryResolver) Product(ctx context.Context, id string) (*entity.Product, error) {
	loggedUser := ctx.Value("user").(*entity.User)
	return r.productInteractor.GetByID(ctx, loggedUser, id)
}

func (r *queryResolver) Products(ctx context.Context) ([]*entity.Product, error) {
	loggedUser := ctx.Value("user").(*entity.User)
	return r.productInteractor.FindAll(ctx, loggedUser)
}

func (r *queryResolver) Version(ctx context.Context, productID string, tag *string) (*entity.Version, error) {
	loggedUser := ctx.Value("user").(*entity.User)

	if tag == nil {
		return r.versionInteractor.GetLatest(ctx, loggedUser, productID)
	} else {
		return r.versionInteractor.GetByTag(ctx, loggedUser, productID, *tag)
	}
}

func (r *queryResolver) Versions(ctx context.Context, productID string) ([]*entity.Version, error) {
	loggedUser := ctx.Value("user").(*entity.User)
	return r.versionInteractor.ListVersionsByProduct(ctx, loggedUser, productID)
}

func (r *queryResolver) RegisteredProcesses(
	ctx context.Context, productID string,
	processName, version, processType *string,
) ([]*entity.RegisteredProcess, error) {
	loggedUser := ctx.Value("user").(*entity.User)

	var filter entity.SearchFilter

	if processName != nil {
		filter.ProcessName = *processName
	}

	if version != nil {
		filter.Version = *version
	}

	if processType != nil {
		filter.ProcessType = entity.ProcessType(*processType)
	}

	return r.processHandler.Search(ctx, loggedUser, productID, &filter)
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
	filters entity.LogFilters,
) ([]*entity.Log, error) {
	return r.logsService.GetLogs(filters)
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

func (r *registeredProcessResolver) UploadDate(_ context.Context, obj *entity.RegisteredProcess) (string, error) {
	return obj.UploadDate.Format(time.RFC3339), nil
}

func (r *registeredProcessResolver) Type(_ context.Context, obj *entity.RegisteredProcess) (string, error) {
	return obj.Type.String(), nil
}

func (r *logFiltersResolver) From(_ context.Context, obj *entity.LogFilters, from string) error {
	var err error
	obj.From, err = time.Parse(time.RFC3339, from)

	return err
}

func (r *logFiltersResolver) To(_ context.Context, obj *entity.LogFilters, to string) error {
	var err error
	obj.To, err = time.Parse(time.RFC3339, to)

	return err
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Product returns ProductResolver implementation.
func (r *Resolver) Product() ProductResolver { return &productResolver{r} }

// UserActivity returns UserActivityResolver implementation.
func (r *Resolver) UserActivity() UserActivityResolver { return &userActivityResolver{r} }

// Version returns VersionResolver implementation.
func (r *Resolver) Version() VersionResolver { return &versionResolver{r} }

// RegisteredProcess returns RegisteredProcessResolver implementation.
func (r *Resolver) RegisteredProcess() RegisteredProcessResolver {
	return &registeredProcessResolver{r}
}

// LogFilters returns LogFiltersResolver implementation.
func (r *Resolver) LogFilters() LogFiltersResolver { return &logFiltersResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type productResolver struct{ *Resolver }
type userActivityResolver struct{ *Resolver }
type versionResolver struct{ *Resolver }
type registeredProcessResolver struct{ *Resolver }

type logFiltersResolver struct{ *Resolver }
