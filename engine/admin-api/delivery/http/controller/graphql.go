package controller

import (
	"context"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/gql"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/logs"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/version"
	"github.com/labstack/echo"
)

//go:generate mockgen -source=${GOFILE} -destination=../../../mocks/controller_${GOFILE} -package=mocks

const UserContextKey = "user"

type GraphQL interface {
	GraphQLHandler(c echo.Context) error
	PlaygroundHandler(c echo.Context) error
}

type GraphQLController struct {
	logger                 logr.Logger
	productInteractor      *usecase.ProductInteractor
	userInteractor         *usecase.UserInteractor
	userActivityInteractor usecase.UserActivityInteracter
	versionInteractor      *version.Handler
	serverInfoGetter       *usecase.ServerInfoGetter
	processService         *usecase.ProcessService
	LogsUsecase            logs.LogsUsecase
}

type Params struct {
	Logger                 logr.Logger
	ProductInteractor      *usecase.ProductInteractor
	UserInteractor         *usecase.UserInteractor
	UserActivityInteractor usecase.UserActivityInteracter
	VersionInteractor      *version.Handler
	ServerInfoGetter       *usecase.ServerInfoGetter
	ProcessService         *usecase.ProcessService
	LogsUsecase            logs.LogsUsecase
}

func NewGraphQLController(
	params Params,
) *GraphQLController {
	return &GraphQLController{
		params.Logger,
		params.ProductInteractor,
		params.UserInteractor,
		params.UserActivityInteractor,
		params.VersionInteractor,
		params.ServerInfoGetter,
		params.ProcessService,
		params.LogsUsecase,
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
		ProductInteractor:      g.productInteractor,
		UserInteractor:         g.userInteractor,
		UserActivityInteractor: g.userActivityInteractor,
		VersionInteractor:      g.versionInteractor,
		ServerInfoGetter:       g.serverInfoGetter,
		ProcessService:         g.processService,
		LogsUsecase:            g.LogsUsecase,
	})

	h.ServeHTTP(c.Response(), r.WithContext(ctx))

	return nil
}

func (g *GraphQLController) PlaygroundHandler(c echo.Context) error {
	h := playground.Handler("GraphQL playground", "/graphql")
	h.ServeHTTP(c.Response(), c.Request())

	return nil
}
