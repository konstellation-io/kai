package gql

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gorilla/websocket"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/auth"
	internalerrors "github.com/konstellation-io/kai/engine/admin-api/domain/usecase/errors"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/logging"
)

func NewHTTPHandler(
	logger logging.Logger,
	runtimeInteractor *usecase.ProductInteractor,
	userInteractor *usecase.UserInteractor,
	userActivityInteractor usecase.UserActivityInteracter,
	versionInteractor *usecase.VersionInteractor,
	metricsInteractor *usecase.MetricsInteractor,
	serverInfoGetter *usecase.ServerInfoGetter,
	cfg *config.Config,
) http.Handler {
	graphQLResolver := NewGraphQLResolver(
		logger,
		runtimeInteractor,
		userInteractor,
		userActivityInteractor,
		versionInteractor,
		metricsInteractor,
		serverInfoGetter,
		cfg,
	)

	var mb int64 = 1 << 20
	maxUploadSize := 500 * mb
	maxMemory := 500 * mb

	srv := handler.New(NewExecutableSchema(Config{Resolvers: graphQLResolver}))

	srv.SetErrorPresenter(errorPresenter)

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.Use(extension.Introspection{})

	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})

	srv.AddTransport(transport.MultipartForm{
		MaxUploadSize: maxUploadSize,
		MaxMemory:     maxMemory,
	})

	return srv
}

func errorPresenter(ctx context.Context, e error) *gqlerror.Error {
	err := graphql.DefaultErrorPresenter(ctx, e)

	var errInvalidKRT internalerrors.KRTValidationError
	if errors.As(err, &errInvalidKRT) {
		return &gqlerror.Error{
			Message: errInvalidKRT.Error(),
			Extensions: map[string]interface{}{
				"code": "krt_validation_error",
			},
		}
	}

	var errUnauthorized auth.UnauthorizedError
	if errors.As(err, &errUnauthorized) {
		return &gqlerror.Error{
			Message: errUnauthorized.Error(),
			Extensions: map[string]interface{}{
				"code": "unauthorized",
			},
		}
	}

	return err
}
