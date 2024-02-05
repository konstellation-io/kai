package http

import (
	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/delivery/http/controller"
	kaimiddleware "github.com/konstellation-io/kai/engine/admin-api/delivery/http/middleware"
	"github.com/konstellation-io/kai/engine/admin-api/delivery/http/token"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

// App is the top-level struct.
type App struct {
	server *echo.Echo
	logger logr.Logger
}

const logFormat = "${time_rfc3339} INFO remote_ip=${remote_ip}, method=${method}, uri=${uri}, status=${status}" +
	", bytes_in=${bytes_in}, bytes_out=${bytes_out}, latency=${latency}, referer=${referer}" +
	", user_agent=${user_agent}, error=${error}\n"

// NewApp creates a new App instance.
func NewApp(
	logger logr.Logger,
	gqlController controller.GraphQL,
) *App {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Validator = newCustomValidator()

	e.Use(
		middleware.RequestID(),
		middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: logFormat,
		}),
	)

	if viper.GetBool(config.CORSEnabledKey) {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowCredentials: true,
		}))
	}

	tokenParser := token.NewParser()
	graphqlOperationMiddleware := kaimiddleware.NewGraphQLOperationMiddleware(logger)
	jwtAuthMiddleware := kaimiddleware.NewJwtAuthMiddleware(logger, tokenParser)

	r := e.Group("/graphql")
	r.Use(graphqlOperationMiddleware, jwtAuthMiddleware)
	r.Any("", gqlController.GraphQLHandler)
	r.GET("/playground", gqlController.PlaygroundHandler)

	return &App{
		e,
		logger,
	}
}

// Start runs the HTTP server.
func (a *App) Start() {
	a.logger.Info("HTTP server started", "port", viper.GetString(config.ApplicationPortKey))
	a.server.Logger.Fatal(a.server.Start(":" + viper.GetString(config.ApplicationPortKey)))
}
