package usecase

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/konstellation-io/krt/pkg/parse"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/repository"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/auth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/errors"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/logging"
)

// VersionInteractor contains app logic about Version entities.
type VersionInteractor struct {
	cfg                    *config.Config
	logger                 logging.Logger
	versionRepo            repository.VersionRepo
	productRepo            repository.ProductRepo
	k8sService             service.K8sService
	natsManagerService     service.NatsManagerService
	userActivityInteractor UserActivityInteracter
	accessControl          auth.AccessControl
	dashboardService       service.DashboardService
	processLogRepo         repository.ProcessLogRepository
}

// NewVersionInteractor creates a new interactor.
func NewVersionInteractor(
	cfg *config.Config,
	logger logging.Logger,
	versionRepo repository.VersionRepo,
	productRepo repository.ProductRepo,
	k8sService service.K8sService,
	natsManagerService service.NatsManagerService,
	userActivityInteractor UserActivityInteracter,
	accessControl auth.AccessControl,
	dashboardService service.DashboardService,
	processLogRepo repository.ProcessLogRepository,
) *VersionInteractor {
	return &VersionInteractor{
		cfg,
		logger,
		versionRepo,
		productRepo,
		k8sService,
		natsManagerService,
		userActivityInteractor,
		accessControl,
		dashboardService,
		processLogRepo,
	}
}

// GetByID returns a Version by its unique ID.
func (i *VersionInteractor) GetByID(productID, versionID string) (*entity.Version, error) {
	return i.versionRepo.GetByID(productID, versionID)
}

// GetByName returns a Version by its unique name.
func (i *VersionInteractor) GetByName(ctx context.Context, user *entity.User, productID, name string) (*entity.Version, error) {
	return i.versionRepo.GetByName(ctx, productID, name)
}

// GetByProduct returns all Versions of the given Product.
func (i *VersionInteractor) GetVersionsByProduct(ctx context.Context, user *entity.User, productID string) ([]*entity.Version, error) {
	versions, err := i.versionRepo.GetVersionsByProduct(ctx, productID)
	if err != nil {
		return nil, err
	}

	return versions, nil
}

func (i *VersionInteractor) copyStreamToTempFile(krtFile io.Reader) (*os.File, error) {
	tmpFile, err := os.CreateTemp("", "version")

	if err != nil {
		return nil, fmt.Errorf("error creating temp file for version: %w", err)
	}

	_, err = io.Copy(tmpFile, krtFile)
	if err != nil {
		return nil, fmt.Errorf("error copying temp file for version: %w", err)
	}

	i.logger.Infof("Created temp file: %s", tmpFile.Name())

	return tmpFile, nil
}

// Create creates a Version on the DB based on the content of a KRT file.
func (i *VersionInteractor) Create(
	ctx context.Context,
	user *entity.User,
	productID string,
	krtFile io.Reader,
) (*entity.Version, chan *entity.Version, error) {
	if err := i.accessControl.CheckProductGrants(user, productID, auth.ActCreateVersion); err != nil {
		return nil, nil, err
	}

	product, err := i.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, nil, fmt.Errorf("error product repo GetById: %w", err)
	}

	tmpDir, err := os.MkdirTemp("", "version")
	if err != nil {
		return nil, nil, fmt.Errorf("error creating temp dir for version: %w", err)
	}

	i.logger.Info("Created temp dir to extract the KRT files at " + tmpDir)

	tmpKrtFile, err := i.copyStreamToTempFile(krtFile)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating temp krt file for version: %w", err)
	}

	krtYml, err := parse.ParseFile(tmpKrtFile.Name())
	if err != nil {
		return nil, nil, errors.ParsingKRTFileError(err)
	}

	err = krtYml.Validate()
	if err != nil {
		return nil, nil, errors.NewErrInvalidKRT(
			"create version: invalid KRT file",
			err,
		)
	}

	_, err = i.versionRepo.GetByName(ctx, productID, krtYml.Name)
	if err != nil && !errors.Is(err, errors.ErrVersionNotFound) {
		return nil, nil, fmt.Errorf("error version repo GetByName: %w", err)
	} else if err == nil {
		return nil, nil, errors.ErrVersionDuplicated
	}

	versionCreated, err := i.versionRepo.Create(
		user.ID,
		productID,
		service.MapKrtYamlToVersion(krtYml),
	)

	if err != nil {
		return nil, nil, err
	}

	i.logger.Info("Version created")

	notifyStatusCh := make(chan *entity.Version, 1)

	go i.completeVersionCreation(
		user.ID, tmpKrtFile, tmpDir, product, versionCreated, notifyStatusCh,
	)

	return versionCreated, notifyStatusCh, nil
}

func (i *VersionInteractor) completeVersionCreation(
	loggedUserID string,
	tmpKrtFile *os.File,
	tmpDir string,
	product *entity.Product,
	versionCreated *entity.Version,
	notifyStatusCh chan *entity.Version,
) {
	ctx := context.Background()

	defer close(notifyStatusCh)

	defer func() {
		err := tmpKrtFile.Close()
		if err != nil {
			i.logger.Errorf("error closing file: %s", err)
			return
		}

		err = os.Remove(tmpKrtFile.Name())
		if err != nil {
			i.logger.Errorf("error removing file: %s", err)
		}
	}()

	var contentErrors []error

	dashboardsFolder := path.Join(tmpDir, "metrics/dashboards")
	contentErrors = i.saveKRTDashboards(ctx, dashboardsFolder, product, versionCreated, contentErrors)

	err := i.versionRepo.UploadKRTFile(product.ID, versionCreated, tmpKrtFile.Name())
	if err != nil {
		contentErrors = append(contentErrors, errors.ErrStoringKRTFile)
	}

	if len(contentErrors) > 0 {
		i.setStatusError(ctx, product.ID, versionCreated, contentErrors, notifyStatusCh)
		return
	}

	err = i.versionRepo.SetStatus(ctx, product.ID, versionCreated.ID, entity.VersionStatusCreated)
	if err != nil {
		versionCreated.Status = entity.VersionStatusError

		i.logger.Errorf("error setting version status: %s", err)
		notifyStatusCh <- versionCreated

		return
	}

	// Notify state
	versionCreated.Status = entity.VersionStatusCreated
	notifyStatusCh <- versionCreated

	err = i.userActivityInteractor.RegisterCreateAction(loggedUserID, product.ID, versionCreated)
	if err != nil {
		i.logger.Errorf("error registering activity: %s", err)
	}
}

// TODO discuss what will happen with this.
//
//nolint:godox // To be done.
func (i *VersionInteractor) saveKRTDashboards(
	ctx context.Context,
	dashboardsFolder string,
	product *entity.Product,
	versionCreated *entity.Version,
	contentErrors []error,
) []error {
	if _, err := os.Stat(path.Join(dashboardsFolder)); err == nil {
		err := i.storeDashboards(ctx, dashboardsFolder, product.ID, versionCreated.Name)
		if err != nil {
			contentErrors = append(contentErrors, errors.ErrCreatingDashboard)
		}
	}

	return contentErrors
}

// Start create the resources of the given Version.
func (i *VersionInteractor) Start(
	ctx context.Context,
	user *entity.User,
	productID string,
	versionName string,
	comment string,
) (*entity.Version, chan *entity.Version, error) {
	if err := i.accessControl.CheckProductGrants(user, productID, auth.ActStartVersion); err != nil {
		return nil, nil, err
	}

	i.logger.Infof("The user %q is starting version %q on product %q", user.ID, versionName, productID)

	v, err := i.versionRepo.GetByName(ctx, productID, versionName)
	if err != nil {
		return nil, nil, err
	}

	if !v.CanBeStarted() {
		return nil, nil, errors.ErrInvalidVersionStatusBeforeStarting
	}

	notifyStatusCh := make(chan *entity.Version, 1)

	err = i.versionRepo.SetStatus(ctx, productID, v.ID, entity.VersionStatusStarting)
	if err != nil {
		return nil, nil, err
	}

	// Notify intermediate state
	v.Status = entity.VersionStatusStarting
	notifyStatusCh <- v

	err = i.userActivityInteractor.RegisterStartAction(user.ID, productID, v, comment)
	if err != nil {
		return nil, nil, err
	}

	versionStreamCfg, err := i.natsManagerService.CreateStreams(ctx, productID, v)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating streams for version %q: %w", v.Name, err)
	}

	objectStoreCfg, err := i.natsManagerService.CreateObjectStores(ctx, productID, v)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating objects stores for version %q: %w", v.Name, err)
	}

	kvStoreCfg, err := i.natsManagerService.CreateKeyValueStores(ctx, productID, v)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating key-value stores for version %q: %w", v.Name, err)
	}

	versionCfg := entity.NewVersionConfig(versionStreamCfg, objectStoreCfg, kvStoreCfg)

	go i.startAndNotify(productID, v, versionCfg, notifyStatusCh)

	return v, notifyStatusCh, nil
}

// Stop removes the resources of the given Version.
func (i *VersionInteractor) Stop(
	ctx context.Context,
	user *entity.User,
	productID string,
	versionName string,
	comment string,
) (*entity.Version, chan *entity.Version, error) {
	if err := i.accessControl.CheckProductGrants(user, productID, auth.ActStartVersion); err != nil {
		return nil, nil, err
	}

	i.logger.Infof("The user %q is stopping version %q on product %q", user.ID, versionName, productID)

	v, err := i.versionRepo.GetByName(ctx, productID, versionName)
	if err != nil {
		return nil, nil, err
	}

	if !v.CanBeStopped() {
		return nil, nil, errors.ErrInvalidVersionStatusBeforeStopping
	}

	err = i.versionRepo.SetStatus(ctx, productID, v.ID, entity.VersionStatusStopping)
	if err != nil {
		return nil, nil, err
	}

	notifyStatusCh := make(chan *entity.Version, 1)

	// Notify intermediate state
	v.Status = entity.VersionStatusStopping
	notifyStatusCh <- v

	err = i.userActivityInteractor.RegisterStopAction(user.ID, productID, v, comment)
	if err != nil {
		return nil, nil, err
	}

	err = i.natsManagerService.DeleteStreams(ctx, productID, versionName)
	if err != nil {
		return nil, nil, fmt.Errorf("error stopping version %q: %w", versionName, err)
	}

	err = i.natsManagerService.DeleteObjectStores(ctx, productID, versionName)
	if err != nil {
		return nil, nil, fmt.Errorf("error stopping version %q: %w", versionName, err)
	}

	go i.stopAndNotify(productID, v, notifyStatusCh)

	return v, notifyStatusCh, nil
}

func (i *VersionInteractor) startAndNotify(
	productID string,
	vers *entity.Version,
	versionConfig *entity.VersionConfig,
	notifyStatusCh chan *entity.Version,
) {
	// WARNING: This function doesn't handle error because there is no  ERROR status defined for a Version
	ctx, cancel := context.WithTimeout(context.Background(), i.cfg.Application.VersionStatusTimeout)
	defer func() {
		cancel()
		close(notifyStatusCh)
		i.logger.Debug("[versionInteractor.startAndNotify] channel closed")
	}()

	err := i.k8sService.Start(ctx, productID, vers, versionConfig)
	if err != nil {
		i.logger.Errorf("[versionInteractor.startAndNotify] error starting version %q: %s", vers.Name, err)
	}

	err = i.versionRepo.SetStatus(ctx, productID, vers.ID, entity.VersionStatusStarted)
	if err != nil {
		i.logger.Errorf("[versionInteractor.startAndNotify] error starting version %q: %s", vers.Name, err)
	}

	vers.Status = entity.VersionStatusStarted
	notifyStatusCh <- vers
	i.logger.Infof("[versionInteractor.startAndNotify] version %q started", vers.Name)
}

func (i *VersionInteractor) stopAndNotify(
	productID string,
	vers *entity.Version,
	notifyStatusCh chan *entity.Version,
) {
	// WARNING: This function doesn't handle error because there is no  ERROR status defined for a Version
	ctx, cancel := context.WithTimeout(context.Background(), i.cfg.Application.VersionStatusTimeout)
	defer func() {
		cancel()
		close(notifyStatusCh)
		i.logger.Debug("[versionInteractor.stopAndNotify] channel closed")
	}()

	err := i.k8sService.Stop(ctx, productID, vers)
	if err != nil {
		i.logger.Errorf("[versionInteractor.stopAndNotify] error stopping version %q: %s", vers.Name, err)
	}

	err = i.versionRepo.SetStatus(ctx, productID, vers.ID, entity.VersionStatusStopped)
	if err != nil {
		i.logger.Errorf("[versionInteractor.stopAndNotify] error stopping version %q: %s", vers.Name, err)
	}

	vers.Status = entity.VersionStatusStopped
	notifyStatusCh <- vers
	i.logger.Infof("[versionInteractor.stopAndNotify] version %q stopped", vers.Name)
}

// Publish set a Version as published on DB and K8s.
func (i *VersionInteractor) Publish(
	ctx context.Context,
	user *entity.User,
	productID,
	versionName,
	comment string,
) (*entity.Version, error) {
	if err := i.accessControl.CheckProductGrants(user, productID, auth.ActPublishVersion); err != nil {
		return nil, err
	}

	i.logger.Infof("The user %s is publishing version %s on product", user.ID, versionName, productID)

	v, err := i.versionRepo.GetByName(ctx, productID, versionName)
	if err != nil {
		return nil, err
	}

	if v.Status != entity.VersionStatusStarted {
		return nil, errors.ErrInvalidVersionStatusBeforePublishing
	}

	err = i.k8sService.Publish(productID, v)
	if err != nil {
		return nil, err
	}

	previousPublishedVersion, err := i.versionRepo.ClearPublishedVersion(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("error unpublishing previous version: %w", err)
	}

	now := time.Now()
	v.PublicationDate = &now
	v.PublicationAuthor = &user.ID
	v.Status = entity.VersionStatusPublished

	err = i.versionRepo.Update(productID, v)
	if err != nil {
		return nil, err
	}

	err = i.userActivityInteractor.RegisterPublishAction(user.ID, productID, v, previousPublishedVersion, comment)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// Unpublish set a Version as not published on DB and K8s.
func (i *VersionInteractor) Unpublish(
	ctx context.Context,
	user *entity.User,
	productID,
	versionName,
	comment string,
) (*entity.Version, error) {
	if err := i.accessControl.CheckProductGrants(user, productID, auth.ActUnpublishVersion); err != nil {
		return nil, err
	}

	i.logger.Infof("The user %s is unpublishing version %s on product %s", user.ID, versionName, productID)

	v, err := i.versionRepo.GetByName(ctx, productID, versionName)
	if err != nil {
		return nil, err
	}

	if v.Status != entity.VersionStatusPublished {
		return nil, errors.ErrInvalidVersionStatusBeforeUnpublishing
	}

	err = i.k8sService.Unpublish(productID, v)
	if err != nil {
		return nil, err
	}

	v.PublicationAuthor = nil
	v.PublicationDate = nil
	v.Status = entity.VersionStatusStarted

	err = i.versionRepo.Update(productID, v)
	if err != nil {
		return nil, err
	}

	err = i.userActivityInteractor.RegisterUnpublishAction(user.ID, productID, v, comment)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (i *VersionInteractor) WatchProcessStatus(
	ctx context.Context,
	user *entity.User,
	productID,
	versionName string,
) (<-chan *entity.Process, error) {
	if err := i.accessControl.CheckProductGrants(user, productID, auth.ActViewProduct); err != nil {
		return nil, err
	}

	v, err := i.versionRepo.GetByName(ctx, productID, versionName)
	if err != nil {
		return nil, err
	}

	return i.k8sService.WatchProcessStatus(ctx, productID, v.Name)
}

func (i *VersionInteractor) WatchProcessLogs(
	ctx context.Context,
	user *entity.User,
	productID,
	versionName string,
	filters entity.LogFilters,
) (<-chan *entity.ProcessLog, error) {
	if err := i.accessControl.CheckProductGrants(user, productID, auth.ActViewVersion); err != nil {
		return nil, err
	}

	return i.processLogRepo.WatchProcessLogs(ctx, productID, versionName, filters)
}

func (i *VersionInteractor) SearchLogs(
	ctx context.Context,
	user *entity.User,
	productID string,
	filters entity.LogFilters,
	cursor *string,
) (*entity.SearchLogsResult, error) {
	if err := i.accessControl.CheckProductGrants(user, productID, auth.ActViewVersion); err != nil {
		return nil, err
	}

	startDate, err := time.Parse(time.RFC3339, filters.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %w", err)
	}

	var endDate time.Time
	if filters.EndDate != nil {
		endDate, err = time.Parse(time.RFC3339, *filters.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date: %w", err)
		}
	} else {
		endDate = time.Now()
	}

	options := &entity.SearchLogsOptions{
		Cursor:         cursor,
		StartDate:      startDate,
		EndDate:        endDate,
		Search:         filters.Search,
		ProcessIDs:     filters.ProcessIDs,
		Levels:         filters.Levels,
		VersionsIDs:    filters.VersionsIDs,
		WorkflowsNames: filters.WorkflowsNames,
	}

	return i.processLogRepo.PaginatedSearch(ctx, productID, options)
}

func (i *VersionInteractor) setStatusError(
	ctx context.Context,
	productID string,
	vers *entity.Version,
	errs []error,
	notifyCh chan *entity.Version,
) {
	errorMessages := make([]string, len(errs))
	for idx, err := range errs {
		errorMessages[idx] = err.Error()
	}

	i.logger.Errorf("The version %q has the following errors: %s", vers.Name,
		strings.Join(errorMessages, "\n"))

	versionWithError, err := i.versionRepo.SetErrors(ctx, productID, vers, errorMessages)
	if err != nil {
		i.logger.Errorf("error saving version error state: %s", err)
	}

	notifyCh <- versionWithError
}
