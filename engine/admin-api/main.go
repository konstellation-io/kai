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
	"github.com/konstellation-io/kai/engine/admin-api/adapter/repository/objectstorage"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/natsmanager"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/natspb"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/versionpb"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/user"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/versionservice"
	"github.com/konstellation-io/kai/engine/admin-api/delivery/http"
	"github.com/konstellation-io/kai/engine/admin-api/delivery/http/controller"
	logging2 "github.com/konstellation-io/kai/engine/admin-api/domain/service/logging"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/version"
	"github.com/sethvargo/go-password/password"
	"github.com/spf13/viper"
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

//nolint:funlen //future refactor
func initGraphqlController(
	cfg *config.Config, oldLogger logging2.Logger, logger logr.Logger, mongodbClient *mongo.Client,
) *controller.GraphQLController {
	productRepo := mongodb.NewProductRepoMongoDB(logger, mongodbClient)
	userActivityRepo := mongodb.NewUserActivityRepoMongoDB(cfg, oldLogger, mongodbClient)
	versionMongoRepo := versionrepository.New(cfg, oldLogger, mongodbClient)
	processLogRepo := mongodb.NewProcessLogMongoDBRepo(cfg, oldLogger, mongodbClient)
	processRepo := processrepository.New(cfg, oldLogger, mongodbClient)
	metricRepo := mongodb.NewMetricMongoDBRepo(cfg, oldLogger, mongodbClient)

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

	natsManagerService, err := natsmanager.NewClient(cfg, oldLogger, natsManagerClient)
	if err != nil {
		log.Fatal(err)
	}

	accessControl, err := casbinauth.NewCasbinAccessControl(oldLogger, "./casbin_rbac_model.conf", "./casbin_rbac_policy.csv")
	if err != nil {
		log.Fatal(err)
	}

	keycloakUserRegistry, err := user.NewKeycloakUserRegistry(user.WithClient(viper.GetString(config.KeycloakURLKey)))
	if err != nil {
		log.Fatal(err)
	}

	userActivityInteractor := usecase.NewUserActivityInteractor(logger, userActivityRepo, accessControl)

	minioClient, err := objectstorage.NewMinioClient()
	if err != nil {
		log.Fatal(err)
	}

	minioAdminClient, err := objectstorage.NewAdminMinioClient()
	if err != nil {
		log.Fatal(err)
	}

	passwordGenerator, err := password.NewGenerator(&password.GeneratorInput{})
	if err != nil {
		log.Fatal(err)
	}

	minioOjectStorage := objectstorage.NewMinioObjectStorage(logger, minioClient, minioAdminClient)

	productInteractor := usecase.NewProductInteractor(&usecase.ProductInteractorOpts{
		Logger:            logger,
		ProductRepo:       productRepo,
		MeasurementRepo:   influx.NewMeasurementRepoInfluxDB(cfg, oldLogger),
		VersionRepo:       versionMongoRepo,
		MetricRepo:        metricRepo,
		ProcessLogRepo:    processLogRepo,
		ProcessRepo:       processRepo,
		UserActivity:      userActivityInteractor,
		AccessControl:     accessControl,
		ObjectStorage:     minioOjectStorage,
		NatsService:       natsManagerService,
		UserRegistry:      keycloakUserRegistry,
		PasswordGenerator: passwordGenerator,
	})

	userInteractor := usecase.NewUserInteractor(
		oldLogger,
		accessControl,
		userActivityInteractor,
		keycloakUserRegistry,
	)

	versionInteractor := version.NewHandler(
		&version.HandlerParams{
			Logger:                 logger,
			VersionRepo:            versionMongoRepo,
			ProductRepo:            productRepo,
			K8sService:             k8sService,
			NatsManagerService:     natsManagerService,
			UserActivityInteractor: userActivityInteractor,
			AccessControl:          accessControl,
			ProcessLogRepo:         processLogRepo,
		},
	)

	metricsInteractor := usecase.NewMetricsInteractor(
		oldLogger,
		productRepo,
		accessControl,
		metricRepo,
	)

	serverInfoGetter := usecase.NewServerInfoGetter(logger, accessControl)

	processService := usecase.NewProcessService(logger, k8sService, processRepo, minioOjectStorage)

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
