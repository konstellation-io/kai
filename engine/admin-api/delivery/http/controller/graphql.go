package controller

import (
	"context"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/gql"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/logging"
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
	runtimeInteractor      *usecase.ProductInteractor
	userInteractor         *usecase.UserInteractor
	userActivityInteractor usecase.UserActivityInteracter
	versionInteractor      *usecase.VersionInteractor
	metricsInteractor      *usecase.MetricsInteractor
	serverInfoGetter       *usecase.ServerInfoGetter
}

func NewGraphQLController(
	cfg *config.Config,
	logger logging.Logger,
	runtimeInteractor *usecase.ProductInteractor,
	userInteractor *usecase.UserInteractor,
	userActivityInteractor usecase.UserActivityInteracter,
	versionInteractor *usecase.VersionInteractor,
	metricsInteractor *usecase.MetricsInteractor,
	serverInfoGetter *usecase.ServerInfoGetter,
) *GraphQLController {
	return &GraphQLController{
		cfg,
		logger,
		runtimeInteractor,
		userInteractor,
		userActivityInteractor,
		versionInteractor,
		metricsInteractor,
		serverInfoGetter,
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

	h := gql.NewHTTPHandler(
		g.logger,
		g.runtimeInteractor,
		g.userInteractor,
		g.userActivityInteractor,
		g.versionInteractor,
		g.metricsInteractor,
		g.serverInfoGetter,
		g.cfg,
	)

	h.ServeHTTP(c.Response(), r.WithContext(ctx))

	return nil
}

func (g *GraphQLController) PlaygroundHandler(c echo.Context) error {
	h := playground.Handler("GraphQL playground", "/graphql")
	h.ServeHTTP(c.Response(), c.Request())

	return nil
}
