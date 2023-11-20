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
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/logging"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/version"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase"
)

type Params struct {
	Logger                 logging.Logger
	Cfg                    *config.Config
	ProductInteractor      *usecase.ProductInteractor
	UserInteractor         *usecase.UserInteractor
	UserActivityInteractor usecase.UserActivityInteracter
	VersionInteractor      *version.Handler
	ServerInfoGetter       *usecase.ServerInfoGetter
	ProcessService         *usecase.ProcessService
}

func NewHTTPHandler(params Params) http.Handler {
	graphQLResolver := NewGraphQLResolver(params)

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

	var errInvalidKRT version.KRTValidationError
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
