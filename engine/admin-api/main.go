package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/casbinauth"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/repository/mongodb"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/repository/mongodb/processrepository"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/repository/mongodb/versionrepository"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/repository/objectstorage"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/repository/redis"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/loki"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/natsmanager"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/natspb"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/versionpb"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/registry"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/user"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/versionservice"
	"github.com/konstellation-io/kai/engine/admin-api/delivery/http"
	"github.com/konstellation-io/kai/engine/admin-api/delivery/http/controller"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/logs"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/process"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/version"
	"github.com/minio/minio-go/v7"
	"github.com/sethvargo/go-password/password"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	zapLog, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}

	logger := zapr.NewLogger(zapLog)

	db := mongodb.NewMongoDB(logger)

	mongodbClient, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Disconnect()

	graphqlController := initGraphqlController(logger, mongodbClient)

	app := http.NewApp(
		logger,
		graphqlController,
	)

	app.Start()
}

//nolint:funlen // Future refactor
func initGraphqlController(logger logr.Logger, mongodbClient *mongo.Client) *controller.GraphQLController {
	productRepo := mongodb.NewProductRepoMongoDB(logger, mongodbClient)
	userActivityRepo := mongodb.NewUserActivityRepoMongoDB(logger, mongodbClient)
	versionMongoRepo := versionrepository.New(logger, mongodbClient)
	processRepo := processrepository.New(logger, mongodbClient)
	logsService := loki.NewClient()

	ccK8sManager, err := grpc.Dial(
		viper.GetString(config.K8sManagerEndpointKey),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}

	k8sManagerClient := versionpb.NewVersionServiceClient(ccK8sManager)

	k8sService, err := versionservice.New(logger, k8sManagerClient)
	if err != nil {
		log.Fatal(err)
	}

	ccNatsManager, err := grpc.Dial(
		viper.GetString(config.NatsManagerEndpointKey),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}

	natsManagerClient := natspb.NewNatsManagerServiceClient(ccNatsManager)

	natsManagerService, err := natsmanager.NewClient(logger, natsManagerClient)
	if err != nil {
		log.Fatal(err)
	}

	accessControl, err := casbinauth.NewCasbinAccessControl(logger, "./casbin_rbac_model.conf", "./casbin_rbac_policy.csv")
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

	err = ensureKAIBucketExists(logger, minioOjectStorage, minioClient, keycloakUserRegistry)
	if err != nil {
		log.Fatal(err)
	}

	predictionRepo := redis.NewPredictionRepository(redis.NewRedisClient())

	err = predictionRepo.EnsurePredictionIndexCreated()
	if err != nil {
		log.Fatal(err)
	}

	productInteractor := usecase.NewProductInteractor(&usecase.ProductInteractorOpts{
		Logger:               logger,
		ProductRepo:          productRepo,
		VersionRepo:          versionMongoRepo,
		ProcessRepo:          processRepo,
		UserActivity:         userActivityInteractor,
		AccessControl:        accessControl,
		ObjectStorage:        minioOjectStorage,
		NatsService:          natsManagerService,
		UserRegistry:         keycloakUserRegistry,
		PasswordGenerator:    passwordGenerator,
		PredictionRepository: predictionRepo,
	})

	userInteractor := usecase.NewUserInteractor(
		logger,
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
		},
	)

	processHandler := process.NewHandler(
		&process.HandlerParams{
			Logger:            logger,
			VersionService:    k8sService,
			ProcessRepository: processRepo,
			ObjectStorage:     minioOjectStorage,
			AccessControl:     accessControl,
			ProcessRegistry:   registry.NewProcessRegistry(),
			ProductRepository: productRepo,
		},
	)

	logsUseCase := logs.NewLogsInteractor(logsService)

	return controller.NewGraphQLController(
		controller.Params{
			Logger:                 logger,
			ProductInteractor:      productInteractor,
			UserInteractor:         userInteractor,
			UserActivityInteractor: userActivityInteractor,
			VersionInteractor:      versionInteractor,
			ProcessHandler:         processHandler,
			LogsUsecase:            logsUseCase,
		},
	)
}

func ensureKAIBucketExists(
	logger logr.Logger,
	storage *objectstorage.MinioObjectStorage,
	minoClient *minio.Client,
	userRegistry *user.KeycloakUserRegistry,
) error {
	ctx := context.Background()
	kaiBucket := viper.GetString(config.GlobalRegistryKey)

	exists, err := minoClient.BucketExists(ctx, kaiBucket)
	if err != nil {
		return fmt.Errorf("checking if kai bucket exists: %w", err)
	}

	var policyName string

	if !exists {
		logger.Info("Creating KAI bucket in object storage", "bucket", kaiBucket)

		err = storage.CreateBucket(ctx, kaiBucket)
		if err != nil {
			return fmt.Errorf("creating kai bucket: %w", err)
		}

		policyName, err = storage.CreateBucketPolicy(ctx, kaiBucket)
		if err != nil {
			return fmt.Errorf("creating kai object storage policy: %w", err)
		}
	}

	groupExists, err := userRegistry.GroupExists(ctx, kaiBucket)
	if err != nil {
		return fmt.Errorf("checking if kai group exists in user registry: %w", err)
	}

	if !groupExists {
		logger.Info("Creating KAI group in user registry", "group", kaiBucket)

		err = userRegistry.CreateGroupWithPolicy(ctx, kaiBucket, policyName)
		if err != nil {
			return fmt.Errorf("creating kai group in user registry: %w", err)
		}
	}

	return nil
}
