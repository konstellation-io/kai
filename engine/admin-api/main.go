package main

import (
	"log"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/auth"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/repository/influx"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/repository/mongodb"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/repository/mongodb/version"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/k8smanager"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/natsmanager"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/natspb"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/versionpb"
	"github.com/konstellation-io/kai/engine/admin-api/delivery/http"
	"github.com/konstellation-io/kai/engine/admin-api/delivery/http/controller"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/logging"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func initApp(
	cfg *config.Config,
	logger logging.Logger,
	mongodbClient *mongo.Client,
) (
	usecase.UserActivityInteracter,
	*usecase.ProductInteractor,
	*usecase.UserInteractor,
	*usecase.VersionInteractor,
	*usecase.MetricsInteractor,
) {
	productRepo := mongodb.NewProductRepoMongoDB(cfg, logger, mongodbClient)
	userActivityRepo := mongodb.NewUserActivityRepoMongoDB(cfg, logger, mongodbClient)
	versionMongoRepo := version.NewVersionRepoMongoDB(cfg, logger, mongodbClient)
	processLogRepo := mongodb.NewProcessLogMongoDBRepo(cfg, logger, mongodbClient)
	metricRepo := mongodb.NewMetricMongoDBRepo(cfg, logger, mongodbClient)
	measurementRepo := influx.NewMeasurementRepoInfluxDB(cfg, logger)

	ccK8sManager, err := grpc.Dial(cfg.Services.K8sManager, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	k8sManagerClient := versionpb.NewVersionServiceClient(ccK8sManager)

	k8sService, err := k8smanager.NewK8sVersionClient(cfg, logger, k8sManagerClient)
	if err != nil {
		log.Fatal(err)
	}

	ccNatsManager, err := grpc.Dial(cfg.Services.NatsManager, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)

	}

	natsManagerClient := natspb.NewNatsManagerServiceClient(ccNatsManager)

	natsManagerService, err := natsmanager.NewNatsManagerClient(cfg, logger, natsManagerClient)
	if err != nil {
		log.Fatal(err)
	}

	accessControl, err := auth.NewCasbinAccessControl(logger, "./casbin_rbac_model.conf", "./casbin_rbac_policy.csv")
	if err != nil {
		log.Fatal(err)
	}

	keycloakCfg := service.KeycloakConfig{
		Realm:         cfg.Keycloak.Realm,
		MasterRealm:   cfg.Keycloak.MasterRealm,
		AdminUsername: cfg.Keycloak.AdminUsername,
		AdminPassword: cfg.Keycloak.AdminPassword,
		AdminClientID: cfg.Keycloak.AdminClientID,
	}

	gocloakUserRegistry, err := service.NewGocloakUserRegistry(service.WithClient(cfg.Keycloak.URL), &keycloakCfg)
	if err != nil {
		log.Fatal(err)
	}

	userActivityInteractor := usecase.NewUserActivityInteractor(logger, userActivityRepo, accessControl)

	ps := usecase.ProductInteractorOpts{
		Cfg:             cfg,
		Logger:          logger,
		ProductRepo:     productRepo,
		MeasurementRepo: measurementRepo,
		VersionRepo:     versionMongoRepo,
		MetricRepo:      metricRepo,
		ProcessLogRepo:  processLogRepo,
		UserActivity:    userActivityInteractor,
		AccessControl:   accessControl,
	}
	productInteractor := usecase.NewProductInteractor(&ps)

	userInteractor := usecase.NewUserInteractor(
		logger,
		accessControl,
		userActivityInteractor,
		gocloakUserRegistry,
	)

	chronografDashboard := service.CreateDashboardService(cfg, logger)
	versionInteractor := usecase.NewVersionInteractor(
		cfg,
		logger,
		versionMongoRepo,
		productRepo,
		k8sService,
		natsManagerService,
		userActivityInteractor,
		accessControl,
		chronografDashboard,
		processLogRepo,
	)

	metricsInteractor := usecase.NewMetricsInteractor(
		logger,
		productRepo,
		accessControl,
		metricRepo,
	)

	return userActivityInteractor, productInteractor, userInteractor, versionInteractor, metricsInteractor
}
