package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/konstellation-io/kre/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kre/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kre/engine/admin-api/domain/repository"
	"github.com/konstellation-io/kre/engine/admin-api/domain/service"
	"github.com/konstellation-io/kre/engine/admin-api/domain/usecase/auth"
	"github.com/konstellation-io/kre/engine/admin-api/domain/usecase/krt"
	"github.com/konstellation-io/kre/engine/admin-api/domain/usecase/krt/parser"
	"github.com/konstellation-io/kre/engine/admin-api/domain/usecase/krt/validator"
	"github.com/konstellation-io/kre/engine/admin-api/domain/usecase/logging"
	"github.com/konstellation-io/kre/engine/admin-api/domain/usecase/version"
)

var (
	// ErrVersionNotFound error.
	ErrVersionNotFound = errors.New("error version not found")
	// ErrVersionDuplicated error.
	ErrVersionDuplicated = errors.New("error version duplicated")
	// ErrVersionConfigIncomplete error.
	ErrVersionConfigIncomplete = errors.New("version config is incomplete")
	// ErrVersionConfigInvalidKey error.
	ErrVersionConfigInvalidKey = errors.New("version config contains an unknown key")
	// ErrUpdatingStartedVersionConfig error.
	ErrUpdatingStartedVersionConfig = errors.New("config can't be incomplete for started version")
	// ErrInvalidVersionStatusBeforeStarting error.
	ErrInvalidVersionStatusBeforeStarting = errors.New("the version must be stopped before starting")
	// ErrInvalidVersionStatusBeforeStopping error.
	ErrInvalidVersionStatusBeforeStopping = errors.New("the version must be started before stopping")
	// ErrInvalidVersionStatusBeforePublishing error.
	ErrInvalidVersionStatusBeforePublishing = errors.New("the version must be started before publishing")
	// ErrInvalidVersionStatusBeforeUnpublishing error.
	ErrInvalidVersionStatusBeforeUnpublishing = errors.New("the version must be published before unpublishing")
)

// VersionInteractor contains app logic about Version entities.
type VersionInteractor struct {
	cfg                    *config.Config
	logger                 logging.Logger
	versionRepo            repository.VersionRepo
	runtimeRepo            repository.RuntimeRepo
	versionService         service.VersionService
	natsManagerService     service.NatsManagerService
	userActivityInteractor UserActivityInteracter
	accessControl          auth.AccessControl
	idGenerator            version.IDGenerator
	docGenerator           version.DocGenerator
	dashboardService       service.DashboardService
	nodeLogRepo            repository.NodeLogRepository
}

// NewVersionInteractor creates a new interactor.
func NewVersionInteractor(
	cfg *config.Config,
	logger logging.Logger,
	versionRepo repository.VersionRepo,
	runtimeRepo repository.RuntimeRepo,
	versionService service.VersionService,
	natsManagerService service.NatsManagerService,
	userActivityInteractor UserActivityInteracter,
	accessControl auth.AccessControl,
	idGenerator version.IDGenerator,
	docGenerator version.DocGenerator,
	dashboardService service.DashboardService,
	nodeLogRepo repository.NodeLogRepository,
) *VersionInteractor {
	return &VersionInteractor{
		cfg,
		logger,
		versionRepo,
		runtimeRepo,
		versionService,
		natsManagerService,
		userActivityInteractor,
		accessControl,
		idGenerator,
		docGenerator,
		dashboardService,
		nodeLogRepo,
	}
}

func (i *VersionInteractor) filterConfigVars(loggedUserID string, vers *entity.Version) {
	if err := i.accessControl.CheckPermission(loggedUserID, auth.ResVersion, auth.ActEdit); err != nil {
		vers.Config.Vars = nil
	}
}

// GetByName returns a Version by its unique name.
func (i *VersionInteractor) GetByName(ctx context.Context, loggedUserID,
	runtimeID, name string) (*entity.Version, error) {
	v, err := i.versionRepo.GetByName(ctx, runtimeID, name)
	if err != nil {
		return nil, err
	}

	i.filterConfigVars(loggedUserID, v)

	return v, nil
}

func (i *VersionInteractor) GetByID(runtimeID, versionID string) (*entity.Version, error) {
	return i.versionRepo.GetByID(runtimeID, versionID)
}

// GetByRuntime returns all Versions of the given Runtime.
func (i *VersionInteractor) GetByRuntime(loggedUserID, runtimeID string) ([]*entity.Version, error) {
	versions, err := i.versionRepo.GetByRuntime(runtimeID)
	if err != nil {
		return nil, err
	}

	for _, v := range versions {
		i.filterConfigVars(loggedUserID, v)
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
func (i *VersionInteractor) Create(ctx context.Context,
	loggedUserID, runtimeID string, krtFile io.Reader) (*entity.Version, chan *entity.Version, error) {
	if err := i.accessControl.CheckPermission(loggedUserID, auth.ResVersion, auth.ActEdit); err != nil {
		return nil, nil, err
	}

	runtime, err := i.runtimeRepo.GetByID(ctx, runtimeID)
	if err != nil {
		return nil, nil, fmt.Errorf("error runtime repo GetById: %w", err)
	}
	// Check if the version is duplicated
	versions, err := i.versionRepo.GetByRuntime(runtimeID)
	if err != nil {
		return nil, nil, ErrVersionDuplicated
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

	valuesValidator := validator.NewYamlFieldsValidator()

	krtYml, err := parser.ProcessAndValidateKrt(i.logger, valuesValidator, tmpKrtFile.Name(), tmpDir)
	if err != nil {
		return nil, nil, err
	}

	duplicatedVersion, err := i.versionRepo.GetByName(ctx, runtimeID, krtYml.Version)
	if err != nil && !errors.Is(err, ErrVersionNotFound) {
		return nil, nil, fmt.Errorf("error version repo GetByName: %w", err)
	}

	if duplicatedVersion != nil {
		return nil, nil, ErrVersionDuplicated
	}

	workflows, err := i.generateWorkflows(krtYml)
	if err != nil {
		return nil, nil, fmt.Errorf("error generating workflows: %w", err)
	}

	existingConfig := readExistingConf(versions)
	cfg := fillNewConfWithExisting(existingConfig, krtYml)

	krtVersion, ok := entity.ParseKRTVersionFromString(krtYml.KrtVersion)
	if !ok {
		//nolint:goerr113 // errors dynamically generated
		return nil, nil, fmt.Errorf("error krtVersion input from krt.yml not valid: %s", krtYml.KrtVersion)
	}

	versionCreated, err := i.versionRepo.Create(loggedUserID, runtimeID, &entity.Version{
		KrtVersion:  krtVersion,
		Name:        krtYml.Version,
		Description: krtYml.Description,
		Config:      cfg,
		Entrypoint: entity.Entrypoint{
			ProtoFile: krtYml.Entrypoint.Proto,
			Image:     krtYml.Entrypoint.Image,
		},
		Workflows: workflows,
	})
	if err != nil {
		return nil, nil, err
	}

	i.logger.Info("Version created")

	notifyStatusCh := make(chan *entity.Version, 1)

	go i.completeVersionCreation(loggedUserID, tmpKrtFile, krtYml, tmpDir, runtime, versionCreated, notifyStatusCh)

	return versionCreated, notifyStatusCh, nil
}

func (i *VersionInteractor) completeVersionCreation(
	loggedUserID string,
	tmpKrtFile *os.File,
	krtYml *krt.Krt,
	tmpDir string,
	runtime *entity.Runtime,
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

	contentErrors := parser.ProcessContent(i.logger, krtYml, tmpKrtFile.Name(), tmpDir)
	if len(contentErrors) > 0 {
		errorMessage := "error processing krt"
		//nolint:goerr113 // errors dynamically generated
		contentErrors = append([]error{errors.New(errorMessage)}, contentErrors...)
	}

	dashboardsFolder := path.Join(tmpDir, "metrics/dashboards")
	contentErrors = i.saveKRTDashboards(ctx, dashboardsFolder, runtime, versionCreated, contentErrors)

	docFolder := path.Join(tmpDir, "docs")
	contentErrors = i.saveKRTDoc(runtime.ID, docFolder, versionCreated, contentErrors, ctx)

	err := i.versionRepo.UploadKRTFile(runtime.ID, versionCreated, tmpKrtFile.Name())
	if err != nil {
		errorMessage := "error storing KRT file"
		//nolint:goerr113 // errors dynamically generated
		contentErrors = append([]error{errors.New(errorMessage)}, contentErrors...)
	}

	if len(contentErrors) > 0 {
		i.setStatusError(ctx, runtime.ID, versionCreated, contentErrors, notifyStatusCh)
		return
	}

	err = i.versionRepo.SetStatus(ctx, runtime.ID, versionCreated.ID, entity.VersionStatusCreated)
	if err != nil {
		i.logger.Errorf("error setting version status: %s", err)
		return
	}

	// Notify state
	versionCreated.Status = entity.VersionStatusCreated
	notifyStatusCh <- versionCreated

	err = i.userActivityInteractor.RegisterCreateAction(loggedUserID, runtime.ID, versionCreated)
	if err != nil {
		i.logger.Errorf("error registering activity: %s", err)
	}
}

func (i *VersionInteractor) saveKRTDashboards(ctx context.Context, dashboardsFolder string,
	runtime *entity.Runtime, versionCreated *entity.Version, contentErrors []error) []error {
	if _, err := os.Stat(path.Join(dashboardsFolder)); err == nil {
		err := i.storeDashboards(ctx, dashboardsFolder, runtime.ID, versionCreated.Name)
		if err != nil {
			errorMessage := "error creating dashboard"
			//nolint:goerr113 // errors dynamically generated
			contentErrors = append(contentErrors, errors.New(errorMessage))
		}
	}

	return contentErrors
}

func (i *VersionInteractor) saveKRTDoc(runtimeID, docFolder string, versionCreated *entity.Version,
	contentErrors []error, ctx context.Context) []error {
	if _, err := os.Stat(path.Join(docFolder, "README.md")); err == nil {
		err = i.docGenerator.Generate(versionCreated.Name, docFolder)
		if err != nil {
			errorMessage := "error generating version doc"
			//nolint:goerr113 // errors dynamically generated
			contentErrors = append(contentErrors, errors.New(errorMessage))
		}

		err = i.versionRepo.SetHasDoc(ctx, runtimeID, versionCreated.ID, true)
		if err != nil {
			errorMessage := "error updating has doc field"
			//nolint:goerr113 // errors dynamically generated
			contentErrors = append(contentErrors, errors.New(errorMessage))
		}
	} else {
		i.logger.Infof("No documentation found inside the krt files")
	}

	return contentErrors
}

func (i *VersionInteractor) generateWorkflows(krtYml *krt.Krt) ([]*entity.Workflow, error) {
	workflows := make([]*entity.Workflow, 0, len(krtYml.Workflows))
	if len(krtYml.Workflows) == 0 {
		//nolint:goerr113 // errors dynamically generated
		return workflows, fmt.Errorf("error generating workflows: there are no defined workflows")
	}

	for _, w := range krtYml.Workflows {
		var nodes []entity.Node

		if len(w.Nodes) == 0 {
			//nolint:goerr113 // errors dynamically generated
			return nil, fmt.Errorf("error generating workflows: workflow %q doesn't have nodes defined", w.Name)
		}

		for _, node := range w.Nodes {
			replicas := int32(1)
			if node.Replicas != 0 {
				replicas = node.Replicas
			}

			nodeToAppend := entity.Node{
				ID:            i.idGenerator.NewID(),
				Name:          node.Name,
				Image:         node.Image,
				Src:           node.Src,
				GPU:           node.GPU,
				Subscriptions: node.Subscriptions,
				Replicas:      replicas,
			}

			if node.ObjectStore != nil {
				if node.ObjectStore.Scope == "" {
					node.ObjectStore.Scope = krt.ObjectStoreConfigDefaultScope
				}

				nodeToAppend.ObjectStore = &entity.ObjectStore{
					Name:  node.ObjectStore.Name,
					Scope: node.ObjectStore.Scope,
				}
			}

			nodes = append(nodes, nodeToAppend)
		}

		workflows = append(workflows, &entity.Workflow{
			ID:         i.idGenerator.NewID(),
			Name:       w.Name,
			Entrypoint: w.Entrypoint,
			Nodes:      nodes,
			Exitpoint:  w.Exitpoint,
		})
	}

	return workflows, nil
}

// Start create the resources of the given Version.
func (i *VersionInteractor) Start(
	ctx context.Context,
	loggedUserID string,
	runtimeID string,
	versionName string,
	comment string,
) (*entity.Version, chan *entity.Version, error) {
	if err := i.accessControl.CheckPermission(loggedUserID, auth.ResVersion, auth.ActEdit); err != nil {
		return nil, nil, err
	}

	i.logger.Infof("The user %q is starting version %q", loggedUserID, versionName)

	v, err := i.versionRepo.GetByName(ctx, runtimeID, versionName)
	if err != nil {
		return nil, nil, err
	}

	if !v.CanBeStarted() {
		return nil, nil, ErrInvalidVersionStatusBeforeStarting
	}

	if !v.Config.Completed {
		return nil, nil, ErrVersionConfigIncomplete
	}

	notifyStatusCh := make(chan *entity.Version, 1)

	err = i.versionRepo.SetStatus(ctx, runtimeID, v.ID, entity.VersionStatusStarting)
	if err != nil {
		return nil, nil, err
	}

	// Notify intermediate state
	v.Status = entity.VersionStatusStarting
	notifyStatusCh <- v

	err = i.userActivityInteractor.RegisterStartAction(loggedUserID, runtimeID, v, comment)
	if err != nil {
		return nil, nil, err
	}

	versionStreamCfg, err := i.natsManagerService.CreateStreams(ctx, runtimeID, v)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating streams for version %q: %w", v.Name, err)
	}

	objectStoreCfg, err := i.natsManagerService.CreateObjectStores(ctx, runtimeID, v)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating objects stores for version %q: %w", v.Name, err)
	}

	kvStoreCfg, err := i.natsManagerService.CreateKeyValueStores(ctx, runtimeID, v)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating key-value stores for version %q: %w", v.Name, err)
	}

	versionCfg := entity.NewVersionConfig(versionStreamCfg, objectStoreCfg, kvStoreCfg)

	go i.startAndNotify(runtimeID, v, versionCfg, notifyStatusCh)

	return v, notifyStatusCh, nil
}

// Stop removes the resources of the given Version.
func (i *VersionInteractor) Stop(
	ctx context.Context,
	loggedUserID string,
	runtimeID string,
	versionName string,
	comment string,
) (*entity.Version, chan *entity.Version, error) {
	if err := i.accessControl.CheckPermission(loggedUserID, auth.ResVersion, auth.ActEdit); err != nil {
		return nil, nil, err
	}

	i.logger.Infof("The user %s is stopping version %s on runtime %s", loggedUserID, versionName, runtimeID)

	v, err := i.versionRepo.GetByName(ctx, runtimeID, versionName)
	if err != nil {
		return nil, nil, err
	}

	if !v.CanBeStopped() {
		return nil, nil, ErrInvalidVersionStatusBeforeStopping
	}

	err = i.versionRepo.SetStatus(ctx, runtimeID, v.ID, entity.VersionStatusStopping)
	if err != nil {
		return nil, nil, err
	}

	notifyStatusCh := make(chan *entity.Version, 1)

	// Notify intermediate state
	v.Status = entity.VersionStatusStopping
	notifyStatusCh <- v

	err = i.userActivityInteractor.RegisterStopAction(loggedUserID, runtimeID, v, comment)
	if err != nil {
		return nil, nil, err
	}

	err = i.natsManagerService.DeleteStreams(ctx, runtimeID, versionName)
	if err != nil {
		return nil, nil, fmt.Errorf("error stopping version %q: %w", versionName, err)
	}

	err = i.natsManagerService.DeleteObjectStores(ctx, runtimeID, versionName)
	if err != nil {
		return nil, nil, fmt.Errorf("error stopping version %q: %w", versionName, err)
	}

	go i.stopAndNotify(runtimeID, v, notifyStatusCh)

	return v, notifyStatusCh, nil
}

func (i *VersionInteractor) startAndNotify(
	runtimeID string,
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

	err := i.versionService.Start(ctx, runtimeID, vers, versionConfig)
	if err != nil {
		i.logger.Errorf("[versionInteractor.startAndNotify] error starting version %q: %s", vers.Name, err)
	}

	err = i.versionRepo.SetStatus(ctx, runtimeID, vers.ID, entity.VersionStatusStarted)
	if err != nil {
		i.logger.Errorf("[versionInteractor.startAndNotify] error starting version %q: %s", vers.Name, err)
	}

	vers.Status = entity.VersionStatusStarted
	notifyStatusCh <- vers
	i.logger.Infof("[versionInteractor.startAndNotify] version %q started", vers.Name)
}

func (i *VersionInteractor) stopAndNotify(
	runtimeID string,
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

	err := i.versionService.Stop(ctx, runtimeID, vers)
	if err != nil {
		i.logger.Errorf("[versionInteractor.stopAndNotify] error stopping version %q: %s", vers.Name, err)
	}

	err = i.versionRepo.SetStatus(ctx, runtimeID, vers.ID, entity.VersionStatusStopped)
	if err != nil {
		i.logger.Errorf("[versionInteractor.stopAndNotify] error stopping version %q: %s", vers.Name, err)
	}

	vers.Status = entity.VersionStatusStopped
	notifyStatusCh <- vers
	i.logger.Infof("[versionInteractor.stopAndNotify] version %q stopped", vers.Name)
}

// Publish set a Version as published on DB and K8s.
func (i *VersionInteractor) Publish(ctx context.Context, loggedUserID, runtimeID,
	versionName, comment string) (*entity.Version, error) {
	if err := i.accessControl.CheckPermission(loggedUserID, auth.ResVersion, auth.ActEdit); err != nil {
		return nil, err
	}

	i.logger.Infof("The user %s is publishing version %s", loggedUserID, versionName)

	v, err := i.versionRepo.GetByName(ctx, runtimeID, versionName)
	if err != nil {
		return nil, err
	}

	if v.Status != entity.VersionStatusStarted {
		return nil, ErrInvalidVersionStatusBeforePublishing
	}

	err = i.versionService.Publish(runtimeID, v)
	if err != nil {
		return nil, err
	}

	previousPublishedVersion, err := i.versionRepo.ClearPublishedVersion(ctx, runtimeID)
	if err != nil {
		return nil, fmt.Errorf("error unpublishing previous version: %w", err)
	}

	now := time.Now()
	v.PublicationDate = &now
	v.PublicationUserID = &loggedUserID
	v.Status = entity.VersionStatusPublished

	err = i.versionRepo.Update(runtimeID, v)
	if err != nil {
		return nil, err
	}

	err = i.userActivityInteractor.RegisterPublishAction(loggedUserID, runtimeID, v, previousPublishedVersion, comment)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// Unpublish set a Version as not published on DB and K8s.
func (i *VersionInteractor) Unpublish(ctx context.Context, loggedUserID, runtimeID,
	versionName, comment string) (*entity.Version, error) {
	if err := i.accessControl.CheckPermission(loggedUserID, auth.ResVersion, auth.ActEdit); err != nil {
		return nil, err
	}

	i.logger.Infof("The user %s is unpublishing version %s", loggedUserID, versionName)

	v, err := i.versionRepo.GetByName(ctx, runtimeID, versionName)
	if err != nil {
		return nil, err
	}

	if v.Status != entity.VersionStatusPublished {
		return nil, ErrInvalidVersionStatusBeforeUnpublishing
	}

	err = i.versionService.Unpublish(runtimeID, v)
	if err != nil {
		return nil, err
	}

	v.PublicationUserID = nil
	v.PublicationDate = nil
	v.Status = entity.VersionStatusStarted

	err = i.versionRepo.Update(runtimeID, v)
	if err != nil {
		return nil, err
	}

	err = i.userActivityInteractor.RegisterUnpublishAction(loggedUserID, runtimeID, v, comment)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (i *VersionInteractor) UpdateVersionConfig(ctx context.Context, loggedUserID, runtimeID string,
	vrs *entity.Version, conf []*entity.ConfigurationVariable) (*entity.Version, error) {
	if err := i.accessControl.CheckPermission(loggedUserID, auth.ResVersion, auth.ActEdit); err != nil {
		return nil, err
	}

	err := i.validateNewConfig(vrs.Config.Vars, conf)
	if err != nil {
		return nil, err
	}

	isStarted := vrs.PublishedOrStarted()

	newConfig, newConfigIsComplete := generateNewConfig(vrs.Config.Vars, conf)

	if isStarted && !newConfigIsComplete {
		return nil, ErrUpdatingStartedVersionConfig
	}

	vrs.Config.Vars = newConfig
	vrs.Config.Completed = newConfigIsComplete

	// No need to restart PODs if there are no resources running
	if isStarted {
		err = i.versionService.UpdateConfig(runtimeID, vrs)
		if err != nil {
			return nil, err
		}
	}

	err = i.versionRepo.Update(runtimeID, vrs)
	if err != nil {
		return nil, err
	}

	return vrs, nil
}

func (i *VersionInteractor) WatchNodeStatus(ctx context.Context, loggedUserID,
	runtimeID, versionName string) (<-chan *entity.Node, error) {
	if err := i.accessControl.CheckPermission(loggedUserID, auth.ResVersion, auth.ActView); err != nil {
		return nil, err
	}

	v, err := i.versionRepo.GetByName(ctx, runtimeID, versionName)
	if err != nil {
		return nil, err
	}

	return i.versionService.WatchNodeStatus(ctx, runtimeID, v.Name)
}

func (i *VersionInteractor) WatchNodeLogs(
	ctx context.Context,
	loggedUserID, runtimeID, versionName string,
	filters entity.LogFilters,
) (<-chan *entity.NodeLog, error) {
	if err := i.accessControl.CheckPermission(loggedUserID, auth.ResLogs, auth.ActView); err != nil {
		return nil, err
	}

	return i.nodeLogRepo.WatchNodeLogs(ctx, runtimeID, versionName, filters)
}

func (i *VersionInteractor) SearchLogs(
	ctx context.Context,
	loggedUserID,
	runtimeID string,
	filters entity.LogFilters,
	cursor *string,
) (*entity.SearchLogsResult, error) {
	if err := i.accessControl.CheckPermission(loggedUserID, auth.ResLogs, auth.ActView); err != nil {
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
		NodeIDs:        filters.NodeIDs,
		Levels:         filters.Levels,
		VersionsIDs:    filters.VersionsIDs,
		WorkflowsNames: filters.WorkflowsNames,
	}

	return i.nodeLogRepo.PaginatedSearch(ctx, runtimeID, options)
}

func (i *VersionInteractor) setStatusError(
	ctx context.Context,
	runtimeID string,
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

	versionWithError, err := i.versionRepo.SetErrors(ctx, runtimeID, vers, errorMessages)
	if err != nil {
		i.logger.Errorf("error saving version error state: %s", err)
	}

	notifyCh <- versionWithError
}
