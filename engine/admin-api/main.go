package main

import (
	"log"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/auth"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/repository/influx"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/repository/mongodb"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/version"
	"github.com/konstellation-io/kai/engine/admin-api/delivery/http"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/logging"
)

func main() {
	cfg := config.NewConfig()
	logger := logging.NewLogger(cfg.LogLevel)

	db := mongodb.NewMongoDB(cfg, logger)

	mongodbClient := db.Connect()
	defer db.Disconnect()

	userActivityInteractor, runtimeInteractor, userInteractor,
		versionInteractor, metricsInteractor := initApp(cfg, logger, mongodbClient)

	app := http.NewApp(
		cfg,
		logger,
		runtimeInteractor,
		userInteractor,
		userActivityInteractor,
		versionInteractor,
		metricsInteractor,
	)
	app.Start()
}

// nolint:godox // this is a TODO to be removed
//
// TODO: Change these to be interfaces
func initApp(cfg *config.Config, logger logging.Logger, mongodbClient *mongo.Client) (usecase.UserActivityInteracter,
	*usecase.RuntimeInteractor, usecase.UserInteractorUsecase, *usecase.VersionInteractor, *usecase.MetricsInteractor) {
	runtimeRepo := mongodb.NewRuntimeRepoMongoDB(cfg, logger, mongodbClient)
	userActivityRepo := mongodb.NewUserActivityRepoMongoDB(cfg, logger, mongodbClient)
	versionMongoRepo := mongodb.NewVersionRepoMongoDB(cfg, logger, mongodbClient)
	measurementRepo := influx.NewMeasurementRepoInfluxDB(cfg, logger)
	metricRepo := mongodb.NewMetricMongoDBRepo(cfg, logger, mongodbClient)
	nodeLogRepo := mongodb.NewNodeLogMongoDBRepo(cfg, logger, mongodbClient)

	versionService, err := service.NewK8sVersionClient(cfg, logger)
	if err != nil {
		log.Fatal(err)
	}

	natsManagerService, err := service.NewNatsManagerClient(cfg, logger)
	if err != nil {
		log.Fatal(err)
	}

	gocloakService, err := service.NewGocloakManager(cfg)
	if err != nil {
		log.Fatal(err)
	}

	accessControl, err := auth.NewCasbinAccessControl(logger, "./casbin_rbac_model.conf", "./casbin_rbac_policy.csv")
	if err != nil {
		log.Fatal(err)
	}

	docGenerator := version.NewHTTPStaticDocGenerator(cfg, logger)
	idGenerator := version.NewIDGenerator()

	userActivityInteractor := usecase.NewUserActivityInteractor(logger, userActivityRepo, accessControl)

	runtimeInteractor := usecase.NewRuntimeInteractor(
		cfg,
		logger,
		runtimeRepo,
		measurementRepo,
		versionMongoRepo,
		metricRepo,
		nodeLogRepo,
		userActivityInteractor,
		accessControl,
	)

	userInteractor := usecase.NewUserInteractor(
		logger,
		userActivityInteractor,
		gocloakService,
	)

	chronografDashboard := service.CreateDashboardService(cfg, logger)
	versionInteractor := usecase.NewVersionInteractor(
		cfg,
		logger,
		versionMongoRepo,
		runtimeRepo,
		versionService,
		natsManagerService,
		userActivityInteractor,
		accessControl,
		idGenerator,
		docGenerator,
		chronografDashboard,
		nodeLogRepo,
	)

	metricsInteractor := usecase.NewMetricsInteractor(
		logger,
		runtimeRepo,
		accessControl,
		metricRepo,
	)

	return userActivityInteractor, runtimeInteractor, userInteractor, versionInteractor, metricsInteractor
}
