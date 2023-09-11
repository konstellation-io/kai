package main

import (
	"log"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/casbinauth"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/repository/influx"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/repository/mongodb"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/repository/mongodb/processrepository"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/repository/mongodb/versionrepository"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/natsmanager"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/natspb"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/versionpb"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/versionservice"
	"github.com/konstellation-io/kai/engine/admin-api/delivery/http"
	"github.com/konstellation-io/kai/engine/admin-api/delivery/http/controller"
	logging2 "github.com/konstellation-io/kai/engine/admin-api/domain/service/logging"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/version"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.NewConfig()

	err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	oldLogger := logging2.NewLogger(cfg.LogLevel)

	zapLog, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}

	logger := zapr.NewLogger(zapLog)

	db := mongodb.NewMongoDB(cfg, oldLogger)

	mongodbClient := db.Connect()
	defer db.Disconnect()

	graphqlController := initGraphqlController(cfg, oldLogger, logger, mongodbClient)

	app := http.NewApp(
		cfg,
		oldLogger,
		graphqlController,
	)

	app.Start()
}

func initGraphqlController(
	cfg *config.Config, oldLogger logging2.Logger, logger logr.Logger, mongodbClient *mongo.Client,
) *controller.GraphQLController {
	productRepo := mongodb.NewProductRepoMongoDB(cfg, oldLogger, mongodbClient)
	userActivityRepo := mongodb.NewUserActivityRepoMongoDB(cfg, oldLogger, mongodbClient)
	versionMongoRepo := versionrepository.New(cfg, oldLogger, mongodbClient)
	processLogRepo := mongodb.NewProcessLogMongoDBRepo(cfg, oldLogger, mongodbClient)
	processRepo := processrepository.New(cfg, oldLogger, mongodbClient)
	metricRepo := mongodb.NewMetricMongoDBRepo(cfg, oldLogger, mongodbClient)
	measurementRepo := influx.NewMeasurementRepoInfluxDB(cfg, oldLogger)

	ccK8sManager, err := grpc.Dial(cfg.Services.K8sManager, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	k8sManagerClient := versionpb.NewVersionServiceClient(ccK8sManager)

	k8sService, err := versionservice.New(cfg, oldLogger, k8sManagerClient)
	if err != nil {
		log.Fatal(err)
	}

	ccNatsManager, err := grpc.Dial(cfg.Services.NatsManager, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	natsManagerClient := natspb.NewNatsManagerServiceClient(ccNatsManager)

	natsManagerService, err := natsmanager.NewNatsManagerClient(cfg, oldLogger, natsManagerClient)
	if err != nil {
		log.Fatal(err)
	}

	accessControl, err := casbinauth.NewCasbinAccessControl(oldLogger, "./casbin_rbac_model.conf", "./casbin_rbac_policy.csv")
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
		Logger:          logger,
		ProductRepo:     productRepo,
		MeasurementRepo: measurementRepo,
		VersionRepo:     versionMongoRepo,
		MetricRepo:      metricRepo,
		ProcessLogRepo:  processLogRepo,
		ProcessRepo:     processRepo,
		UserActivity:    userActivityInteractor,
		AccessControl:   accessControl,
	}
	productInteractor := usecase.NewProductInteractor(&ps)

	userInteractor := usecase.NewUserInteractor(
		oldLogger,
		accessControl,
		userActivityInteractor,
		gocloakUserRegistry,
	)

	chronografDashboard := service.CreateDashboardService(cfg, oldLogger)

	handlerParams := version.HanlderParams{
		Logger:                 logger,
		VersionRepo:            versionMongoRepo,
		ProductRepo:            productRepo,
		K8sService:             k8sService,
		NatsManagerService:     natsManagerService,
		UserActivityInteractor: userActivityInteractor,
		AccessControl:          accessControl,
		DashboardService:       chronografDashboard,
		ProcessLogRepo:         processLogRepo,
	}
	versionInteractor := version.NewHandler(handlerParams)

	metricsInteractor := usecase.NewMetricsInteractor(
		oldLogger,
		productRepo,
		accessControl,
		metricRepo,
	)

	serverInfoGetter := usecase.NewServerInfoGetter(logger, accessControl)

	processService := usecase.NewProcessService(logger, k8sService, processRepo)

	return controller.NewGraphQLController(
		controller.Params{
			Logger:                 oldLogger,
			Cfg:                    cfg,
			ProductInteractor:      productInteractor,
			UserInteractor:         userInteractor,
			UserActivityInteractor: userActivityInteractor,
			VersionInteractor:      versionInteractor,
			MetricsInteractor:      metricsInteractor,
			ServerInfoGetter:       serverInfoGetter,
			ProcessService:         processService,
		},
	)
}
