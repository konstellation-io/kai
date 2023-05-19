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
	"github.com/konstellation-io/kai/engine/admin-api/delivery/http/controller"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/logging"
)

func main() {
	cfg := config.NewConfig()
	logger := logging.NewLogger(cfg.LogLevel)

	db := mongodb.NewMongoDB(cfg, logger)

	mongodbClient := db.Connect()
	defer db.Disconnect()

	userActivityInteractor, productInteractor, userInteractor,
		versionInteractor, metricsInteractor := initApp(cfg, logger, mongodbClient)

	graphqlController := controller.NewGraphQLController(
		cfg,
		logger,
		productInteractor,
		userInteractor,
		userActivityInteractor,
		versionInteractor,
		metricsInteractor,
	)

	app := http.NewApp(
		cfg,
		logger,
		graphqlController,
	)

	app.Start()
}

func initApp(cfg *config.Config, logger logging.Logger, mongodbClient *mongo.Client) (usecase.UserActivityInteracter,
	*usecase.ProductInteractor, *usecase.UserInteractor, *usecase.VersionInteractor, *usecase.MetricsInteractor) {
	productRepo := mongodb.NewProductRepoMongoDB(cfg, logger, mongodbClient)
	userActivityRepo := mongodb.NewUserActivityRepoMongoDB(cfg, logger, mongodbClient)
	versionMongoRepo := mongodb.NewVersionRepoMongoDB(cfg, logger, mongodbClient)
	nodeLogRepo := mongodb.NewNodeLogMongoDBRepo(cfg, logger, mongodbClient)
	metricRepo := mongodb.NewMetricMongoDBRepo(cfg, logger, mongodbClient)
	measurementRepo := influx.NewMeasurementRepoInfluxDB(cfg, logger)

	versionService, err := service.NewK8sVersionClient(cfg, logger)
	if err != nil {
		log.Fatal(err)
	}

	natsManagerService, err := service.NewNatsManagerClient(cfg, logger)
	if err != nil {
		log.Fatal(err)
	}

	accessControl, err := auth.NewCasbinAccessControl(logger, "./casbin_rbac_model.conf", "./casbin_rbac_policy.csv")
	if err != nil {
		log.Fatal(err)
	}

	gocloakService, err := service.NewGocloakManager(cfg)
	if err != nil {
		log.Fatal(err)
	}

	docGenerator := version.NewHTTPStaticDocGenerator(cfg, logger)
	idGenerator := version.NewIDGenerator()

	userActivityInteractor := usecase.NewUserActivityInteractor(logger, userActivityRepo, accessControl)

	ps := usecase.ProductInteractorSetup{
		Cfg:             cfg,
		Logger:          logger,
		ProductRepo:     productRepo,
		MeasurementRepo: measurementRepo,
		VersionRepo:     versionMongoRepo,
		MetricRepo:      metricRepo,
		NodeLogRepo:     nodeLogRepo,
		UserActivity:    userActivityInteractor,
		AccessControl:   accessControl,
	}
	productInteractor := usecase.NewProductInteractor(&ps)

	userInteractor := usecase.NewUserInteractor(
		logger,
		accessControl,
		userActivityInteractor,
		gocloakService,
	)

	chronografDashboard := service.CreateDashboardService(cfg, logger)
	versionInteractor := usecase.NewVersionInteractor(
		cfg,
		logger,
		versionMongoRepo,
		productRepo,
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
		productRepo,
		accessControl,
		metricRepo,
	)

	return userActivityInteractor, productInteractor, userInteractor, versionInteractor, metricsInteractor
}
