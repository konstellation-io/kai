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

const UserIDContextKey = "userID"
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
}

func NewGraphQLController(
	cfg *config.Config,
	logger logging.Logger,
	runtimeInteractor *usecase.ProductInteractor,
	userInteractor *usecase.UserInteractor,
	userActivityInteractor usecase.UserActivityInteracter,
	versionInteractor *usecase.VersionInteractor,
	metricsInteractor *usecase.MetricsInteractor,
) *GraphQLController {
	return &GraphQLController{
		cfg,
		logger,
		runtimeInteractor,
		userInteractor,
		userActivityInteractor,
		versionInteractor,
		metricsInteractor,
	}
}

func (g *GraphQLController) GraphQLHandler(c echo.Context) error {
	user := c.Get("user").(*entity.User)

	g.logger.Info("Request from user " + user.ID)

	h := gql.NewHTTPHandler(
		g.logger,
		g.runtimeInteractor,
		g.userInteractor,
		g.userActivityInteractor,
		g.versionInteractor,
		g.metricsInteractor,
		g.cfg,
	)

	r := c.Request()

	//nolint:staticcheck // legacy code
	ctx := context.WithValue(r.Context(), UserIDContextKey, user.ID)
	ctx = context.WithValue(ctx, UserContextKey, user)

	h.ServeHTTP(c.Response(), r.WithContext(ctx))

	return nil
}

func (g *GraphQLController) PlaygroundHandler(c echo.Context) error {
	h := playground.Handler("GraphQL playground", "/graphql")
	h.ServeHTTP(c.Response(), c.Request())

	return nil
}
