package controller

import (
	"context"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/logging"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/version"
	"github.com/labstack/echo"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/gql"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase"
)

//go:generate mockgen -source=${GOFILE} -destination=../../../mocks/controller_${GOFILE} -package=mocks

const UserContextKey = "user"

type GraphQL interface {
	GraphQLHandler(c echo.Context) error
	PlaygroundHandler(c echo.Context) error
}

type GraphQLController struct {
	cfg                    *config.Config
	logger                 logging.Logger
	productInteractor      *usecase.ProductInteractor
	userInteractor         *usecase.UserInteractor
	userActivityInteractor usecase.UserActivityInteracter
	versionInteractor      *version.Handler
	metricsInteractor      *usecase.MetricsInteractor
	serverInfoGetter       *usecase.ServerInfoGetter
	processService         *usecase.ProcessService
}

type Params struct {
	Logger                 logging.Logger
	Cfg                    *config.Config
	ProductInteractor      *usecase.ProductInteractor
	UserInteractor         *usecase.UserInteractor
	UserActivityInteractor usecase.UserActivityInteracter
	VersionInteractor      *version.Handler
	MetricsInteractor      *usecase.MetricsInteractor
	ServerInfoGetter       *usecase.ServerInfoGetter
	ProcessService         *usecase.ProcessService
}

func NewGraphQLController(
	params Params,
) *GraphQLController {
	return &GraphQLController{
		params.Cfg,
		params.Logger,
		params.ProductInteractor,
		params.UserInteractor,
		params.UserActivityInteractor,
		params.VersionInteractor,
		params.MetricsInteractor,
		params.ServerInfoGetter,
		params.ProcessService,
	}
}

func (g *GraphQLController) GraphQLHandler(c echo.Context) error {
	r := c.Request()

	ctx := r.Context()

	user, ok := c.Get("user").(*entity.User)
	if ok {
		//nolint:staticcheck // legacy code
		ctx = context.WithValue(ctx, UserContextKey, user)
		g.logger.Info("Request from user " + user.ID)
	}

	h := gql.NewHTTPHandler(gql.Params{
		Logger:                 g.logger,
		Cfg:                    g.cfg,
		ProductInteractor:      g.productInteractor,
		UserInteractor:         g.userInteractor,
		UserActivityInteractor: g.userActivityInteractor,
		VersionInteractor:      g.versionInteractor,
		MetricsInteractor:      g.metricsInteractor,
		ServerInfoGetter:       g.serverInfoGetter,
		ProcessService:         g.processService,
	})

	h.ServeHTTP(c.Response(), r.WithContext(ctx))

	return nil
}

func (g *GraphQLController) PlaygroundHandler(c echo.Context) error {
	h := playground.Handler("GraphQL playground", "/graphql")
	h.ServeHTTP(c.Response(), c.Request())

	return nil
}
